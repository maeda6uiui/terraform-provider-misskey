package note

import "github.com/hashicorp/terraform-plugin-framework/types"

type NoteResourceModel struct {
	Id             types.String `tfsdk:"id"`
	Text           types.String `tfsdk:"text"`
	Visibility     types.String `tfsdk:"visibility"`
	VisibleUserIds types.List   `tfsdk:"visible_user_ids"`
}
