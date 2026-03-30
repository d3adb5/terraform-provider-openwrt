// Package cert implements the openwrt_acme_cert resource and data source,
// which manages individual ACME-managed certificates (UCI: /etc/config/acme,
// section type: cert).
package cert

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ORFops/terraform-provider-openwrt/lucirpc"
	"github.com/ORFops/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	acmeServerAttribute            = "acme_server"
	acmeServerAttributeDescription = "URL of the ACME directory endpoint to use as the CA. Defaults to Let's Encrypt when unset."
	acmeServerUCIOption            = "acme_server"

	caliasAttribute            = "calias"
	caliasAttributeDescription = "Challenge alias domain: all DNS challenge records are written under this domain instead of the actual certificate domain. Used with the DNS alias mode."
	caliasUCIOption            = "calias"

	certProfileAttribute            = "cert_profile"
	certProfileAttributeDescription = "Certificate profile to request from the CA (provider-specific, e.g. \"tlsserver\")."
	certProfileUCIOption            = "cert_profile"

	credentialsAttribute            = "credentials"
	credentialsAttributeDescription = "DNS provider API credentials as shell-variable assignments (e.g. [\"CF_Email=user@example.com\", \"CF_Key=abc123\"]). Required when validation_method is \"dns\"."
	credentialsUCIOption            = "credentials"

	daliasAttribute            = "dalias"
	daliasAttributeDescription = "Domain alias: the certificate's DNS records are managed under this domain. Similar to calias but applies to all challenge types."
	daliasUCIOption            = "dalias"

	daysAttribute            = "days"
	daysAttributeDescription = "Number of days before certificate expiry at which renewal is attempted. Default: 60."
	daysUCIOption            = "days"

	dnsAttribute            = "dns"
	dnsAttributeDescription = "DNS API provider to use for DNS-01 challenges (e.g. \"dns_cf\" for Cloudflare, \"dns_duckdns\" for DuckDNS). Required when validation_method is \"dns\"."
	dnsUCIOption            = "dns"

	dnsWaitAttribute            = "dns_wait"
	dnsWaitAttributeDescription = "Seconds to wait for DNS propagation before asking the CA to validate. Increase this if validation fails due to slow DNS propagation."
	dnsWaitUCIOption            = "dns_wait"

	domainsAttribute            = "domains"
	domainsAttributeDescription = "List of domain names to include in the certificate. The first entry becomes the Common Name; subsequent entries become Subject Alternative Names. Wildcard domains (e.g. \"*.example.com\") require DNS validation."
	domainsUCIOption            = "domains"

	enabledAttribute            = "enabled"
	enabledAttributeDescription = "Enable or disable management of this certificate."
	enabledUCIOption            = "enabled"

	keyTypeAttribute            = "key_type"
	keyTypeAttributeDescription = "Key algorithm and size. One of \"rsa2048\", \"rsa3072\", \"rsa4096\", \"ec256\", or \"ec384\"."
	keyTypeUCIOption            = "key_type"

	listenPortAttribute            = "listen_port"
	listenPortAttributeDescription = "TCP port for the temporary ACME challenge HTTP server. Used by standalone and ALPN validation methods. Defaults to 80 (standalone/webroot) or 443 (ALPN)."
	listenPortUCIOption            = "listen_port"

	schemaDescription = "ACME-managed certificate configuration (UCI: /etc/config/acme, section type: cert). Each section manages one certificate with one or more domain names."

	stagingAttribute            = "staging"
	stagingAttributeDescription = "Use the Let's Encrypt staging environment. Certificates issued will not be trusted by browsers but allow testing without rate-limit concerns."
	stagingUCIOption            = "staging"

	uciConfig = "acme"
	uciType   = "cert"

	validationMethodAttribute            = "validation_method"
	validationMethodAttributeDescription = "ACME challenge type: \"standalone\" (built-in HTTP server), \"webroot\" (existing web server), \"dns\" (DNS-01 via provider API), or \"alpn\" (TLS-ALPN-01 on port 443)."
	validationMethodUCIOption            = "validation_method"

	webrootAttribute            = "webroot"
	webrootAttributeDescription = "Filesystem path served by the existing web server where ACME challenge files are written. Used when validation_method is \"webroot\". Default: \"/var/run/acme/challenge/\"."
	webrootUCIOption            = "webroot"
)

