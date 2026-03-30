// Package settings implements the openwrt_acme_acme resource and data source,
// which manages the global ACME account configuration (UCI: /etc/config/acme,
// section type: acme).
package settings

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ORFops/terraform-provider-openwrt/lucirpc"
	"github.com/ORFops/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	accountEmailAttribute            = "account_email"
	accountEmailAttributeDescription = "Email address associated with the ACME account. Used by the CA to send expiry notices."
	accountEmailUCIOption            = "account_email"

	debugAttribute            = "debug"
	debugAttributeDescription = "Enable verbose debug logging for ACME operations."
	debugUCIOption            = "debug"

	schemaDescription = "Global ACME account settings (UCI: /etc/config/acme, section type: acme). Manages the account e-mail and debug flag shared by all certificates."

	uciConfig = "acme"
	uciType   = "acme"
)

var (
	accountEmailSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       accountEmailAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetAccountEmail, accountEmailAttribute, accountEmailUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetAccountEmail, accountEmailAttribute, accountEmailUCIOption),
		Validators:        []validator.String{stringvalidator.LengthAtLeast(1)},
	}

	debugSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       debugAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetDebug, debugAttribute, debugUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetDebug, debugAttribute, debugUCIOption),
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		accountEmailAttribute:   accountEmailSchemaAttribute,
		debugAttribute:          debugSchemaAttribute,
		lucirpcglue.IdAttribute: lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
	}
)

func NewDataSource() datasource.DataSource {
	return lucirpcglue.NewDataSource(
		modelGetId,
		schemaAttributes,
		schemaDescription,
		uciConfig,
		uciType,
	)
}

func NewResource() resource.Resource {
	return lucirpcglue.NewResource(
		modelGetId,
		schemaAttributes,
		schemaDescription,
		uciConfig,
		uciType,
	)
}

type model struct {
	AccountEmail types.String `tfsdk:"account_email"`
	Debug        types.Bool   `tfsdk:"debug"`
	Id           types.String `tfsdk:"id"`
}

func modelGetAccountEmail(m model) types.String { return m.AccountEmail }
func modelGetDebug(m model) types.Bool          { return m.Debug }
func modelGetId(m model) types.String           { return m.Id }

func modelSetAccountEmail(m *model, value types.String) { m.AccountEmail = value }
func modelSetDebug(m *model, value types.Bool)          { m.Debug = value }
func modelSetId(m *model, value types.String)           { m.Id = value }
