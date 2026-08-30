package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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
	resp.Schema = schema.Schema{}
}

func (p *MisskeyProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse) {

}

func (p *MisskeyProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *MisskeyProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
