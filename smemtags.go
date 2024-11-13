package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"fmt"

	"github.com/go-pdf/fpdf"
	"github.com/skip2/go-qrcode"
)

func smolid(b int) string {
	buf := make([]byte, b)
	_, err := rand.Read(buf)
	if err != nil {
		panic("couldn't read random bytes")
	}
	enc := base32.StdEncoding.WithPadding(base32.NoPadding)
	id := enc.EncodeToString(buf)
	return id
}

func wraptext(s string) string {
	var out string
	for i, c := range s {
		out += string(c)
		if i+1 == len(s)/2 {
			out += "\n"
		}
	}
	return out
}

func main() {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetMargins(9.0, 13.0, 9.0)
	pdf.SetFont("Courier", "", 8)

	rows := 27
	cols := 7
	qrsize := 8.0
	colspace := 3.0
	leftmargin := 9.0
	topmargin := 13.0

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {

			id := smolid(8)
			id_split := wraptext(id)
			fmt.Println(id)
			qr, err := qrcode.Encode(id, qrcode.Medium, 100)
			if err != nil {
				panic("couldn't encode qr")
			}

			offset_x := leftmargin + float64(c)*(25.98+colspace)
			offset_y := topmargin + float64(r)*10.0
			text_x := offset_x + 13.0
			text_y := offset_y

			// put image on pdf
			var opt fpdf.ImageOptions
			opt.ImageType = "png"
			opt.AllowNegativePosition = false
			_ = pdf.RegisterImageOptionsReader(id, opt, bytes.NewReader(qr))
			pdf.ImageOptions(id, offset_x, offset_y, qrsize, qrsize, false, opt, 0, "")

			fmt.Printf("Placing (%f, %f), (%f, %f)\n", offset_x, offset_y, text_x, text_y)
			pdf.SetXY(offset_x, offset_y)
			pdf.CellFormat(25.4, 10.0, id_split, "1", 0, "LM", false, 0, "")
		}
	}
	e := pdf.Error()
	if e != nil {
		fmt.Println(e)
		panic("error in pdf")
	}
	pdf.OutputFileAndClose("out.pdf")

}
