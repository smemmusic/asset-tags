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
	enc := base32.NewEncoding("ABCDEFGHJKMNPQRSTUVWXYZ123456789").WithPadding(base32.NoPadding)
	id := enc.EncodeToString(buf)
	return id
}

func main() {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetMargins(9.0, 14.0, 9.0)
	pdf.SetAutoPageBreak(false, 14.0)
	pdf.SetFont("Courier", "", 10)

	rows := 27
	cols := 7
	qrsize := 9.5
	colspace := 2.5
	width := 25.4
	height := 10.0
	leftmargin := 10.0
	topmargin := 14.0

	var opt fpdf.ImageOptions
	opt.ImageType = "png"
	opt.AllowNegativePosition = false

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {

			id1 := smolid(3)
			id2 := smolid(3)
			id := id1 + id2
			fmt.Println(id1, id2)
			qr, err := qrcode.Encode(id, qrcode.Highest, 256)
			if err != nil {
				panic("couldn't encode qr")
			}

			offset_x := leftmargin + float64(c)*(width+colspace)
			offset_y := topmargin + float64(r)*height

			// put image on pdf
			fmt.Printf("Placing (%f, %f)\n", offset_x, offset_y)
			_ = pdf.RegisterImageOptionsReader(id, opt, bytes.NewReader(qr))
			pdf.SetXY(offset_x+qrsize, offset_y)
			pdf.CellFormat(width-qrsize, height/2, id1, "0", 0, "CB", false, 0, "")
			pdf.SetXY(offset_x+qrsize, offset_y+5.0)
			pdf.CellFormat(width-qrsize, height/2, id2, "0", 0, "CT", false, 0, "")
			pdf.ImageOptions(id, offset_x, offset_y, qrsize, qrsize, false, opt, 0, "")
		}
	}
	e := pdf.Error()
	if e != nil {
		fmt.Println(e)
		panic("error in pdf")
	}
	pdf.OutputFileAndClose("out.pdf")

}
