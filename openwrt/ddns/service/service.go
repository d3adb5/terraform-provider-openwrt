package service

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	caCertAttribute            = "cacert"
	caCertAttributeDescription = "Path to the CA certificate file or directory for TLS verification."
	caCertUCIOption            = "cacert"

	checkIntervalAttribute            = "check_interval"
	checkIntervalAttributeDescription = "How often to check whether an update is needed, expressed in check_unit units."
	checkIntervalUCIOption            = "check_interval"

	checkUnitAttribute            = "check_unit"
	checkUnitAttributeDescription = `Time unit for check_interval. One of "seconds", "minutes", "hours", or "days".`
	checkUnitUCIOption            = "check_unit"

	domainAttribute            = "domain"
	domainAttributeDescription = "The fully qualified domain name (FQDN) to keep updated."
	domainUCIOption            = "domain"

	enabledAttribute            = "enabled"
	enabledAttributeDescription = "Enable or disable this DDNS service entry."
	enabledUCIOption            = "enabled"

	forceIntervalAttribute            = "force_interval"
	forceIntervalAttributeDescription = "Force an update even if the IP address has not changed, expressed in force_unit units."
	forceIntervalUCIOption            = "force_interval"

	forceUnitAttribute            = "force_unit"
	forceUnitAttributeDescription = `Time unit for force_interval. One of "seconds", "minutes", "hours", or "days".`
	forceUnitUCIOption            = "force_unit"

	interfaceAttribute            = "interface"
	interfaceAttributeDescription = `Network interface to monitor for IP changes (e.g. "wan"). An update is triggered when this interface comes up or its address changes.`
	interfaceUCIOption            = "interface"

	ipNetworkAttribute            = "ip_network"
	ipNetworkAttributeDescription = `Network interface to read the current IP address from. Used when ip_source is "network".`
	ipNetworkUCIOption            = "ip_network"

	ipScriptAttribute            = "ip_script"
	ipScriptAttributeDescription = `Path to a script that outputs the current IP address on stdout. Used when ip_source is "script".`
	ipScriptUCIOption            = "ip_script"

	ipSourceAttribute            = "ip_source"
	ipSourceAttributeDescription = `How to detect the current public IP address: "network" (read from a local interface), "web" (query an external URL), or "script" (run a custom script).`
	ipSourceUCIOption            = "ip_source"

	ipUrlAttribute            = "ip_url"
	ipUrlAttributeDescription = `URL of an external service that returns the current public IP address (e.g. "https://api4.my-ip.io/ip"). Used when ip_source is "web".`
	ipUrlUCIOption            = "ip_url"

	lookupHostAttribute            = "lookup_host"
	lookupHostAttributeDescription = "Hostname used for DNS lookup to verify the currently registered IP. Defaults to the value of domain when unset."
	lookupHostUCIOption            = "lookup_host"

	passwordAttribute            = "password"
	passwordAttributeDescription = "Password or API token for the DDNS provider."
	passwordUCIOption            = "password"

	retryCountAttribute            = "retry_count"
	retryCountAttributeDescription = "Number of times to retry a failed update before giving up and waiting for the next check cycle."
	retryCountUCIOption            = "retry_count"

	retryIntervalAttribute            = "retry_interval"
	retryIntervalAttributeDescription = "How long to wait between retries, expressed in retry_unit units."
	retryIntervalUCIOption            = "retry_interval"

	retryUnitAttribute            = "retry_unit"
	retryUnitAttributeDescription = `Time unit for retry_interval. One of "seconds", "minutes", "hours", or "days".`
	retryUnitUCIOption            = "retry_unit"

	schemaDescription = "Dynamic DNS (DDNS) service configuration. Each section configures one DDNS record to keep updated via ddns-scripts."

	serviceNameAttribute            = "service_name"
	serviceNameAttributeDescription = `DDNS provider to use (e.g. "cloudflare.com", "no-ip.com", "dyndns.org"). Use "--custom--" together with update_url for a provider not in the built-in list.`
	serviceNameUCIOption            = "service_name"

	uciConfig = "ddns"
	uciType   = "service"

	updateScriptAttribute            = "update_script"
	updateScriptAttributeDescription = "Path to a custom script that performs the DDNS update. Used instead of the built-in provider logic."
	updateScriptUCIOption            = "update_script"

	updateUrlAttribute            = "update_url"
	updateUrlAttributeDescription = `Custom update URL template used with service_name "--custom--". Placeholders like [USERNAME], [PASSWORD], [DOMAIN], and [IP] are substituted at runtime.`
	updateUrlUCIOption            = "update_url"

	useHttpsAttribute            = "use_https"
	useHttpsAttributeDescription = "Use HTTPS when communicating with the DDNS provider."
	useHttpsUCIOption            = "use_https"

	useIpv6Attribute            = "use_ipv6"
	useIpv6AttributeDescription = "Update the IPv6 AAAA record instead of the IPv4 A record."
	useIpv6UCIOption            = "use_ipv6"

	useLogfileAttribute            = "use_logfile"
	useLogfileAttributeDescription = "Write update activity to a dedicated log file."
	useLogfileUCIOption            = "use_logfile"

	useSyslogAttribute            = "use_syslog"
	useSyslogAttributeDescription = "Syslog verbosity level. 0 = disabled, 1–4 = increasing verbosity."
	useSyslogUCIOption            = "use_syslog"

	usernameAttribute            = "username"
	usernameAttributeDescription = "Username or account identifier for the DDNS provider."
	usernameUCIOption            = "username"
)

