package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"fmt"

	"github.com/go-pdf/fpdf"
	qrcode "github.com/skip2/go-qrcode"
)

func smolid(b int) string {
	buf := make([]byte, b)
	_, err := rand.Read(buf)
	if err != nil {
		panic("couldn't read random bytes")
	}
	id := base32.StdEncoding.EncodeToString(buf)
	return id
}

func main() {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	rows := 27
	cols := 7
	qrsize := 8.0
	colspace := 3.0
	header := 10.0

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {

			id := smolid(8)
			fmt.Println(id)
			qr, err := qrcode.Encode(id, qrcode.Low, 100)
			if err != nil {
				panic("couldn't encode qr")
			}

			offset_x := 10.0 + float64(c)*(25.4+colspace)
			offset_y := header + float64(r)*10.0
			text_x := offset_x + 12.2
			text_y := offset_y

			// put image on pdf
			var opt fpdf.ImageOptions
			opt.ImageType = "png"
			opt.AllowNegativePosition = false
			_ = pdf.RegisterImageOptionsReader(id, opt, bytes.NewReader(qr))
			pdf.ImageOptions(id, offset_x, offset_y, qrsize, qrsize, false, opt, 0, "")
			pdf.Text(text_x, text_y, id)
		}
	}
	pdf.OutputFileAndClose("out.pdf")

}
