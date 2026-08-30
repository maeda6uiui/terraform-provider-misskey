package note

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/maeda6uiui/terraform-provider-misskey/internal/provider/misskey"
	"github.com/maeda6uiui/terraform-provider-misskey/internal/provider/model"
)

type NoteResource struct {
	misskeyClient *misskey.MisskeyHttpClient
}

func NewNoteResource() resource.Resource {
	return &NoteResource{}
}

func (r *NoteResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_note"
}

func (r *NoteResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*model.MisskeyProviderModel)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *MisskeyProviderModel, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	misskeyClient := misskey.NewMisskeyHttpClient(
		data.ServerUrl.ValueString(),
		int(data.TimeoutSeconds.ValueInt32()),
		data.AccessToken.ValueString(),
	)
	r.misskeyClient = misskeyClient
}

func (r *NoteResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Id of the note",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"text": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Text of the note",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"visibility": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "Visibility of the note. " +
					"Allowed values are `public`, `home`, `followers`, and `specified`. ",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("public", "home", "followers", "specified"),
				},
			},
			"visible_user_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				MarkdownDescription: "Note is visible to the users specified. " +
					"Only valid if `visibility` is `specified`.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *NoteResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse) {
	var model NoteResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var visibleUserIds []string
	resp.Diagnostics.Append(model.VisibleUserIds.ElementsAs(ctx, &visibleUserIds, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody := &misskey.CreateNoteRequest{
		Text:           model.Text.ValueString(),
		Visibility:     model.Visibility.ValueString(),
		VisibleUserIds: visibleUserIds,
	}
	respBody, respStatus, err := r.misskeyClient.Post("/api/notes/create", &reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to create a note, got error: %s", err),
		)
		return
	}
	if respStatus != http.StatusOK {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to create a note, got status code: %d", respStatus),
		)
		return
	}

	var respModel *misskey.CreateNoteResponse
	if err := json.Unmarshal(respBody, &respModel); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Response Format",
			"Client returned a response that cannot be parsed",
		)
		return
	}

	model.Id = types.StringValue(respModel.CreatedNote.Id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *NoteResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse) {
	var model NoteResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody := &misskey.ShowNoteRequest{
		NoteId: model.Id.ValueString(),
	}
	respBody, respStatus, err := r.misskeyClient.Post("/api/notes/show", &reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to get a note, got error: %s", err),
		)
		return
	}
	if respStatus != http.StatusOK {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to get a note, got status code: %d", respStatus),
		)
		return
	}

	var respModel *misskey.Note
	if err := json.Unmarshal(respBody, &respModel); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Response Format",
			"Client returned a response that cannot be parsed",
		)
		return
	}

	model.Id = types.StringValue(respModel.Id)
	model.Text = types.StringValue(respModel.Text)
	model.Visibility = types.StringValue(respModel.Visibility)

	visibleUserIds, diags := types.ListValueFrom(ctx, types.StringType, respModel.VisibleUserIds)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	model.VisibleUserIds = visibleUserIds

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *NoteResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse) {
	// Every attribute requires replacement, so this should never be called
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Misskey note cannot be updated in place, it must be recreated instead",
	)
}

func (r *NoteResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse) {
	var model NoteResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody := &misskey.DeleteNoteRequest{
		NoteId: model.Id.ValueString(),
	}
	_, respStatus, err := r.misskeyClient.Post("/api/notes/delete", &reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to delete a note, got error: %s", err),
		)
		return
	}
	if respStatus != http.StatusNoContent {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Failed to delete a note, got status code: %d", respStatus),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
