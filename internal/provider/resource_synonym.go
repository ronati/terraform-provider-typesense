package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/typesense/typesense-go/v3/typesense"
	"github.com/typesense/typesense-go/v3/typesense/api"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SynonymResource{}
var _ resource.ResourceWithImportState = &SynonymResource{}

func NewSynonymResource() resource.Resource {
	return &SynonymResource{}
}

type SynonymResource struct {
	client *typesense.Client
}

type SynonymResourceModel struct {
	Id             types.String   `tfsdk:"id"`
	Name           types.String   `tfsdk:"name"`
	CollectionName types.String   `tfsdk:"collection_name"`
	Root           types.String   `tfsdk:"root"`
	Synonyms       []types.String `tfsdk:"synonyms"`
}

func (r *SynonymResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_synonym"
}

func (r *SynonymResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The synonyms feature allows you to define search terms that should be considered equivalent. For eg: when you define a synonym for sneaker as shoe, searching for sneaker will now return all records with the word shoe in them, in addition to records with the word sneaker.",

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
			"root": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "For 1-way synonyms, indicates the root word that words in the synonyms parameter map to",
			},
			"synonyms": schema.ListAttribute{
				Required:            true,
				ElementType:         types.StringType,
				Validators:          []validator.List{listvalidator.SizeAtLeast(1)},
				MarkdownDescription: "Array of words that should be considered as synonyms.",
			},
		},
	}
}

func (r *SynonymResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SynonymResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SynonymResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	schema := &api.SearchSynonymSchema{}

	schema.Root = data.Root.ValueStringPointer()
	schema.Synonyms = convertTerraformArrayToStringArray(data.Synonyms)

	tflog.Info(ctx, "synonyms: "+fmt.Sprint(schema.Synonyms))
	tflog.Info(ctx, "collection name: "+data.CollectionName.ValueString())

	synonym, err := r.client.Collection(data.CollectionName.ValueString()).Synonyms().Upsert(ctx, data.Name.ValueString(), schema)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create synonym, got error: %s", err))
		return
	}

	data.Id = types.StringValue(createId(data.CollectionName.ValueString(), *synonym.Id))
	data.Root = types.StringPointerValue(synonym.Root)
	data.Synonyms = convertStringArrayToTerraformArray(synonym.Synonyms)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SynonymResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SynonymResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	collectionName, id, parseError := splitCollectionRelatedId(data.Id.ValueString())
	if parseError != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Unable to split resource ID: %s", parseError))
		return
	}

	synonym, err := r.client.Collection(collectionName).Synonym(id).Retrieve(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning("Resource Not Found", fmt.Sprintf("Unable to find synonym %s, removing from state", data.Id.ValueString()))
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to retrieve synonym, got error: %s", err))
		}

		return
	}

	// data.Id = types.StringPointerValue(synonym.Id)
	data.Name = types.StringPointerValue(synonym.Id)
	data.Synonyms = convertStringArrayToTerraformArray(synonym.Synonyms)

	if synonym.Root != nil && *synonym.Root != "" {
		data.Root = types.StringPointerValue(synonym.Root)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SynonymResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SynonymResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	schema := &api.SearchSynonymSchema{}

	schema.Root = data.Root.ValueStringPointer()
	schema.Synonyms = convertTerraformArrayToStringArray(data.Synonyms)

	tflog.Info(ctx, "synonyms: "+fmt.Sprint(schema.Synonyms))
	tflog.Info(ctx, "collection name: "+data.CollectionName.ValueString())

	collectionName, id, parseError := splitCollectionRelatedId(data.Id.ValueString())
	if parseError != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Unable to split resource ID: %s", parseError))
		return
	}

	synonym, err := r.client.Collection(collectionName).Synonyms().Upsert(ctx, id, schema)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create synonym, got error: %s", err))
		return
	}

	// data.Id = types.StringPointerValue(synonym.Id)
	data.Name = types.StringPointerValue(synonym.Id)
	data.Root = types.StringPointerValue(synonym.Root)
	data.Synonyms = convertStringArrayToTerraformArray(synonym.Synonyms)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SynonymResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SynonymResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Warn(ctx, "###Delete Synonym with id="+data.Id.ValueString())

	collectionName, id, parseError := splitCollectionRelatedId(data.Id.ValueString())
	if parseError != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Unable to split resource ID: %s", parseError))
		return
	}

	_, err := r.client.Collection(collectionName).Synonym(id).Delete(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete synonym, got error: %s", err))
		}

		return
	}

	data.Id = types.StringValue("")
}

func (r *SynonymResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// ID format is: collection_name.synonym_id
	collectionName, synonymId, err := splitCollectionRelatedId(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be in format 'collection_name.synonym_id', got: %s", req.ID),
		)
		return
	}

	// Set both the ID and collection_name
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("collection_name"), collectionName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), synonymId)...)
}
