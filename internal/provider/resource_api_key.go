package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/typesense/typesense-go/v3/typesense"
	"github.com/typesense/typesense-go/v3/typesense/api"
)

var _ resource.Resource = &ApiKeyResource{}
var _ resource.ResourceWithImportState = &ApiKeyResource{}

func NewApiKeyResource() resource.Resource {
	return &ApiKeyResource{}
}

type ApiKeyResource struct {
	client *typesense.Client
}

type ApiKeyResourceModel struct {
	Id          types.String   `tfsdk:"id"`
	Description types.String   `tfsdk:"description"`
	Actions     []types.String `tfsdk:"actions"`
	Collections []types.String `tfsdk:"collections"`
	ExpiresAt   types.Int64    `tfsdk:"expires_at"`
	Value       types.String   `tfsdk:"value"`
	ValuePrefix types.String   `tfsdk:"value_prefix"`
}

func (r *ApiKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key"
}

func (r *ApiKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "API Key resource for accessing Typesense collections with specific permissions",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "ID identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the API key",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"actions": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of actions this API key can perform (e.g., documents:search, documents:insert, collections:create, * for all actions)",
				Required:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"collections": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of collections this API key can access (use ['*'] for all collections)",
				Required:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"expires_at": schema.Int64Attribute{
				MarkdownDescription: "Unix timestamp when the API key expires (optional)",
				Optional:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "The actual API key value. If not provided, Typesense will auto-generate one.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"value_prefix": schema.StringAttribute{
				MarkdownDescription: "First few characters of the API key for identification",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *ApiKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*typesense.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *typesense.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *ApiKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ApiKeyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	actions := []string{}
	for _, action := range data.Actions {
		actions = append(actions, action.ValueString())
	}

	collections := []string{}
	for _, collection := range data.Collections {
		collections = append(collections, collection.ValueString())
	}

	keySchema := &api.ApiKeySchema{
		Description: data.Description.ValueString(),
		Actions:     actions,
		Collections: collections,
	}

	if !data.Value.IsNull() && data.Value.ValueString() != "" {
		value := data.Value.ValueString()
		keySchema.Value = &value
	}

	if !data.ExpiresAt.IsNull() {
		expiresAt := data.ExpiresAt.ValueInt64()
		keySchema.ExpiresAt = &expiresAt
	}

	apiKey, err := r.client.Keys().Create(ctx, keySchema)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create API key, got error: %s", err))
		return
	}

	data.Id = types.StringValue(strconv.FormatInt(*apiKey.Id, 10))
	data.Description = types.StringValue(apiKey.Description)

	data.Actions = []types.String{}
	for _, action := range apiKey.Actions {
		data.Actions = append(data.Actions, types.StringValue(action))
	}

	data.Collections = []types.String{}
	for _, collection := range apiKey.Collections {
		data.Collections = append(data.Collections, types.StringValue(collection))
	}

	// Only set expires_at if it was explicitly set in the configuration
	// Typesense returns a very high default value (64723363199) when no expiration is set
	if !data.ExpiresAt.IsNull() && apiKey.ExpiresAt != nil {
		data.ExpiresAt = types.Int64Value(*apiKey.ExpiresAt)
	}
	// Keep expires_at null if it wasn't set by user, even if API returns default value

	if apiKey.Value != nil {
		data.Value = types.StringValue(*apiKey.Value)
	}

	if apiKey.ValuePrefix != nil {
		data.ValuePrefix = types.StringValue(*apiKey.ValuePrefix)
	} else {
		data.ValuePrefix = types.StringValue("")
	}

	tflog.Info(ctx, "Created API key with ID: "+data.Id.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApiKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ApiKeyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueString()
	keyId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse API key ID: %s", err))
		return
	}

	apiKey, err := r.client.Key(keyId).Retrieve(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to retrieve API key, got error: %s", err))
		}
		return
	}

	data.Description = types.StringValue(apiKey.Description)

	data.Actions = []types.String{}
	for _, action := range apiKey.Actions {
		data.Actions = append(data.Actions, types.StringValue(action))
	}

	data.Collections = []types.String{}
	for _, collection := range apiKey.Collections {
		data.Collections = append(data.Collections, types.StringValue(collection))
	}

	// Only set expires_at if it was explicitly set in the configuration
	// Typesense returns a very high default value (64723363199) when no expiration is set
	if !data.ExpiresAt.IsNull() && apiKey.ExpiresAt != nil {
		data.ExpiresAt = types.Int64Value(*apiKey.ExpiresAt)
	}
	// Keep expires_at null if it wasn't set by user, even if API returns default value

	if apiKey.ValuePrefix != nil {
		data.ValuePrefix = types.StringValue(*apiKey.ValuePrefix)
	} else {
		data.ValuePrefix = types.StringValue("")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApiKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"API keys cannot be updated. Please delete and recreate the resource to make changes.",
	)
}

func (r *ApiKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ApiKeyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueString()
	keyId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse API key ID: %s", err))
		return
	}

	tflog.Info(ctx, "Delete API key with id="+id)

	_, err = r.client.Key(keyId).Delete(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete API key, got error: %s", err))
		}
		return
	}

	data.Id = types.StringValue("")
}

func (r *ApiKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
