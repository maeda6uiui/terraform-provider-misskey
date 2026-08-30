package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/maeda6uiui/terraform-provider-misskey/internal/provider/model"
	"github.com/maeda6uiui/terraform-provider-misskey/internal/provider/note"
)

type MisskeyProvider struct {
	version string
}

var _ provider.Provider = &MisskeyProvider{}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MisskeyProvider{
			version: version,
		}
	}
}

func (p *MisskeyProvider) Metadata(
	ctx context.Context,
	req provider.MetadataRequest,
	resp *provider.MetadataResponse) {
	resp.TypeName = "misskey"
	resp.Version = p.version
}

func (p *MisskeyProvider) Schema(
	ctx context.Context,
	req provider.SchemaRequest,
	resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"server_url": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "URL of the Misskey server",
			},
			"access_token": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Access token for the Misskey account",
			},
		},
	}
}

func (p *MisskeyProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse) {
	var model model.MisskeyProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var accessToken string
	if !model.AccessToken.IsNull() && !model.AccessToken.IsUnknown() {
		accessToken = model.AccessToken.ValueString()
	} else {
		accessToken = os.Getenv("MISSKEY_ACCESS_TOKEN")
	}

	if accessToken == "" {
		resp.Diagnostics.AddError(
			"Missing Access Token",
			"The provider argument 'access_token' or the environment variable 'MISSKEY_ACCESS_TOKEN' must be set.",
		)
		return
	}

	model.AccessToken = types.StringValue(accessToken)

	resp.ResourceData = model
}

func (p *MisskeyProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		note.NewNoteResource,
	}
}

func (p *MisskeyProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
