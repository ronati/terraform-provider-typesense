package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/typesense/typesense-go/v3/typesense"
	"github.com/typesense/typesense-go/v3/typesense/api"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &CollectionResource{}
var _ resource.ResourceWithImportState = &CollectionResource{}
var _ resource.ResourceWithModifyPlan = &CollectionResource{}

func NewCollectionResource() resource.Resource {
	return &CollectionResource{}
}

type CollectionResource struct {
	client *typesense.Client
}

type CollectionResourceModel struct {
	Id                  types.String                   `tfsdk:"id"`
	Name                types.String                   `tfsdk:"name"`
	DefaultSortingField types.String                   `tfsdk:"default_sorting_field"`
	Fields              []CollectionResourceFieldModel `tfsdk:"fields"`
	EnableNestedFields  types.Bool                     `tfsdk:"enable_nested_fields"`
	SymbolsToIndex      []types.String                 `tfsdk:"symbols_to_index"`
	TokenSeparators     []types.String                 `tfsdk:"token_separators"`
}

type CollectionResourceFieldModel struct {
	Name           types.String `tfsdk:"name"`
	Facet          types.Bool   `tfsdk:"facet"`
	Index          types.Bool   `tfsdk:"index"`
	Optional       types.Bool   `tfsdk:"optional"`
	Sort           types.Bool   `tfsdk:"sort"`
	Infix          types.Bool   `tfsdk:"infix"`
	Type           types.String `tfsdk:"type"`
	Stem           types.Bool   `tfsdk:"stem"`
	StemDictionary types.String `tfsdk:"stem_dictionary"`
	Locale         types.String `tfsdk:"locale"`
	Store          types.Bool   `tfsdk:"store"`
	NumDim         types.Int64  `tfsdk:"num_dim"`
}

func (r *CollectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collection"
}

func (r *CollectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Group of related documents which are roughly equivalent to a table in a relational database. Terraform will still remove auto-created fields for collections with auto-type, so you need to manually update the collection schema to match generated fields",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Id identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Collection name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"default_sorting_field": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Default sorting field",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enable_nested_fields": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Enable nested fields, must be enabled to use object/object[] types",
				Default:             booldefault.StaticBool(false),
			},
			"symbols_to_index": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "List of symbols to index",
				Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"token_separators": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "List of token separators",
				Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"fields": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
						},
						"facet": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Facet field. Defaults to false.",
						},
						"index": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Index field. Defaults to true.",
						},
						"optional": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Optional field. Defaults to false.",
						},
						"sort": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Sort field. Defaults to false.",
						},
						"infix": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Infix field. Defaults to false.",
						},
						"type": schema.StringAttribute{
							Required:    true,
							Description: "Field type.",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"string",
									"int32",
									"int64",
									"float",
									"bool",
									"geopoint",
									"object",
									"string[]",
									"int32[]",
									"int64[]",
									"float[]",
									"bool[]",
									"geopoint[]",
									"object[]",
									"string*",
									"auto",
								),
							},
						},
						"stem": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Enable stemming on field. Defaults to false.",
						},
						"stem_dictionary": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Custom stemming dictionary. Defaults to empty string.",
						},
						"locale": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Locale for language-specific tokenization. Defaults to empty string.",
						},
						"store": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Store field value on disk. Defaults to true.",
						},
						"num_dim": schema.Int64Attribute{
							Optional:    true,
							Description: "Number of dimensions for vector fields (float[] type). Required for vector search.",
						},
					},
				},
			},
		},
	}
}

