package invoice

import (
	"errors"
	"strings"
	"time"

	"github.com/klauspost/lctime"
)

var English = Locale{
	DefaultTitle:    "INVOICE",
	DefaultItem:     "Paper Cranes",
	DefaultFrom:     "Project Folded, Inc.",
	DefaultTo:       "Untitled Corporation, Inc.",
	DefaultCurrency: "USD",
	Locale:          "en_US",
	RateLabel:       "RATE",
	QuantityLabel:   "QTY",
	AmountLabel:     "AMOUNT",
	ItemLabel:       "ITEM",
	ToLabel:         "BILL TO",
	NotesLabel:      "Notes",
	SubtotalLabel:   "Subtotal",
	TaxLabel:        "Tax",
	DiscountLabel:   "Discount",
	TotalLabel:      "Total",
	DueLabel:        "Due Date",
}

var Dutch = Locale{
	DefaultTitle:    "Factuur",
	DefaultItem:     "Speelgoed brandweerwagen",
	DefaultFrom:     "Toet toet auto's B.V.",
	DefaultTo:       "Naamloos bedrijf B.V.",
	DefaultCurrency: "EUR",
	Locale:          "nl_NL",
	RateLabel:       "Stuksprijs",
	QuantityLabel:   "Aantal",
	AmountLabel:     "Totaal",
	ItemLabel:       "Omschrijving",
	ToLabel:         "",
	NotesLabel:      "Opmerkingen",
	SubtotalLabel:   "Subtotaal",
	TaxLabel:        "BTW",
	DiscountLabel:   "Korting",
	TotalLabel:      "Te voldoen",
	DueLabel:        "Voldoen voor",
}

type Locale struct {
	DefaultTitle    string
	DefaultItem     string
	DefaultFrom     string
	DefaultTo       string
	DefaultCurrency string
	Locale          string
	RateLabel       string
	QuantityLabel   string
	AmountLabel     string
	ItemLabel       string
	ToLabel         string
	NotesLabel      string
	SubtotalLabel   string
	TaxLabel        string
	DiscountLabel   string
	TotalLabel      string
	DueLabel        string
}

func DefaultLocale() Locale {
	return English
}

func GetLocale(name string) (Locale, error) {
	switch strings.ToLower(name) {
	case "english", "en_us", "en":
		return English, nil
	case "dutch", "nl_nl", "nl":
		return Dutch, nil
	}

	return DefaultLocale(), errors.New("unsupported locale")
}

type Invoice struct {
	Id    string `json:"id" yaml:"id"`
	Title string `json:"title" yaml:"title"`

	Logo string `json:"logo" yaml:"logo"`
	From string `json:"from" yaml:"from"`
	To   string `json:"to" yaml:"to"`
	Date string `json:"date" yaml:"date"`
	Due  string `json:"due" yaml:"due"`

	Items      []string  `json:"items" yaml:"items"`
	Quantities []float64 `json:"quantities" yaml:"quantities"`
	Rates      []float64 `json:"rates" yaml:"rates"`

	Tax      float64 `json:"tax" yaml:"tax"`
	Discount float64 `json:"discount" yaml:"discount"`
	Currency string  `json:"currency" yaml:"currency"`

	Note string `json:"note" yaml:"note"`
}

func DefaultInvoice(locale Locale) Invoice {
	now := time.Now()

	date, err := lctime.StrftimeLoc(locale.Locale, "%x", now)
	if err != nil {
		date = now.Format("2006/01/02")
	}

	dueDate := now.AddDate(0, 0, 14)
	due, err := lctime.StrftimeLoc(locale.Locale, "%x", dueDate)
	if err != nil {
		due = dueDate.Format("2006/01/02")
	}

	return Invoice{
		Id:         now.Format("20060102"),
		Title:      locale.DefaultTitle,
		Rates:      []float64{25},
		Quantities: []float64{2},
		Items:      []string{locale.DefaultItem},
		From:       locale.DefaultFrom,
		To:         locale.DefaultTo,
		Date:       date,
		Due:        due,
		Tax:        0,
		Discount:   0,
		Currency:   locale.DefaultCurrency,
	}
}
