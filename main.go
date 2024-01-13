package main

import (
	"io"
	"net/http"
	"os"
	"time"

	gofpdf "github.com/go-pdf/fpdf"
	"github.com/go-pdf/fpdf/contrib/gofpdi"
)

func main() {

	updatePDF()

}

func updatePDF() ([]byte, error) {

	pdf := gofpdf.New("P", "mm", "A4", "")
	//	Download a PDF
	fileUrl := "https://tcpdf.org/files/examples/example_027.pdf"
	if err := DownloadFile("example.pdf", fileUrl); err != nil {
		panic(err)
	}
	// create a new Importer instance
	imp := gofpdi.NewImporter()
	//pdf.SetFont("Arial", "", 10)

	// pdf.AddUTF8Font("amatic", "", "AmaticSC-Regular.ttf")
	// pdf.SetFont("amatic", "", 10)
	pdf.SetFont("Helvetica", "", 10)
	pdf.AddPage()
	tpl := imp.ImportPage(pdf, "example.pdf", 1, "/MediaBox")

	nrPages := len(imp.GetPageSizes())

	//add download date on first page of pdf

	imp.UseImportedTemplate(pdf, tpl, 0, 0, 210, 297)

	format := "15:04:05 UTC 02 January 2006"
	downloadTime := time.Now().UTC().Format(format)

	pdf.Ln(15)
	pdf.Cellf(120, 0, "downloaded at %s", downloadTime)

	// add download date on first page of pdf
	pdf.AddPage()
	imp.UseImportedTemplate(pdf, tpl, 0, 0, 210, 297)

	// add all pages from template pdf
	if nrPages > 1 {
		for i := 2; i <= nrPages; i++ {
			pdf.AddPage()
			tpl := imp.ImportPage(pdf, "example.pdf", i, "/MediaBox")
			imp.UseImportedTemplate(pdf, tpl, 0, 0, 210, 297)
		}
	}

	err := pdf.OutputFileAndClose("with-download-date.pdf")
	if err != nil {
		return nil, err
	}
	//convert downloaded file to bytes
	bytes, err := os.ReadFile("with-download-date.pdf")
	if err != nil {
		return nil, err
	}

	return bytes, nil

}
func DownloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