func (r *CollectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CollectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CollectionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	schema := &api.CollectionSchema{}
	schema.Name = data.Name.ValueString()
	schema.DefaultSortingField = data.DefaultSortingField.ValueStringPointer()
	schema.EnableNestedFields = data.EnableNestedFields.ValueBoolPointer()

	symbolsToIndex := []string{}
	for _, symbol := range data.SymbolsToIndex {
		symbolsToIndex = append(symbolsToIndex, symbol.ValueString())
	}
	schema.SymbolsToIndex = &symbolsToIndex

	tokensSeparators := []string{}
	for _, token := range data.TokenSeparators {
		tokensSeparators = append(tokensSeparators, token.ValueString())
	}
	schema.TokenSeparators = &tokensSeparators

	fields := []api.Field{}

	for _, field := range data.Fields {
		fields = append(fields, filedModelToApiField(field))
	}

	schema.Fields = fields
	collection, err := r.client.Collections().Create(ctx, schema)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create collection, got error: %s", err))
		return
	}

	data.Id = types.StringValue(collection.Name)
	data.Name = types.StringValue(collection.Name)

	if collection.DefaultSortingField != nil && *collection.DefaultSortingField != "" {
		data.DefaultSortingField = types.StringPointerValue(collection.DefaultSortingField)
	}

	data.EnableNestedFields = types.BoolPointerValue(collection.EnableNestedFields)
	data.Fields = flattenCollectionFields(collection.Fields)

	data.SymbolsToIndex = []types.String{}
	if collection.SymbolsToIndex != nil {
		for _, symbol := range *collection.SymbolsToIndex {
			data.SymbolsToIndex = append(data.SymbolsToIndex, types.StringValue(symbol))
		}
	}

	data.TokenSeparators = []types.String{}
	if collection.TokenSeparators != nil {
		for _, token := range *collection.TokenSeparators {
			data.TokenSeparators = append(data.TokenSeparators, types.StringValue(token))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CollectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CollectionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueString()

	collection, err := r.client.Collection(id).Retrieve(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to retrieve collection, got error: %s", err))
		}

		return
	}

	tflog.Info(ctx, "###Got collection name:"+collection.Name)

	data.Id = types.StringValue(collection.Name)
	data.Name = types.StringValue(collection.Name)

	if collection.DefaultSortingField != nil && *collection.DefaultSortingField != "" {
		data.DefaultSortingField = types.StringPointerValue(collection.DefaultSortingField)
	}

	data.EnableNestedFields = types.BoolPointerValue(collection.EnableNestedFields)
	data.Fields = flattenCollectionFields(collection.Fields)

	if collection.SymbolsToIndex != nil {
		data.SymbolsToIndex = []types.String{}
		if collection.SymbolsToIndex != nil {
			for _, symbol := range *collection.SymbolsToIndex {
				data.SymbolsToIndex = append(data.SymbolsToIndex, types.StringValue(symbol))
			}
		}
	}

	if collection.TokenSeparators != nil {
		data.TokenSeparators = []types.String{}
		if collection.TokenSeparators != nil {
			for _, token := range *collection.TokenSeparators {
				data.TokenSeparators = append(data.TokenSeparators, types.StringValue(token))
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func boolPointerValueWithDefault(ptr *bool, defaultVal bool) types.Bool {
	if ptr == nil {
		return types.BoolValue(defaultVal)
	}
	return types.BoolValue(*ptr)
}

func stringPointerValueWithDefault(ptr *string, defaultVal string) types.String {
	if ptr == nil {
		return types.StringValue(defaultVal)
	}
	return types.StringValue(*ptr)
}

func intPointerValue(ptr *int) types.Int64 {
	if ptr == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*ptr))
}

func flattenCollectionFields(fields []api.Field) []CollectionResourceFieldModel {
	if fields != nil {
		fis := make([]CollectionResourceFieldModel, len(fields))

		for i, fieldResponse := range fields {
			var field CollectionResourceFieldModel
			field.Name = types.StringValue(fieldResponse.Name)
			field.Facet = boolPointerValueWithDefault(fieldResponse.Facet, false)
			field.Index = boolPointerValueWithDefault(fieldResponse.Index, true)
			field.Optional = boolPointerValueWithDefault(fieldResponse.Optional, false)
			field.Sort = boolPointerValueWithDefault(fieldResponse.Sort, false)
			field.Infix = boolPointerValueWithDefault(fieldResponse.Infix, false)
			field.Type = types.StringValue(fieldResponse.Type)
			field.Stem = boolPointerValueWithDefault(fieldResponse.Stem, false)
			field.StemDictionary = stringPointerValueWithDefault(fieldResponse.StemDictionary, "")
			field.Locale = stringPointerValueWithDefault(fieldResponse.Locale, "")
			field.Store = boolPointerValueWithDefault(fieldResponse.Store, true)
			field.NumDim = intPointerValue(fieldResponse.NumDim)
			fis[i] = field
		}

		return fis
	}

	return make([]CollectionResourceFieldModel, 0)
}

func (r *CollectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CollectionResourceModel
	var state CollectionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	stateItems := make(map[string]CollectionResourceFieldModel)

	for i := 0; i < len(state.Fields); i += 1 {
		stateItems[state.Fields[i].Name.ValueString()] = state.Fields[i]
	}

	schema := &api.CollectionUpdateSchema{}

	var drop = new(bool)
	*drop = true

	for _, field := range plan.Fields {
		// item not exists, need to create
		if _, ok := stateItems[field.Name.ValueString()]; !ok {
			schema.Fields = append(schema.Fields, filedModelToApiField(field))

			tflog.Info(ctx, "###Field will be created: "+field.Name.ValueString())

		} else if stateItems[field.Name.ValueString()] != field {
			// item was changed, need to update

			schema.Fields = append(schema.Fields,
				api.Field{
					Drop: drop,
					Name: field.Name.ValueString(),
				},
				filedModelToApiField(field))
			tflog.Info(ctx, "###Field will be updated: "+field.Name.ValueString())

		} else {
			// item was not changed, do nothing
			tflog.Info(ctx, "###Field remaining the same: "+field.Name.ValueString())
		}

		// delete processed field from the state object
		delete(stateItems, field.Name.ValueString())
	}

	for _, field := range stateItems {
		schema.Fields = append(schema.Fields,
			api.Field{
				Drop: drop,
				Name: field.Name.ValueString(),
			})
		tflog.Info(ctx, "###Field will be deleted: "+field.Name.ValueString())
	}

	_, err := r.client.Collection(state.Id.ValueString()).Update(ctx, schema)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update collection, got error: %s", err))
		return
	}

	// Read back the updated collection to get all computed field attributes
	collection, err := r.client.Collection(state.Id.ValueString()).Retrieve(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to retrieve updated collection, got error: %s", err))
		return
	}

	plan.Id = types.StringValue(collection.Name)
	plan.Name = types.StringValue(collection.Name)

	if collection.DefaultSortingField != nil && *collection.DefaultSortingField != "" {
		plan.DefaultSortingField = types.StringPointerValue(collection.DefaultSortingField)
	}

	plan.EnableNestedFields = types.BoolPointerValue(collection.EnableNestedFields)
	plan.Fields = flattenCollectionFields(collection.Fields)

	plan.SymbolsToIndex = []types.String{}
	if collection.SymbolsToIndex != nil {
		for _, symbol := range *collection.SymbolsToIndex {
			plan.SymbolsToIndex = append(plan.SymbolsToIndex, types.StringValue(symbol))
		}
	}

	plan.TokenSeparators = []types.String{}
	if collection.TokenSeparators != nil {
		for _, token := range *collection.TokenSeparators {
			plan.TokenSeparators = append(plan.TokenSeparators, types.StringValue(token))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CollectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CollectionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Warn(ctx, "###Delete collection with id="+data.Id.ValueString())

	_, err := r.client.Collection(data.Id.ValueString()).Delete(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete collection, got error: %s", err))
		}

		return
	}

	data.Id = types.StringValue("")
}

func (r *CollectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *CollectionResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan CollectionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	modified := false
	for i := range plan.Fields {
		if plan.Fields[i].Facet.IsUnknown() || plan.Fields[i].Facet.IsNull() {
			plan.Fields[i].Facet = types.BoolValue(false)
			modified = true
		}
		if plan.Fields[i].Index.IsUnknown() || plan.Fields[i].Index.IsNull() {
			plan.Fields[i].Index = types.BoolValue(true)
			modified = true
		}
		if plan.Fields[i].Optional.IsUnknown() || plan.Fields[i].Optional.IsNull() {
			plan.Fields[i].Optional = types.BoolValue(false)
			modified = true
		}
		if plan.Fields[i].Sort.IsUnknown() || plan.Fields[i].Sort.IsNull() {
			plan.Fields[i].Sort = types.BoolValue(false)
			modified = true
		}
		if plan.Fields[i].Infix.IsUnknown() || plan.Fields[i].Infix.IsNull() {
			plan.Fields[i].Infix = types.BoolValue(false)
			modified = true
		}
		if plan.Fields[i].Stem.IsUnknown() || plan.Fields[i].Stem.IsNull() {
			plan.Fields[i].Stem = types.BoolValue(false)
			modified = true
		}
		if plan.Fields[i].StemDictionary.IsUnknown() || plan.Fields[i].StemDictionary.IsNull() {
			plan.Fields[i].StemDictionary = types.StringValue("")
			modified = true
		}
		if plan.Fields[i].Locale.IsUnknown() || plan.Fields[i].Locale.IsNull() {
			plan.Fields[i].Locale = types.StringValue("")
			modified = true
		}
		if plan.Fields[i].Store.IsUnknown() || plan.Fields[i].Store.IsNull() {
			plan.Fields[i].Store = types.BoolValue(true)
			modified = true
		}
	}

	if modified {
		resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
	}
}

func filedModelToApiField(field CollectionResourceFieldModel) api.Field {
	apiField := api.Field{
		Name:           field.Name.ValueString(),
		Facet:          field.Facet.ValueBoolPointer(),
		Index:          field.Index.ValueBoolPointer(),
		Optional:       field.Optional.ValueBoolPointer(),
		Sort:           field.Sort.ValueBoolPointer(),
		Infix:          field.Infix.ValueBoolPointer(),
		Type:           field.Type.ValueString(),
		Stem:           field.Stem.ValueBoolPointer(),
		StemDictionary: field.StemDictionary.ValueStringPointer(),
		Locale:         field.Locale.ValueStringPointer(),
		Store:          field.Store.ValueBoolPointer(),
	}

	if !field.NumDim.IsNull() && !field.NumDim.IsUnknown() {
		numDim := int(field.NumDim.ValueInt64())
		apiField.NumDim = &numDim
	}

	return apiField
}
