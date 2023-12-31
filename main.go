package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type Invoice struct {
	date          time.Time
	value         float64
	currency      string
	client        string
	invoiceID     string
	services      map[string]float64
	originName    string
	originAddress string
}

func (i *Invoice) GeneratePDF() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)

	// Add a title
	pdf.CellFormat(190, 10, "Invoice", "0", 1, "CM", false, 0, "")
	pdf.Ln(5)

	// Write the origin details
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(40, 10, fmt.Sprintf("%s", i.originName), "0", 0, "", false, 0, "")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(40, 10, i.originAddress, "0", 0, "", false, 0, "")
	pdf.Ln(20)

	// Write the invoice details
	pdf.CellFormat(40, 10, fmt.Sprintf("Invoice #: %s", i.invoiceID), "0", 0, "", false, 0, "")
	pdf.Ln(10)
	pdf.CellFormat(40, 10, fmt.Sprintf("Invoice Date: %s", i.date.Format("02-01-2006")), "0", 0, "", false, 0, "")
	pdf.Ln(10)
	pdf.CellFormat(40, 10, fmt.Sprintf("Bill To: %s", i.client), "0", 0, "", false, 0, "")
	pdf.Ln(20)

	// Add a table for the services
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(90, 10, "Service", "1", 0, "", false, 0, "")
	pdf.CellFormat(100, 10, "Price", "1", 0, "", false, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	for service, price := range i.services {
		pdf.CellFormat(90, 10, service, "1", 0, "", false, 0, "")
		pdf.CellFormat(100, 10, fmt.Sprintf("%s %.2f", i.currency, price), "1", 0, "", false, 0, "")
		pdf.Ln(10)
	}
	pdf.Ln(10)

	pdf.CellFormat(40, 10, fmt.Sprintf("Total Value: %s %.2f", i.currency, i.value), "0", 0, "", false, 0, "")

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
	client := flag.String("client", "", "Client name")
	invoiceID := flag.String("id", "", "Invoice ID")
	services := flag.String("services", "", "Comma separated list of services in the format 'name:price'")
	originName := flag.String("originName", "", "Origin Name")
	originAddress := flag.String("originAddress", "", "Origin Address")

	flag.Parse()

	if *day == 0 || *month == 0 || *year == 0 || *client == "" || *invoiceID == "" || *services == "" || *originName == "" || *originAddress == "" {
		fmt.Println("Please provide all required parameters.")
		return
	}

	value := 0.0
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
		value += servicePrice
	}

	date := time.Date(*year, time.Month(*month), *day, 0, 0, 0, 0, time.UTC)

	invoice := &Invoice{
		date:          date,
		value:         value,
		currency:      "$",
		client:        *client,
		invoiceID:     *invoiceID,
		services:      serviceList,
		originName:    *originName,
		originAddress: *originAddress,
	}

	invoice.GeneratePDF()
}
