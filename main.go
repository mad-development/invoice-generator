package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type Invoice struct {
	date     time.Time
	value    float64
	currency string
}

func (i *Invoice) GeneratePDF() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Invoice")
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Date: %s", i.date.Format("2006-01-02")))
	pdf.Ln(5)
	pdf.Cell(40, 10, fmt.Sprintf("Value: %.2f %s", i.value, i.currency))
	pdf.Ln(5)

	err := pdf.OutputFileAndClose("invoice.pdf")
	if err != nil {
		panic(err)
	}
}

func main() {
	day := flag.Int("day", 0, "Invoice day")
	month := flag.Int("month", 0, "Invoice month")
	year := flag.Int("year", 0, "Invoice year")
	value := flag.Float64("value", 0.0, "Invoice value")

	flag.Parse()

	if *day == 0 || *month == 0 || *year == 0 || *value == 0.0 {
		fmt.Println("Please provide all required parameters.")
		return
	}

	date := time.Date(*year, time.Month(*month), *day, 0, 0, 0, 0, time.UTC)

	invoice := &Invoice{
		date:     date,
		value:    *value,
		currency: "USD",
	}

	invoice.GeneratePDF()
}
