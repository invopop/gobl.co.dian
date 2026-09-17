package dian_test

import (
	"os"
	"testing"

	"github.com/invopop/gobl"
	dian "github.com/invopop/gobl.co.dian/addon"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func invoiceForCorrection(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Regime: tax.WithRegime("CO"),
		Addons: tax.WithAddons(dian.V2),
		Series: "TEST",
		Code:   "123",
		Tax: &bill.Tax{
			PricesInclude: tax.CategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: "CO",
				Code:    "9014586527",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: "CO",
				Code:    "8001345363",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "general",
					},
				},
				Discounts: []*bill.LineDiscount{
					{
						Reason:  "Testing",
						Percent: num.NewPercentage(10, 2),
					},
				},
			},
		},
	}
}

func TestInvoiceCorrect(t *testing.T) {
	i := invoiceForCorrection(t)
	err := i.Correct(bill.Credit)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing stamp")

	stamps := []*head.Stamp{
		{
			Provider: dian.StampCUDE,
			Value:    "FOOO",
		},
		{
			Provider: dian.StampQR, // not copied!
			Value:    "BARRRR",
		},
	}

	i = invoiceForCorrection(t)
	err = i.Correct(bill.Corrective, bill.WithStamps(stamps))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid correction type: corrective")

	i = invoiceForCorrection(t)
	err = i.Correct(
		bill.Credit,
		bill.WithStamps(stamps),
		bill.WithReason("test refund"),
		bill.WithExtension(dian.ExtKeyCreditCode, "2"),
	)
	require.NoError(t, err)
	assert.Equal(t, i.Type, bill.InvoiceTypeCreditNote)
	pre := i.Preceding[0]
	require.Len(t, pre.Stamps, 1)
	assert.Equal(t, pre.Stamps[0].Provider, dian.StampCUDE)
}

func TestEnvelopeCorrectWithStamps(t *testing.T) {
	data, err := os.ReadFile("../examples/out/simple.json")
	require.NoError(t, err)
	out, err := gobl.Parse(data)
	require.NoError(t, err)
	env, ok := out.(*gobl.Envelope)
	require.True(t, ok)
	env.Head.AddStamp(&head.Stamp{
		Provider: dian.StampCUDE,
		Value:    "1234567890",
	})

	_, err = env.Correct()
	assert.ErrorContains(t, err, "validation: missing correction type")

	e2, err := env.Correct(bill.Credit, bill.WithReason("test"))
	require.NoError(t, err)
	doc := e2.Extract().(*bill.Invoice)
	assert.Equal(t, "1234567890", doc.Preceding[0].Stamps[0].Value, "should copy stamps")
}
