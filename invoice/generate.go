package invoice

import (
	"strings"

	"dario.cat/mergo"
	"github.com/signintech/gopdf"
)

func MergeInvoices(invoice, defaultInvoice Invoice, importedInvoice *Invoice) (*Invoice, error) {
	finalInvoice := defaultInvoice

	if importedInvoice != nil {
		err := mergo.Merge(&finalInvoice, *importedInvoice, mergo.WithOverride)
		if err != nil {
			return nil, err
		}
	}

	err := mergo.Merge(&finalInvoice, invoice, mergo.WithOverride)
	if err != nil {
		return nil, err
	}

	return &finalInvoice, nil
}

func Generate(invoice Invoice, defaultInvoice Invoice, importPath string, locale Locale, output string) (*string, error) {
	var importedInvoice *Invoice

	if importPath != "" {
		err := importStruct(importPath, &importedInvoice)
		if err != nil {
			return nil, err
		}
	}
	newInvoice, err := MergeInvoices(invoice, defaultInvoice, importedInvoice)
	if err != nil {
		return nil, err
	}
	invoice = *newInvoice

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{
		PageSize: *gopdf.PageSizeA4,
	})
	pdf.SetMargins(40, 40, 40, 40)
	pdf.AddPage()
	err = pdf.AddTTFFontData("Inter", interFont)
	if err != nil {
		return nil, err
	}

	err = pdf.AddTTFFontData("Inter-Bold", interBoldFont)
	if err != nil {
		return nil, err
	}

	writeLogo(&pdf, invoice.Logo, invoice.From)
	writeTitle(&pdf, invoice.Title, invoice.Id, invoice.Date)
	writeBillTo(&pdf, locale, invoice.To)
	writeHeaderRow(&pdf, locale)

	subtotal := 0.0
	for i := range invoice.Items {
		q := 1.0
		if len(invoice.Quantities) > i {
			q = invoice.Quantities[i]
		}

		r := 0.0
		if len(invoice.Rates) > i {
			r = invoice.Rates[i]
		}

		if pdf.GetY() >= 650 {
			writeFooter(&pdf, invoice.Id)
			pdf.AddPage()
		}

		writeRow(&pdf, invoice.Items[i], q, r, invoice.Currency)
		subtotal += q * r
	}

	if invoice.Note != "" {
		writeNotes(&pdf, locale, invoice.Note)
	}

	writeTotals(&pdf, locale, subtotal, subtotal*invoice.Tax, subtotal*invoice.Discount, invoice.Currency)
	if invoice.Due != "" {
		writeDueDate(&pdf, locale, invoice.Due)
	}

	writeFooter(&pdf, invoice.Id)

	output = strings.TrimSuffix(output, ".pdf") + ".pdf"
	err = pdf.WritePdf(output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}
