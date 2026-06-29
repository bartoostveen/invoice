package internal

var currencySymbols = map[string]string{
	"USD": "$",
	"EUR": "€",
	"GBP": "£",
	"JPY": "¥",
	"CNY": "¥",
	"INR": "₹",
	"RUB": "₽",
	"KRW": "₩",
	"BRL": "R$",
	"SGD": "SGD$",
}

func getCurrencySymbol(currency string) string {
	symbol, ok := currencySymbols[currency]
	if !ok {
		return currency
	}
	return symbol
}
