package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type Invoice struct {
	date      time.Time
	value     float64
	currency  string
	client    string
	invoiceID string
	services  map[string]float64
}

func (i *Invoice) GeneratePDF() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)

	// Adiciona um título
	pdf.CellFormat(190, 10, "Invoice", "0", 1, "CM", false, 0, "")
	pdf.Ln(5)

	// Escreve os detalhes da fatura
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(40, 10, fmt.Sprintf("Invoice ID: %s", i.invoiceID), "0", 0, "", false, 0, "")
	pdf.Ln(10)
	pdf.CellFormat(40, 10, fmt.Sprintf("Invoice Date: %s", i.date.Format("02-01-2006")), "0", 0, "", false, 0, "")
	pdf.Ln(10)
	pdf.CellFormat(40, 10, fmt.Sprintf("Client: %s", i.client), "0", 0, "", false, 0, "")
	pdf.Ln(20)

	// Adiciona uma tabela para os serviços
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(90, 10, "Service", "1", 0, "", false, 0, "")
	pdf.CellFormat(100, 10, "Price", "1", 0, "", false, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	for service, price := range i.services {
		pdf.CellFormat(90, 10, service, "1", 0, "", false, 0, "")
		pdf.CellFormat(100, 10, fmt.Sprintf("%.2f %s", price, i.currency), "1", 0, "", false, 0, "")
		pdf.Ln(10)
	}
	pdf.Ln(10)

	pdf.CellFormat(40, 10, fmt.Sprintf("Total Value: %.2f %s", i.value, i.currency), "0", 0, "", false, 0, "")

	err := pdf.OutputFileAndClose("invoice.pdf")
	if err != nil {
		panic(err)
	}
}

func main() {
	now := time.Now()

	day := flag.Int("day", now.Day(), "Invoice day")
	month := flag.Int("month", int(now.Month()), "Invoice month")
	year := flag.Int("year", now.Year(), "Invoice year")
	value := flag.Float64("value", 0.0, "Invoice value")
	client := flag.String("client", "", "Client name")
	invoiceID := flag.String("id", "", "Invoice ID")
	services := flag.String("services", "", "Comma separated list of services in the format 'name:price'")

	flag.Parse()

	if *day == 0 || *month == 0 || *year == 0 || *value == 0.0 || *client == "" || *invoiceID == "" || *services == "" {
		fmt.Println("Please provide all required parameters.")
		return
	}

	serviceList := make(map[string]float64)
	servicePairs := strings.Split(*services, ",")
	for _, pair := range servicePairs {
		splitPair := strings.Split(pair, ":")
		if len(splitPair) != 2 {
			fmt.Printf("Invalid service: %s\n", pair)
			return
		}
		serviceName := splitPair[0]
		var servicePrice float64
		_, err := fmt.Sscanf(splitPair[1], "%f", &servicePrice)
		if err != nil {
			fmt.Printf("Invalid price for service: %s\n", pair)
			return
		}
		serviceList[serviceName] = servicePrice
	}

	date := time.Date(*year, time.Month(*month), *day, 0, 0, 0, 0, time.UTC)

	invoice := &Invoice{
		date:      date,
		value:     *value,
		currency:  "USD",
		client:    *client,
		invoiceID: *invoiceID,
		services:  serviceList,
	}

	invoice.GeneratePDF()
}