var (
	timeUnitValidators = []validator.String{
		stringvalidator.OneOf("seconds", "minutes", "hours", "days"),
	}

	caCertSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       caCertAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetCACert, caCertAttribute, caCertUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetCACert, caCertAttribute, caCertUCIOption),
	}

	checkIntervalSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       checkIntervalAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetCheckInterval, checkIntervalAttribute, checkIntervalUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetCheckInterval, checkIntervalAttribute, checkIntervalUCIOption),
		Validators:        []validator.Int64{int64validator.AtLeast(1)},
	}

	checkUnitSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       checkUnitAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetCheckUnit, checkUnitAttribute, checkUnitUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetCheckUnit, checkUnitAttribute, checkUnitUCIOption),
		Validators:        timeUnitValidators,
	}

	domainSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       domainAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetDomain, domainAttribute, domainUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetDomain, domainAttribute, domainUCIOption),
		Validators:        []validator.String{stringvalidator.LengthAtLeast(1)},
	}

	enabledSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       enabledAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetEnabled, enabledAttribute, enabledUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetEnabled, enabledAttribute, enabledUCIOption),
	}

	forceIntervalSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       forceIntervalAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetForceInterval, forceIntervalAttribute, forceIntervalUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetForceInterval, forceIntervalAttribute, forceIntervalUCIOption),
		Validators:        []validator.Int64{int64validator.AtLeast(1)},
	}

	forceUnitSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       forceUnitAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetForceUnit, forceUnitAttribute, forceUnitUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetForceUnit, forceUnitAttribute, forceUnitUCIOption),
		Validators:        timeUnitValidators,
	}

	interfaceSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       interfaceAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetInterface, interfaceAttribute, interfaceUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetInterface, interfaceAttribute, interfaceUCIOption),
		Validators:        []validator.String{stringvalidator.LengthAtLeast(1)},
	}

	ipNetworkSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       ipNetworkAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetIPNetwork, ipNetworkAttribute, ipNetworkUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetIPNetwork, ipNetworkAttribute, ipNetworkUCIOption),
		Validators:        []validator.String{stringvalidator.LengthAtLeast(1)},
	}

	ipScriptSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       ipScriptAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetIPScript, ipScriptAttribute, ipScriptUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetIPScript, ipScriptAttribute, ipScriptUCIOption),
	}

	ipSourceSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       ipSourceAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetIPSource, ipSourceAttribute, ipSourceUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetIPSource, ipSourceAttribute, ipSourceUCIOption),
		Validators:        []validator.String{stringvalidator.OneOf("network", "web", "script")},
	}

	ipUrlSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       ipUrlAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetIPUrl, ipUrlAttribute, ipUrlUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetIPUrl, ipUrlAttribute, ipUrlUCIOption),
	}

	lookupHostSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       lookupHostAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetLookupHost, lookupHostAttribute, lookupHostUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetLookupHost, lookupHostAttribute, lookupHostUCIOption),
		Validators:        []validator.String{stringvalidator.LengthAtLeast(1)},
	}

	passwordSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       passwordAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetPassword, passwordAttribute, passwordUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		Sensitive:         true,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetPassword, passwordAttribute, passwordUCIOption),
	}

	retryCountSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       retryCountAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetRetryCount, retryCountAttribute, retryCountUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetRetryCount, retryCountAttribute, retryCountUCIOption),
		Validators:        []validator.Int64{int64validator.AtLeast(0)},
	}

	retryIntervalSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       retryIntervalAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetRetryInterval, retryIntervalAttribute, retryIntervalUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetRetryInterval, retryIntervalAttribute, retryIntervalUCIOption),
		Validators:        []validator.Int64{int64validator.AtLeast(1)},
	}

	retryUnitSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       retryUnitAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetRetryUnit, retryUnitAttribute, retryUnitUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetRetryUnit, retryUnitAttribute, retryUnitUCIOption),
		Validators:        timeUnitValidators,
	}

	serviceNameSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       serviceNameAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetServiceName, serviceNameAttribute, serviceNameUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetServiceName, serviceNameAttribute, serviceNameUCIOption),
		Validators:        []validator.String{stringvalidator.LengthAtLeast(1)},
	}

	updateScriptSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       updateScriptAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetUpdateScript, updateScriptAttribute, updateScriptUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetUpdateScript, updateScriptAttribute, updateScriptUCIOption),
	}

	updateUrlSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       updateUrlAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetUpdateURL, updateUrlAttribute, updateUrlUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetUpdateURL, updateUrlAttribute, updateUrlUCIOption),
	}

	useHttpsSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       useHttpsAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetUseHTTPS, useHttpsAttribute, useHttpsUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetUseHTTPS, useHttpsAttribute, useHttpsUCIOption),
	}

	useIpv6SchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       useIpv6AttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetUseIPv6, useIpv6Attribute, useIpv6UCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetUseIPv6, useIpv6Attribute, useIpv6UCIOption),
	}

	useLogfileSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       useLogfileAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetUseLogfile, useLogfileAttribute, useLogfileUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetUseLogfile, useLogfileAttribute, useLogfileUCIOption),
	}

	useSyslogSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       useSyslogAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetUseSyslog, useSyslogAttribute, useSyslogUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetUseSyslog, useSyslogAttribute, useSyslogUCIOption),
		Validators:        []validator.Int64{int64validator.Between(0, 4)},
	}

	usernameSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       usernameAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetUsername, usernameAttribute, usernameUCIOption),
		ResourceExistence: lucirpcglue.Optional,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetUsername, usernameAttribute, usernameUCIOption),
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		caCertAttribute:         caCertSchemaAttribute,
		checkIntervalAttribute:  checkIntervalSchemaAttribute,
		checkUnitAttribute:      checkUnitSchemaAttribute,
		domainAttribute:         domainSchemaAttribute,
		enabledAttribute:        enabledSchemaAttribute,
		forceIntervalAttribute:  forceIntervalSchemaAttribute,
		forceUnitAttribute:      forceUnitSchemaAttribute,
		lucirpcglue.IdAttribute: lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
		interfaceAttribute:      interfaceSchemaAttribute,
		ipNetworkAttribute:      ipNetworkSchemaAttribute,
		ipScriptAttribute:       ipScriptSchemaAttribute,
		ipSourceAttribute:       ipSourceSchemaAttribute,
		ipUrlAttribute:          ipUrlSchemaAttribute,
		lookupHostAttribute:     lookupHostSchemaAttribute,
		passwordAttribute:       passwordSchemaAttribute,
		retryCountAttribute:     retryCountSchemaAttribute,
		retryIntervalAttribute:  retryIntervalSchemaAttribute,
		retryUnitAttribute:      retryUnitSchemaAttribute,
		serviceNameAttribute:    serviceNameSchemaAttribute,
		updateScriptAttribute:   updateScriptSchemaAttribute,
		updateUrlAttribute:      updateUrlSchemaAttribute,
		useHttpsAttribute:       useHttpsSchemaAttribute,
		useIpv6Attribute:        useIpv6SchemaAttribute,
		useLogfileAttribute:     useLogfileSchemaAttribute,
		useSyslogAttribute:      useSyslogSchemaAttribute,
		usernameAttribute:       usernameSchemaAttribute,
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
	CACert        types.String `tfsdk:"cacert"`
	CheckInterval types.Int64  `tfsdk:"check_interval"`
	CheckUnit     types.String `tfsdk:"check_unit"`
	Domain        types.String `tfsdk:"domain"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	ForceInterval types.Int64  `tfsdk:"force_interval"`
	ForceUnit     types.String `tfsdk:"force_unit"`
	Id            types.String `tfsdk:"id"`
	Interface     types.String `tfsdk:"interface"`
	IPNetwork     types.String `tfsdk:"ip_network"`
	IPScript      types.String `tfsdk:"ip_script"`
	IPSource      types.String `tfsdk:"ip_source"`
	IPUrl         types.String `tfsdk:"ip_url"`
	LookupHost    types.String `tfsdk:"lookup_host"`
	Password      types.String `tfsdk:"password"`
	RetryCount    types.Int64  `tfsdk:"retry_count"`
	RetryInterval types.Int64  `tfsdk:"retry_interval"`
	RetryUnit     types.String `tfsdk:"retry_unit"`
	ServiceName   types.String `tfsdk:"service_name"`
	UpdateScript  types.String `tfsdk:"update_script"`
	UpdateURL     types.String `tfsdk:"update_url"`
	UseHTTPS      types.Bool   `tfsdk:"use_https"`
	UseIPv6       types.Bool   `tfsdk:"use_ipv6"`
	UseLogfile    types.Bool   `tfsdk:"use_logfile"`
	UseSyslog     types.Int64  `tfsdk:"use_syslog"`
	Username      types.String `tfsdk:"username"`
}

func modelGetCACert(m model) types.String        { return m.CACert }
func modelGetCheckInterval(m model) types.Int64  { return m.CheckInterval }
func modelGetCheckUnit(m model) types.String     { return m.CheckUnit }
func modelGetDomain(m model) types.String        { return m.Domain }
func modelGetEnabled(m model) types.Bool         { return m.Enabled }
func modelGetForceInterval(m model) types.Int64  { return m.ForceInterval }
func modelGetForceUnit(m model) types.String     { return m.ForceUnit }
func modelGetId(m model) types.String            { return m.Id }
func modelGetInterface(m model) types.String     { return m.Interface }
func modelGetIPNetwork(m model) types.String     { return m.IPNetwork }
func modelGetIPScript(m model) types.String      { return m.IPScript }
func modelGetIPSource(m model) types.String      { return m.IPSource }
func modelGetIPUrl(m model) types.String         { return m.IPUrl }
func modelGetLookupHost(m model) types.String    { return m.LookupHost }
func modelGetPassword(m model) types.String      { return m.Password }
func modelGetRetryCount(m model) types.Int64     { return m.RetryCount }
func modelGetRetryInterval(m model) types.Int64  { return m.RetryInterval }
func modelGetRetryUnit(m model) types.String     { return m.RetryUnit }
func modelGetServiceName(m model) types.String   { return m.ServiceName }
func modelGetUpdateScript(m model) types.String  { return m.UpdateScript }
func modelGetUpdateURL(m model) types.String     { return m.UpdateURL }
func modelGetUseHTTPS(m model) types.Bool        { return m.UseHTTPS }
func modelGetUseIPv6(m model) types.Bool         { return m.UseIPv6 }
func modelGetUseLogfile(m model) types.Bool      { return m.UseLogfile }
func modelGetUseSyslog(m model) types.Int64      { return m.UseSyslog }
func modelGetUsername(m model) types.String      { return m.Username }

func modelSetCACert(m *model, value types.String)        { m.CACert = value }
func modelSetCheckInterval(m *model, value types.Int64)  { m.CheckInterval = value }
func modelSetCheckUnit(m *model, value types.String)     { m.CheckUnit = value }
func modelSetDomain(m *model, value types.String)        { m.Domain = value }
func modelSetEnabled(m *model, value types.Bool)         { m.Enabled = value }
func modelSetForceInterval(m *model, value types.Int64)  { m.ForceInterval = value }
func modelSetForceUnit(m *model, value types.String)     { m.ForceUnit = value }
func modelSetId(m *model, value types.String)            { m.Id = value }
func modelSetInterface(m *model, value types.String)     { m.Interface = value }
func modelSetIPNetwork(m *model, value types.String)     { m.IPNetwork = value }
func modelSetIPScript(m *model, value types.String)      { m.IPScript = value }
func modelSetIPSource(m *model, value types.String)      { m.IPSource = value }
func modelSetIPUrl(m *model, value types.String)         { m.IPUrl = value }
func modelSetLookupHost(m *model, value types.String)    { m.LookupHost = value }
func modelSetPassword(m *model, value types.String)      { m.Password = value }
func modelSetRetryCount(m *model, value types.Int64)     { m.RetryCount = value }
func modelSetRetryInterval(m *model, value types.Int64)  { m.RetryInterval = value }
func modelSetRetryUnit(m *model, value types.String)     { m.RetryUnit = value }
func modelSetServiceName(m *model, value types.String)   { m.ServiceName = value }
func modelSetUpdateScript(m *model, value types.String)  { m.UpdateScript = value }
func modelSetUpdateURL(m *model, value types.String)     { m.UpdateURL = value }
func modelSetUseHTTPS(m *model, value types.Bool)        { m.UseHTTPS = value }
func modelSetUseIPv6(m *model, value types.Bool)         { m.UseIPv6 = value }
func modelSetUseLogfile(m *model, value types.Bool)      { m.UseLogfile = value }
func modelSetUseSyslog(m *model, value types.Int64)      { m.UseSyslog = value }
func modelSetUsername(m *model, value types.String)      { m.Username = value }
