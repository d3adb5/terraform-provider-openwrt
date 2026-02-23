# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```sh
make build                        # go build ./...
make test                         # build + docs check + unit + acceptance tests
go test ./...                     # unit tests only (fast, no Docker)
TF_ACC=1 go test -tags=acceptance.test ./...  # acceptance tests (requires Docker)

# Single package / single test
go test -v ./lucirpc -run TestClientCreateSection
TF_ACC=1 go test -tags=acceptance.test -v ./openwrt/firewall/rule -run TestAcceptance

make start-acceptance-test-server # spin up OpenWrt Docker container for acceptance tests
make clean                        # tear down Docker containers and cache

make docs                         # regenerate docs/
```

Docs are auto-generated (`go generate ./...` via tfplugindocs) and committed. CI verifies they are up to date — run `make docs` after any schema change.

## Architecture

This is a Terraform provider for OpenWrt devices. Communication happens through the LuCI JSON-RPC API.

### Layers

**`lucirpc/`** — Low-level RPC client. `client.go` handles authentication and the four UCI operations (`CreateSection`, `DeleteSection`, `GetSection`, `UpdateSection`, `CommitChanges`). `options.go` defines `Options` (`map[string]Option`) and the four concrete option types (`optionBoolean`, `optionInteger`, `optionListString`, `optionString`). Option types coerce between Go types — e.g. a single-element `optionListString` can return as a string, and an `optionString` can return as a one-element slice.

**`openwrt/internal/lucirpcglue/`** — Generic glue between the Terraform Plugin Framework and `lucirpc`. The key file is `attribute.go`, which defines `SchemaAttribute[Model, Request, Response]` — an interface implemented by `BoolSchemaAttribute`, `StringSchemaAttribute`, `Int64SchemaAttribute`, `ListStringSchemaAttribute`, and `SetStringSchemaAttribute`. Each attribute knows how to render itself into a Terraform schema, read from a UCI response into the model, and serialize from the model into a UCI upsert request. `resource.go` and `data_source.go` provide generic `NewResource` / `NewDataSource` constructors. `model.go` has `ReadModel` and `GenerateUpsertBody` that iterate over the attribute map.

**`openwrt/{dhcp,firewall,network,system,wireless}/`** — Domain-specific resource packages. Each sub-package (e.g. `openwrt/firewall/rule/`) is self-contained and follows an identical pattern:

```
rule.go        — constants (attribute names, UCI option names, descriptions), var block of
                 SchemaAttribute values, schemaAttributes map, NewResource/NewDataSource,
                 model struct, modelGet*/modelSet* funcs
```

`openwrt/provider.go` registers all resources and data sources and plumbs the LuCI client through `provider_data`.

### Adding a new resource

1. Create `openwrt/<domain>/<type>/` with a single `<type>.go` following the pattern above.
2. Define `uciConfig` and `uciType` constants for the UCI config file and section type.
3. For each attribute: three constants (`fooAttribute`, `fooAttributeDescription`, `fooUCIOption`), one `*SchemaAttribute` var, and getter/setter funcs on the model.
4. Use `ResourceExistence: lucirpcglue.Required` for required fields, `lucirpcglue.Optional` for user-controlled optional fields, and `lucirpcglue.NoValidation` only for attributes the provider itself computes (triggers `Computed+Optional` + `UseStateForUnknown`).
5. Register in `openwrt/provider.go`.
6. Run `make docs` to regenerate documentation.

### Key conventions

- **UCI option types**: OpenWrt can store a value as either `option` (string) or `list`. `optionString.AsListString()` returns a one-element slice; `optionListString.AsString()` returns the value only when the list has exactly one element. Use `ListStringSchemaAttribute` for anything OpenWrt may normalize into a multi-element list (e.g. `proto`).
- **Port ranges**: OpenWrt uses hyphen notation (`67-68`). The `portValidators` regex accepts both `[:-]`.
- **`NoValidation` vs `Optional`**: `NoValidation` = Computed+Optional (with `UseStateForUnknown`). Use it only when the provider or OpenWrt auto-generates the value. For user-controlled optional fields that OpenWrt never auto-populates, use `Optional` to avoid `(known after apply)` in plans.
- **Diagnostics**: All internal functions return `diag.Diagnostics`. In `ReadModel`, always append to `allDiagnostics` — returning the last-iteration `diagnostics` variable silently drops earlier errors.
