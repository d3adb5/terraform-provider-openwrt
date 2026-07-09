package lucirpcglue

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/ORFops/terraform-provider-openwrt/lucirpc"
)

const (
	orderingIdAttributeDescription = "Id of this ordering. Always the UCI config and section type."

	orderingIdsAttribute            = "ids"
	orderingIdsAttributeDescription = "Ids of the sections in the desired order. The sections are moved to the end of the config in this order; sections not listed are left untouched."
)

var (
	_ frameworkresource.Resource              = &orderingResource{}
	_ frameworkresource.ResourceWithConfigure = &orderingResource{}
)

// NewOrderingResource constructs a resource that manages the relative order
// of the sections of a single UCI type within a config.
func NewOrderingResource(
	schemaDescription string,
	uciConfig string,
	uciType string,
) frameworkresource.Resource {
	return &orderingResource{
		schemaDescription: schemaDescription,
		uciConfig:         uciConfig,
		uciType:           uciType,
	}
}

type orderingResource struct {
	client            lucirpc.Client
	fullTypeName      string
	schemaDescription string
	uciConfig         string
	uciType           string
}

type orderingModel struct {
	Id  types.String `tfsdk:"id"`
	Ids types.List   `tfsdk:"ids"`
}

// Configure adds the provider configured client to the resource.
func (d *orderingResource) Configure(
	ctx context.Context,
	req frameworkresource.ConfigureRequest,
	res *frameworkresource.ConfigureResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Configuring %s.%s ordering resource", d.uciConfig, d.uciType))
	if req.ProviderData == nil {
		tflog.Debug(ctx, "No provider data")
		return
	}

	providerData, diagnostics := ParseProviderData(ConfigureRequest(req))
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	d.client = providerData.Client
	d.fullTypeName = d.getFullTypeName(providerData.TypeName)
}

// Create enforces the desired order and sets the initial Terraform state.
func (d *orderingResource) Create(
	ctx context.Context,
	req frameworkresource.CreateRequest,
	res *frameworkresource.CreateResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Creating %s resource", d.fullTypeName))

	var model orderingModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if res.Diagnostics.HasError() {
		return
	}

	d.reorder(ctx, model, &res.Diagnostics)
	if res.Diagnostics.HasError() {
		return
	}

	model.Id = types.StringValue(fmt.Sprintf("%s.%s", d.uciConfig, d.uciType))
	res.Diagnostics.Append(res.State.Set(ctx, model)...)
}

// Delete removes the Terraform state.
// The sections keep their current order.
func (d *orderingResource) Delete(
	ctx context.Context,
	req frameworkresource.DeleteRequest,
	res *frameworkresource.DeleteResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Deleting %s resource", d.fullTypeName))
}

// Metadata sets the resource type name.
func (d *orderingResource) Metadata(
	ctx context.Context,
	req frameworkresource.MetadataRequest,
	res *frameworkresource.MetadataResponse,
) {
	res.TypeName = d.getFullTypeName(req.ProviderTypeName)
}

// Read refreshes the Terraform state with the current order of the sections.
func (d *orderingResource) Read(
	ctx context.Context,
	req frameworkresource.ReadRequest,
	res *frameworkresource.ReadResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Reading %s resource", d.fullTypeName))

	var model orderingModel
	res.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if res.Diagnostics.HasError() {
		return
	}

	var ids []string
	res.Diagnostics.Append(model.Ids.ElementsAs(ctx, &ids, false)...)
	if res.Diagnostics.HasError() {
		return
	}

	sections, err := d.client.GetSections(ctx, d.uciConfig, d.uciType)
	if err != nil {
		res.Diagnostics.AddError(
			fmt.Sprintf("problem getting %s.%s sections", d.uciConfig, d.uciType),
			err.Error(),
		)
		return
	}

	// Keep only the sections this resource manages,
	// in the order they currently appear in the config.
	managed := map[string]bool{}
	for _, id := range ids {
		managed[id] = true
	}

	current := []string{}
	for _, section := range sections {
		name, err := section.GetString(idUCISection)
		if err != nil {
			res.Diagnostics.AddError(
				fmt.Sprintf("problem parsing %s.%s section name", d.uciConfig, d.uciType),
				err.Error(),
			)
			return
		}

		if managed[name] {
			current = append(current, name)
		}
	}

	if len(current) == 0 {
		tflog.Warn(ctx, fmt.Sprintf("No %s.%s sections managed by this resource exist; removing from state", d.uciConfig, d.uciType))
		res.State.RemoveResource(ctx)
		return
	}

	value, diagnostics := types.ListValueFrom(ctx, types.StringType, current)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	model.Ids = value
	res.Diagnostics.Append(res.State.Set(ctx, model)...)
}

// Schema defines the schema for the resource.
func (d *orderingResource) Schema(
	ctx context.Context,
	req frameworkresource.SchemaRequest,
	res *frameworkresource.SchemaResponse,
) {
	res.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			IdAttribute: schema.StringAttribute{
				Computed:    true,
				Description: orderingIdAttributeDescription,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			orderingIdsAttribute: schema.ListAttribute{
				Description: orderingIdsAttributeDescription,
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.UniqueValues(),
				},
			},
		},
		Description: d.schemaDescription,
	}
}

// Update enforces the desired order and sets the Terraform state on success.
func (d *orderingResource) Update(
	ctx context.Context,
	req frameworkresource.UpdateRequest,
	res *frameworkresource.UpdateResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Updating %s resource", d.fullTypeName))

	var model orderingModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if res.Diagnostics.HasError() {
		return
	}

	d.reorder(ctx, model, &res.Diagnostics)
	if res.Diagnostics.HasError() {
		return
	}

	model.Id = types.StringValue(fmt.Sprintf("%s.%s", d.uciConfig, d.uciType))
	res.Diagnostics.Append(res.State.Set(ctx, model)...)
}

func (d *orderingResource) getFullTypeName(
	providerTypeName string,
) string {
	uciConfig := strings.ReplaceAll(d.uciConfig, "-", "_")
	uciType := strings.ReplaceAll(d.uciType, "-", "_")
	return fmt.Sprintf("%s_%s_%s_ordering", providerTypeName, uciConfig, uciType)
}

func (d *orderingResource) reorder(
	ctx context.Context,
	model orderingModel,
	diagnostics *diag.Diagnostics,
) {
	var ids []string
	diagnostics.Append(model.Ids.ElementsAs(ctx, &ids, false)...)
	if diagnostics.HasError() {
		return
	}

	result, err := d.client.ReorderSections(ctx, d.uciConfig, d.uciType, ids)
	if err != nil {
		diagnostics.AddError(
			fmt.Sprintf("problem reordering %s.%s sections", d.uciConfig, d.uciType),
			err.Error(),
		)
		return
	}

	if !result {
		diagnostics.AddError(
			fmt.Sprintf("Could not reorder %s.%s sections", d.uciConfig, d.uciType),
			"It is not currently known why this happens. It is unclear if this is a problem with the provider. Please double check the values provided are acceptable.",
		)
	}
}
