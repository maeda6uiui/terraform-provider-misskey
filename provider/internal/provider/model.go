package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type MisskeyProviderModel struct {
	ServerUrl      types.String `tfsdk:"server_url"`
	AccessToken    types.String `tfsdk:"access_token"`
	TimeoutSeconds types.Int32  `tfsdk:"timeout_seconds"`
}
