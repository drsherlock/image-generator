package main

import (
	"fmt"
	"github.com/fogleman/gg"
	"image/color"
	"os"
	"path/filepath"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	dc := gg.NewContext(3475, 3503)

	backgroundImage, err := gg.LoadImage("in.jpg")
	if err != nil {
		return fmt.Errorf("load background image %w", err)
	}
	dc.DrawImage(backgroundImage, 0, 0)

	margin := 50
	x := margin
	y := margin
	w := dc.Width() - (2.0 * margin)
	h := dc.Height() - (2.0 * margin)
	dc.SetColor(color.RGBA{0, 0, 0, 150})
	dc.DrawRectangle(float64(x), float64(y), float64(w), float64(h))
	dc.Fill()

	title := "Row! Row! Fight The Powah!"
	textShadowColor := color.Black
	textColor := color.White
	fontPath := filepath.Join("fonts", "BebasNeue-Regular.ttf")
	if err := dc.LoadFontFace(fontPath, 160); err != nil {
		return fmt.Errorf("load font face %w", err)
	}
	textRightMargin := 90
	textTopMargin := 120
	x = textRightMargin
	y = textTopMargin
	maxWidth := dc.Width() - textRightMargin - textRightMargin
	dc.SetColor(textShadowColor)
	dc.DrawStringWrapped(title, float64(x+1), float64(y+1), 0, 0, float64(maxWidth), 1.5, gg.AlignLeft)
	dc.SetColor(textColor)
	dc.DrawStringWrapped(title, float64(x), float64(y), 0, 0, float64(maxWidth), 1.5, gg.AlignLeft)

	if err := dc.SavePNG("out.png"); err != nil {
		return fmt.Errorf("save png %w", err)
	}
	return nil
}
