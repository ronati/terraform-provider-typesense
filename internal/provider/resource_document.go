package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
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
var _ resource.Resource = &DocumentResource{}
var _ resource.ResourceWithImportState = &DocumentResource{}

func NewDocumentResource() resource.Resource {
	return &DocumentResource{}
}

type DocumentResource struct {
	client *typesense.Client
}

type DocumentResourceModel struct {
	Id             types.String         `tfsdk:"id"`
	Name           types.String         `tfsdk:"name"`
	CollectionName types.String         `tfsdk:"collection_name"`
	Document       jsontypes.Normalized `tfsdk:"document"`
}

func (r *DocumentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_document"
}

func (r *DocumentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Every record you index in Typesense is called a Document",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Id identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name identifier, it will be used as id, so needs to be URL-friendly",
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
			"document": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Document object in JSON format",
				CustomType:          jsontypes.NormalizedType{},
			},
		},
	}
}

func (r *DocumentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DocumentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DocumentResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	document, err := parseJsonStringToMap(data.Document.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("JSON format error", fmt.Sprintf("Unable to parse document json, got error: %s", err))
		return
	}

	document["id"] = data.Name.ValueString()

	result, err := r.client.Collection(data.CollectionName.ValueString()).Documents().Create(ctx, document, &api.DocumentIndexParameters{})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create document, got error: %s", err))
		return
	}

	docId := result["id"].(string)
	data.Id = types.StringValue(createId(data.CollectionName.ValueString(), docId))

	// Read back the document to ensure consistent JSON formatting
	retrievedDoc, err := r.client.Collection(data.CollectionName.ValueString()).Document(docId).Retrieve(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to retrieve created document, got error: %s", err))
		return
	}

	data.Name = types.StringValue(retrievedDoc["id"].(string))
	delete(retrievedDoc, "id")

	data.Document, err = parseMapToJsonString(retrievedDoc)
	if err != nil {
		resp.Diagnostics.AddError("JSON format error", fmt.Sprintf("Unable to parse json response, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DocumentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DocumentResourceModel

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

	result, err := r.client.Collection(collectionName).Document(id).Retrieve(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning("Resource Not Found", fmt.Sprintf("Unable to retrieve document, got error: %s", err))
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to retrieve document, got error: %s", err))
		}

		return
	}

	// data.Id = types.StringValue(result["id"].(string))
	data.Name = types.StringValue(result["id"].(string))

	delete(result, "id")

	data.Document, err = parseMapToJsonString(result)

	if err != nil {
		resp.Diagnostics.AddError("JSON format error", fmt.Sprintf("Unable to parse json response, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DocumentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DocumentResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	document, err := parseJsonStringToMap(data.Document.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("JSON format error", fmt.Sprintf("Unable to parse document json, got error: %s", err))
		return
	}

	collectionName, id, parseError := splitCollectionRelatedId(data.Id.ValueString())
	if parseError != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Unable to split resource ID: %s", parseError))
		return
	}

	document["id"] = id

	result, err := r.client.Collection(collectionName).Document(id).Update(ctx, document, &api.DocumentIndexParameters{})
	_ = result // result is empty

	if err != nil {

		//check if error contains 201 response
		if strings.Contains(err.Error(), "201") {
			//ignore, sometimes typesense returns 201 code
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update document, got error: %s", err))
			return
		}

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DocumentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DocumentResourceModel

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

	tflog.Warn(ctx, "###Delete Document with id="+data.Id.ValueString())

	_, err := r.client.Collection(collectionName).Document(id).Delete(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning("Resource Not Found", fmt.Sprintf("Unable to delete document, got error: %s", err))
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete document, got error: %s", err))
		}

		return
	}

	data.Id = types.StringValue("")
}

func (r *DocumentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// ID format is: collection_name.document_id
	collectionName, documentId, err := splitCollectionRelatedId(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Import ID must be in format 'collection_name.document_id', got: %s", req.ID),
		)
		return
	}

	// Set both the ID and collection_name
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("collection_name"), collectionName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), documentId)...)
}
