package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/maaslalani/invoice/invoice"
	"gopkg.in/yaml.v3"

	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

var (
	importPath string
	output     string
	localeName = invoice.English.Locale
	file       = invoice.Invoice{}
)

func init() {
	viper.AutomaticEnv()

	rootCmd.PersistentFlags().StringVar(&localeName, "locale", localeName, "")
	rootCmd.PersistentFlags().StringVar(&importPath, "import", "", "import invoice data (.json/.yaml)")
	rootCmd.PersistentFlags().StringVar(&file.Id, "id", "", "")
	rootCmd.PersistentFlags().StringVar(&file.Title, "title", "", "")

	rootCmd.PersistentFlags().Float64SliceVarP(&file.Rates, "rate", "r", []float64{}, "")
	rootCmd.PersistentFlags().Float64SliceVarP(&file.Quantities, "quantity", "q", []float64{}, "")
	rootCmd.PersistentFlags().StringSliceVarP(&file.Items, "item", "i", []string{}, "")

	rootCmd.PersistentFlags().StringVarP(&file.Logo, "logo", "l", "", "")
	rootCmd.PersistentFlags().StringVarP(&file.From, "from", "f", "", "")
	rootCmd.PersistentFlags().StringVarP(&file.To, "to", "t", "", "")
	rootCmd.PersistentFlags().StringVar(&file.Date, "date", "", "")
	rootCmd.PersistentFlags().StringVar(&file.Due, "due", "", "")

	rootCmd.PersistentFlags().Float64Var(&file.Tax, "tax", 0, "")
	rootCmd.PersistentFlags().Float64VarP(&file.Discount, "discount", "d", 0, "")
	rootCmd.PersistentFlags().StringVarP(&file.Currency, "currency", "c", "", "")

	rootCmd.PersistentFlags().StringVarP(&file.Note, "note", "n", "", "")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "invoice.pdf", "")

	flag.Parse()
}

var rootCmd = &cobra.Command{
	Use:   "invoice",
	Short: "Invoice generates invoices from the command line.",
	Long:  `Invoice generates invoices from the command line.`,
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate an invoice",
	Long:  `Generate an invoice`,
	RunE: func(cmd *cobra.Command, args []string) error {
		locale, err := invoice.GetLocale(localeName)
		if err != nil {
			return err
		}

		defaultInvoice := invoice.DefaultInvoice(locale)
		output, err := invoice.Generate(file, defaultInvoice, importPath, locale, output)
		if err != nil {
			return err
		}
		if output == nil {
			return errors.New("unknown error occurred")
		}

		fmt.Printf("Generated %s\n", *output)

		return nil
	},
}

func mergedInvoice() (*invoice.Invoice, error) {
	locale, err := invoice.GetLocale(localeName)
	if err != nil {
		return nil, err
	}

	defaultInvoice := invoice.DefaultInvoice(locale)
	return invoice.MergeInvoices(file, defaultInvoice, nil)
}

var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "Convert invoice to json",
	Long:  "Convert invoice to json",
	RunE: func(cmd *cobra.Command, args []string) error {
		mergedInvoice, err := mergedInvoice()
		if err != nil {
			return err
		}

		bytes, err := json.Marshal(mergedInvoice)
		if err != nil {
			return err
		}

		if !strings.HasSuffix(output, ".json") {
			output = output + ".json"
		}
		return os.WriteFile(output, bytes, 0666)
	},
}

var yamlCmd = &cobra.Command{
	Use:   "yaml",
	Short: "Convert invoice to yaml",
	Long:  "Convert invoice to yaml",
	RunE: func(cmd *cobra.Command, args []string) error {
		mergedInvoice, err := mergedInvoice()
		if err != nil {
			return err
		}

		bytes, err := yaml.Marshal(mergedInvoice)
		if err != nil {
			return err
		}

		if !strings.HasSuffix(output, ".yml") && !strings.HasSuffix(output, ".yaml") {
			output = output + ".yml"
		}
		return os.WriteFile(output, bytes, 0666)
	},
}

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Show docs",
	Long:  "Generate and render docs for the invoice generate command",
	RunE: func(cmd *cobra.Command, args []string) error {
		var buf = new(strings.Builder)
		err := doc.GenMarkdown(generateCmd, buf)
		if err != nil {
			return err
		}

		var md = buf.String()
		md, _, _ = strings.Cut(md, "\n### SEE ALSO")

		out, err := glamour.Render(md, "dark")
		if err != nil {
			return err
		}

		println(out)
		return nil
	},
}

func main() {
	rootCmd.AddCommand(generateCmd, jsonCmd, yamlCmd, docsCmd)
	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		os.Exit(1)
	}
}
