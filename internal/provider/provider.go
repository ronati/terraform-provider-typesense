package provider

import (
	"context"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/typesense/typesense-go/v3/typesense"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &TypesenseProvider{}
)
var _ provider.ProviderWithFunctions = &TypesenseProvider{}

type TypesenseProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// TypesenseProviderModel is the provider implementation.
type TypesenseProviderModel struct {
	ApiKey     types.String `tfsdk:"api_key"`
	ApiAddress types.String `tfsdk:"api_address"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &TypesenseProvider{
			version: version,
		}
	}
}

// Metadata returns the provider type name.
func (p *TypesenseProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "typesense"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *TypesenseProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "API Key to access the Typesense server. This can also be set via the `TYPESENSE_API_KEY` environment variable.",
			},
			"api_address": schema.StringAttribute{
				Optional:    true,
				Description: "URL of the Typesense server. This can also be set via the `TYPESENSE_API_ADDRESS` environment variable.",
			},
		},
	}
}

// Configure prepares a typesense API client for data sources and resources.
func (p *TypesenseProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data TypesenseProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	api_key := os.Getenv("TYPESENSE_API_KEY")
	api_address := os.Getenv("TYPESENSE_API_ADDRESS")

	if !data.ApiKey.IsNull() {
		api_key = data.ApiKey.ValueString()
	}

	if !data.ApiAddress.IsNull() {
		api_address = data.ApiAddress.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if api_key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing Typesense API Key",
			"The provider cannot create the Typesense API client as there is a missing or empty value for the Typesense API key. "+
				"Set the api_key value in the configuration or use the TYPESENSE_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if api_address == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_address"),
			"Missing Typesense API Address",
			"The provider cannot create the Typesense API client as there is a missing or empty value for the Typesense API host. "+
				"Set the api_address value in the configuration or use the TYPESENSE_API_ADDRESS environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new typesense client using the configuration values
	client := typesense.NewClient(
		typesense.WithServer(api_address),
		typesense.WithAPIKey(api_key),
		typesense.WithConnectionTimeout(5*time.Minute),
		typesense.WithConnectionTimeout(5*time.Second),
		typesense.WithCircuitBreakerMaxRequests(50),
		typesense.WithCircuitBreakerInterval(2*time.Minute),
		typesense.WithCircuitBreakerTimeout(1*time.Minute))

	// Make the Typesense client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// Resources defines the resources implemented in the provider.
func (p *TypesenseProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCollectionResource,
		NewSynonymResource,
		NewDocumentResource,
		NewAliasResource,
		NewApiKeyResource,
	}
}

// DataSources defines the data sources implemented in the provider.
func (p *TypesenseProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Functions implements provider.ProviderWithFunctions.
func (p *TypesenseProvider) Functions(context.Context) []func() function.Function {
	return nil
}
