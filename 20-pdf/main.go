package main

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

type LineItem struct {
	UnitName       string
	PricePerUnit   int
	UnitsPurchased int
}

const (
	bannerHeight = 94.0
	xIndent      = 40.0
	taxRate      = 18
)

func main() {
	lineItems := []LineItem{
		{
			UnitName:       "2x6 Lumber - 8'",
			PricePerUnit:   375,
			UnitsPurchased: 220,
		}, {
			UnitName:       "Drywall Sheet",
			PricePerUnit:   822,
			UnitsPurchased: 50,
		}, {
			UnitName:       "Paint",
			PricePerUnit:   1455,
			UnitsPurchased: 3,
		}, {
			UnitName:       "This is a line item with a very long description to test that our word wrapping is implemented and working as intended",
			PricePerUnit:   3211,
			UnitsPurchased: 3,
		}, {
			UnitName:       "Paint",
			PricePerUnit:   5,
			UnitsPurchased: 3300,
		}, {
			UnitName:       "Paint",
			PricePerUnit:   332,
			UnitsPurchased: 44,
		},
	}
	subtotal := 0
	for _, item := range lineItems {
		subtotal += item.PricePerUnit * item.UnitsPurchased
	}
	tax := int(float64(subtotal) * taxRate)
	total := subtotal + tax
	totalAmount := toInr(total)

	pdf := gofpdf.New(gofpdf.OrientationPortrait, gofpdf.UnitPoint, gofpdf.PageSizeLetter, "")
	w, h := pdf.GetPageSize()
	pdf.AddPage()

	// Header
	pdf.SetFillColor(103, 60, 79)
	pdf.Polygon([]gofpdf.PointType{
		{0, 0},
		{w, 0},
		{w, bannerHeight},
		{0, bannerHeight * 0.9},
	}, "F")
	pdf.Polygon([]gofpdf.PointType{
		{0, h},
		{0, h - (bannerHeight * 0.2)},
		{w, h - (bannerHeight * 0.1)},
		{w, h},
	}, "F")

	pdf.SetFont("Arial", "B", 40)
	pdf.SetTextColor(255, 255, 255)
	_, lineHeight := pdf.GetFontSize()
	pdf.Text(xIndent, bannerHeight-(bannerHeight/2.0)+lineHeight/3, "INVOICE")

	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(255, 255, 255)
	_, lineHeight = pdf.GetFontSize()
	pdf.MoveTo(w-xIndent-2.0*124.0, (bannerHeight-(lineHeight*1.5*3.0))/2.0)
	pdf.MultiCell(124.0, lineHeight*1.5, "(123) 456-7890\nabhishek@gmail.com\nwrongabhishek.com", gofpdf.BorderNone, gofpdf.AlignRight, false)

	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(255, 255, 255)
	_, lineHeight = pdf.GetFontSize()
	pdf.MoveTo(w-xIndent-124.0, (bannerHeight-(lineHeight*1.5*3.0))/2.0)
	pdf.MultiCell(124.0, lineHeight*1.5, "123 Fake St\nSome City, BLR\n123456", gofpdf.BorderNone, gofpdf.AlignRight, false)

	// Invoice Details
	summaryBlock(pdf, xIndent, bannerHeight+lineHeight*2.0, "Billed To", "Client Name", "123 Client Address", "City, State, Country", "Postal Code")
	summaryBlock(pdf, xIndent*2.0+lineHeight*12, bannerHeight+lineHeight*2.0, "Invoice Number", "002148234675")
	_, sy := summaryBlock(pdf, xIndent*2.0+lineHeight*12, bannerHeight+lineHeight*6.0, "Date of Issue", "01/01/2023")

	x, y := w-xIndent-124.0, bannerHeight+lineHeight*2.0
	pdf.MoveTo(x, y)
	pdf.SetFont("Times", "", 14)
	_, lineHeight = pdf.GetFontSize()
	pdf.SetTextColor(180, 180, 180)
	pdf.CellFormat(124.0, lineHeight, "Invoice Total", gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")

	y = y + lineHeight*1.5
	pdf.MoveTo(x, y)
	pdf.SetFont("Times", "", 32)
	_, lineHeight = pdf.GetFontSize()
	pdf.SetTextColor(130, 100, 113)
	pdf.CellFormat(124.0, lineHeight, totalAmount, gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")

	x, y = x-2.0, y+lineHeight*1.25

	if sy > y {
		y = sy
	}
	x, y = xIndent-20.0, y+30.0
	pdf.Rect(x, y, w-(xIndent*2.0)+40.0, 3.0, "F")

	// Invoice Breakdown
	pdf.SetFont("Times", "", 14)
	_, lineHeight = pdf.GetFontSize()
	pdf.SetTextColor(180, 180, 180)
	x, y = xIndent-2.0, y+lineHeight
	pdf.MoveTo(x, y)
	pdf.CellFormat(w/2.5+1.5, lineHeight, "Description", gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignLeft, false, 0, "")
	x = x + w/2.5 + 1.5
	pdf.MoveTo(x, y)
	pdf.CellFormat(100.0, lineHeight, "Price Per Unit", gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	x = x + 100.0
	pdf.MoveTo(x, y)
	pdf.CellFormat(80.0, lineHeight, "Quantity", gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	x = w - xIndent - 2.0 - 119.5
	pdf.MoveTo(x, y)
	pdf.CellFormat(119.5, lineHeight, "Amount", gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")

	y = y + lineHeight
	for _, li := range lineItems {
		x, y = lineItem(pdf, x, y, li)
	}

	x, y = w/1.75, y+lineHeight*2.25
	x, y = trailerLine(pdf, x, y, "Subtotal", subtotal)
	x, y = trailerLine(pdf, x, y, "Tax", tax)
	pdf.SetDrawColor(180, 180, 180)
	pdf.Line(x+20.0, y, x+220.0, y)
	y = y + lineHeight*0.5
	x, y = trailerLine(pdf, x, y, "Total", total)

	err := pdf.OutputFileAndClose("output.pdf")
	if err != nil {
		panic(err)
	}
}

func summaryBlock(pdf *gofpdf.Fpdf, x, y float64, title string, data ...string) (float64, float64) {
	pdf.SetFont("Times", "", 14)
	pdf.SetTextColor(180, 180, 180)
	_, lineHeight := pdf.GetFontSize()
	y = y + lineHeight
	pdf.Text(x, y, title)
	y = y + lineHeight*0.25
	for _, d := range data {
		pdf.SetTextColor(50, 50, 50)
		y = y + lineHeight*1.25
		pdf.Text(x, y, d)
	}

	return x, y
}

func lineItem(pdf *gofpdf.Fpdf, x, y float64, lineItem LineItem) (float64, float64) {
	origX := x
	w, _ := pdf.GetPageSize()
	pdf.SetFont("times", "", 14)
	_, lineHt := pdf.GetFontSize()
	pdf.SetTextColor(50, 50, 50)
	pdf.MoveTo(x, y)
	x, y = xIndent-2.0, y+lineHt*.75
	pdf.MoveTo(x, y)
	pdf.MultiCell(w/2.65+1.5, lineHt, lineItem.UnitName, gofpdf.BorderNone, gofpdf.AlignLeft, false)
	tmp := pdf.SplitLines([]byte(lineItem.UnitName), w/2.65+1.5)
	maxY := y + float64(len(tmp)-1)*lineHt
	x = x + w/2.65 + 1.5
	pdf.MoveTo(x, y)
	pdf.CellFormat(100.0, lineHt, toInr(lineItem.PricePerUnit), gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	x = x + 100.0
	pdf.MoveTo(x, y)
	pdf.CellFormat(80.0, lineHt, fmt.Sprintf("%d", lineItem.UnitsPurchased), gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	x = w - xIndent - 2.0 - 119.5
	pdf.MoveTo(x, y)
	pdf.CellFormat(119.5, lineHt, toInr(lineItem.PricePerUnit*lineItem.UnitsPurchased), gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	if maxY > y {
		y = maxY
	}
	y = y + lineHt*1.75
	pdf.SetDrawColor(180, 180, 180)
	pdf.Line(xIndent-10.0, y, w-xIndent+10.0, y)
	return origX, y
}

func toInr(cents int) string {
	centsStr := fmt.Sprintf("%d", cents%100)
	if len(centsStr) < 2 {
		centsStr = "0" + centsStr
	}
	return fmt.Sprintf("INR %d.%s", cents/100, centsStr)
}

func trailerLine(pdf *gofpdf.Fpdf, x, y float64, label string, amount int) (float64, float64) {
	origX := x
	w, _ := pdf.GetPageSize()
	pdf.SetFont("times", "", 14)
	_, lineHt := pdf.GetFontSize()
	pdf.SetTextColor(180, 180, 180)
	pdf.MoveTo(x, y)
	pdf.CellFormat(80.0, lineHt, label, gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	x = w - xIndent - 2.0 - 119.5
	pdf.MoveTo(x, y)
	pdf.SetTextColor(50, 50, 50)
	pdf.CellFormat(119.5, lineHt, toInr(amount), gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	y = y + lineHt*1.5
	return origX, y
}