var (
	acmeServerSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       acmeServerAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetAcmeServer, acmeServerAttribute, acmeServerUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetAcmeServer, acmeServerAttribute, acmeServerUCIOption),
	}

	caliasSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       caliasAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetCalias, caliasAttribute, caliasUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetCalias, caliasAttribute, caliasUCIOption),
		Validators:        []validator.String{stringvalidator.LengthAtLeast(1)},
	}

	certProfileSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       certProfileAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetCertProfile, certProfileAttribute, certProfileUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetCertProfile, certProfileAttribute, certProfileUCIOption),
		Validators:        []validator.String{stringvalidator.LengthAtLeast(1)},
	}

	credentialsSchemaAttribute = lucirpcglue.ListStringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       credentialsAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionListString(modelSetCredentials, credentialsAttribute, credentialsUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		Sensitive:         true,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionListString(modelGetCredentials, credentialsAttribute, credentialsUCIOption),
	}

	daliasSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       daliasAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetDalias, daliasAttribute, daliasUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetDalias, daliasAttribute, daliasUCIOption),
		Validators:        []validator.String{stringvalidator.LengthAtLeast(1)},
	}

	daysSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       daysAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetDays, daysAttribute, daysUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetDays, daysAttribute, daysUCIOption),
		Validators:        []validator.Int64{int64validator.AtLeast(1)},
	}

	dnsSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       dnsAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetDNS, dnsAttribute, dnsUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetDNS, dnsAttribute, dnsUCIOption),
		Validators:        []validator.String{stringvalidator.LengthAtLeast(1)},
	}

	dnsWaitSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       dnsWaitAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetDNSWait, dnsWaitAttribute, dnsWaitUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetDNSWait, dnsWaitAttribute, dnsWaitUCIOption),
		Validators:        []validator.Int64{int64validator.AtLeast(0)},
	}

	domainsSchemaAttribute = lucirpcglue.ListStringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       domainsAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionListString(modelSetDomains, domainsAttribute, domainsUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionListString(modelGetDomains, domainsAttribute, domainsUCIOption),
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
			listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
		},
	}

	enabledSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       enabledAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetEnabled, enabledAttribute, enabledUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetEnabled, enabledAttribute, enabledUCIOption),
	}

	keyTypeSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       keyTypeAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetKeyType, keyTypeAttribute, keyTypeUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetKeyType, keyTypeAttribute, keyTypeUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf("rsa2048", "rsa3072", "rsa4096", "ec256", "ec384"),
		},
	}

	listenPortSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       listenPortAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetListenPort, listenPortAttribute, listenPortUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetListenPort, listenPortAttribute, listenPortUCIOption),
		Validators:        []validator.Int64{int64validator.Between(1, 65535)},
	}

	stagingSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       stagingAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetStaging, stagingAttribute, stagingUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetStaging, stagingAttribute, stagingUCIOption),
	}

	validationMethodSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       validationMethodAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetValidationMethod, validationMethodAttribute, validationMethodUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetValidationMethod, validationMethodAttribute, validationMethodUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf("standalone", "webroot", "dns", "alpn"),
		},
	}

	webrootSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       webrootAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetWebroot, webrootAttribute, webrootUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetWebroot, webrootAttribute, webrootUCIOption),
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		acmeServerAttribute:       acmeServerSchemaAttribute,
		caliasAttribute:           caliasSchemaAttribute,
		certProfileAttribute:      certProfileSchemaAttribute,
		credentialsAttribute:      credentialsSchemaAttribute,
		daliasAttribute:           daliasSchemaAttribute,
		daysAttribute:             daysSchemaAttribute,
		dnsAttribute:              dnsSchemaAttribute,
		dnsWaitAttribute:          dnsWaitSchemaAttribute,
		domainsAttribute:          domainsSchemaAttribute,
		enabledAttribute:          enabledSchemaAttribute,
		keyTypeAttribute:          keyTypeSchemaAttribute,
		listenPortAttribute:       listenPortSchemaAttribute,
		lucirpcglue.IdAttribute:   lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
		stagingAttribute:          stagingSchemaAttribute,
		validationMethodAttribute: validationMethodSchemaAttribute,
		webrootAttribute:          webrootSchemaAttribute,
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
	AcmeServer       types.String `tfsdk:"acme_server"`
	Calias           types.String `tfsdk:"calias"`
	CertProfile      types.String `tfsdk:"cert_profile"`
	Credentials      types.List   `tfsdk:"credentials"`
	Dalias           types.String `tfsdk:"dalias"`
	Days             types.Int64  `tfsdk:"days"`
	DNS              types.String `tfsdk:"dns"`
	DNSWait          types.Int64  `tfsdk:"dns_wait"`
	Domains          types.List   `tfsdk:"domains"`
	Enabled          types.Bool   `tfsdk:"enabled"`
	Id               types.String `tfsdk:"id"`
	KeyType          types.String `tfsdk:"key_type"`
	ListenPort       types.Int64  `tfsdk:"listen_port"`
	Staging          types.Bool   `tfsdk:"staging"`
	ValidationMethod types.String `tfsdk:"validation_method"`
	Webroot          types.String `tfsdk:"webroot"`
}

