package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/typesense/typesense-go/v3/typesense"
	"github.com/typesense/typesense-go/v3/typesense/api"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AliasResource{}
var _ resource.ResourceWithImportState = &AliasResource{}

func NewAliasResource() resource.Resource {
	return &AliasResource{}
}

type AliasResource struct {
	client *typesense.Client
}

type AliasResourceModel struct {
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	CollectionName types.String `tfsdk:"collection_name"`
}

func (r *AliasResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alias"
}

func (r *AliasResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "An alias is a virtual collection name that points to a real collection. If you're familiar with symbolic links on Linux, it's very similar to that.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Id identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name identifier",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"collection_name": schema.StringAttribute{
				MarkdownDescription: "Collection name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *AliasResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*typesense.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *AliasResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AliasResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := &api.CollectionAliasSchema{CollectionName: data.CollectionName.ValueString()}

	alias, err := r.client.Aliases().Upsert(ctx, data.Name.ValueString(), body)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create alias, got error: %s", err))
		return
	}

	data.Id = types.StringPointerValue(alias.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AliasResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AliasResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	alias, err := r.client.Alias(data.Id.ValueString()).Retrieve(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning("Resource Not Found", fmt.Sprintf("Unable to find alias %s, removing from state", data.Id.ValueString()))
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to retrieve alias, got error: %s", err))
		}

		return
	}

	data.Name = types.StringPointerValue(alias.Name)
	data.CollectionName = types.StringValue(alias.CollectionName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AliasResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AliasResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := &api.CollectionAliasSchema{CollectionName: data.CollectionName.ValueString()}

	alias, err := r.client.Aliases().Upsert(ctx, data.Name.ValueString(), body)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create alias, got error: %s", err))
		return
	}

	data.Name = types.StringPointerValue(alias.Name)
	data.CollectionName = types.StringValue(alias.CollectionName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AliasResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AliasResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Warn(ctx, "###Delete alias with id="+data.Id.ValueString())

	_, err := r.client.Alias(data.Id.ValueString()).Delete(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete alias, got error: %s", err))
		}

		return
	}

	data.Id = types.StringValue("")
}

func (r *AliasResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
