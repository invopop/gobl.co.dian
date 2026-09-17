# GOBL ➡️ Colombia DIAN

Colombia DIAN (Dirección de Impuestos y Aduanas Nacionales) e-invoicing addon for [GOBL](https://github.com/invopop/gobl).

Copyright [Invopop S.L.](https://invopop.com) 2026. Released publicly under the [GNU Affero General Public License v3.0](LICENSE). For commercial licenses please contact the [dev team at invopop](mailto:dev@invopop.com).

[![Lint](https://github.com/invopop/gobl.co.dian/actions/workflows/lint.yaml/badge.svg)](https://github.com/invopop/gobl.co.dian/actions/workflows/lint.yaml)
[![Test Go](https://github.com/invopop/gobl.co.dian/actions/workflows/test.yaml/badge.svg)](https://github.com/invopop/gobl.co.dian/actions/workflows/test.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/invopop/gobl.co.dian)](https://goreportcard.com/report/github.com/invopop/gobl.co.dian)
[![codecov](https://codecov.io/gh/invopop/gobl.co.dian/graph/badge.svg)](https://codecov.io/gh/invopop/gobl.co.dian)
[![GoDoc](https://godoc.org/github.com/invopop/gobl.co.dian?status.svg)](https://godoc.org/github.com/invopop/gobl.co.dian)
![Latest Tag](https://img.shields.io/github/v/tag/invopop/gobl.co.dian)

This module implements the Colombia DIAN e-invoicing requirements as a GOBL
tax addon (`co-dian-v2`), based on UBL 2.1 and the DIAN "Anexo Técnico de la
Factura Electrónica de Venta". It covers sales invoices, credit notes, and
debit notes, with validation rules registered under the `CO-DIAN` namespace.

The addon lives in the `addon/` subpackage and registers into GOBL's global
registry on import; the module root is reserved for future DIAN tooling. The
Colombian tax regime (`regimes/co`) stays in GOBL core.

## Usage

Add a blank import of the addon, then declare it on documents:

```go
import (
	_ "github.com/invopop/gobl.co.dian/addon"
)
```

```json
{
	"$schema": "https://gobl.org/draft-0/bill/invoice",
	"$regime": "CO",
	"$addons": ["co-dian-v2"]
}
```

> **Note**: `co-dian-v2` is an approved external addon key in GOBL core, but
> documents declaring it fail validation unless this module is imported.

## Extensions

| Key | Description |
| --- | --- |
| `co-dian-municipality` | DANE municipality code for party addresses. Required for Colombian parties. |
| `co-dian-fiscal-responsibility` | Party fiscal responsibility (`TaxLevelCode`). Required; defaults to `R-99-PN`. |
| `co-dian-tax-scheme` | Party tributo (`PartyTaxScheme/TaxScheme`): `01`, `04`, `ZA`, `ZZ`. Optional. |
| `co-dian-item-identification` | Product-coding standard for line item codes. Optional. |
| `co-dian-credit-code` | Correction concept, required on credit note `preceding` references. |
| `co-dian-debit-code` | Correction concept, required on debit note `preceding` references. |

Simplified (B2C) customers may be identified with DIAN document types via
`org.Identity` keys (`co-citizen-id`, `co-passport`, …), and corrections via
GOBL's `Correct` process copy the `dian-cude` stamp onto the `preceding`
reference.

## Development

The `go.mod` pins `github.com/invopop/gobl` to a branch pseudo-version until
the external-addon registration ships in a tagged release.

`examples/` holds sample documents with golden envelopes under `examples/out/`;
regenerate them with:

```sh
go test . -run TestExamples -update
```

## Sources

- [Anexo Técnico de la Factura Electrónica de Venta v1.9 (PDF)](https://www.dian.gov.co/impuestos/factura-electronica/Documents/Anexo-Tecnico-Factura-Electronica-de-Venta-vr-1-9.pdf)
- [Caja de Herramientas FE V1.9 (referenced code tables, Excel)](https://www.dian.gov.co/impuestos/factura-electronica/Documents/Caja-de-herramientas-FE-V1-9.zip)