func modelGetAcmeServer(m model) types.String       { return m.AcmeServer }
func modelGetCalias(m model) types.String           { return m.Calias }
func modelGetCertProfile(m model) types.String      { return m.CertProfile }
func modelGetCredentials(m model) types.List        { return m.Credentials }
func modelGetDalias(m model) types.String           { return m.Dalias }
func modelGetDays(m model) types.Int64              { return m.Days }
func modelGetDNS(m model) types.String              { return m.DNS }
func modelGetDNSWait(m model) types.Int64           { return m.DNSWait }
func modelGetDomains(m model) types.List            { return m.Domains }
func modelGetEnabled(m model) types.Bool            { return m.Enabled }
func modelGetId(m model) types.String               { return m.Id }
func modelGetKeyType(m model) types.String          { return m.KeyType }
func modelGetListenPort(m model) types.Int64        { return m.ListenPort }
func modelGetStaging(m model) types.Bool            { return m.Staging }
func modelGetValidationMethod(m model) types.String { return m.ValidationMethod }
func modelGetWebroot(m model) types.String          { return m.Webroot }

func modelSetAcmeServer(m *model, value types.String)       { m.AcmeServer = value }
func modelSetCalias(m *model, value types.String)           { m.Calias = value }
func modelSetCertProfile(m *model, value types.String)      { m.CertProfile = value }
func modelSetCredentials(m *model, value types.List)        { m.Credentials = value }
func modelSetDalias(m *model, value types.String)           { m.Dalias = value }
func modelSetDays(m *model, value types.Int64)              { m.Days = value }
func modelSetDNS(m *model, value types.String)              { m.DNS = value }
func modelSetDNSWait(m *model, value types.Int64)           { m.DNSWait = value }
func modelSetDomains(m *model, value types.List)            { m.Domains = value }
func modelSetEnabled(m *model, value types.Bool)            { m.Enabled = value }
func modelSetId(m *model, value types.String)               { m.Id = value }
func modelSetKeyType(m *model, value types.String)          { m.KeyType = value }
func modelSetListenPort(m *model, value types.Int64)        { m.ListenPort = value }
func modelSetStaging(m *model, value types.Bool)            { m.Staging = value }
func modelSetValidationMethod(m *model, value types.String) { m.ValidationMethod = value }
func modelSetWebroot(m *model, value types.String)          { m.Webroot = value }
