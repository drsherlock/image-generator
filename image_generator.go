package main

import (
	"flag"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"image/color"
	"os"
	"path/filepath"
)

func main() {
	// get file name, title from arguments
	fileName := flag.String("file-name", "input.jpg", "name of image")
	title := flag.String("title", "Yello!", "title of image")
	titleColor := flag.String("title-color", "#ffffff", "color of title")
	flag.Parse()

	if err := run(fileName, title, titleColor); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(fileName *string, title *string, titleColor *string) error {
	// draw the image
	backgroundImage, err := gg.LoadImage(*fileName)

	if err != nil {
		return fmt.Errorf("load background image %w", err)
	}

	dc := gg.NewContextForImage(backgroundImage)

	backgroundImage = imaging.Fill(backgroundImage, dc.Width(), dc.Height(), imaging.Center, imaging.Lanczos)

	dc.DrawImage(backgroundImage, 0, 0)

	// draw the overlay
	imageWidth := float64(dc.Width())
	imageHeight := float64(dc.Height())

	x := imageWidth / 50
	y := imageHeight / 50
	w := imageWidth - (2.0 * x)
	h := imageHeight - (2.0 * y)
	dc.SetColor(color.RGBA{0, 0, 0, 150})
	dc.DrawRectangle(x, y, w, h)
	dc.Fill()

	// add text
	textShadowColor := color.Black
	fontPath := filepath.Join("fonts", "BebasNeue-Regular.ttf")
	if err := dc.LoadFontFace(fontPath, 200); err != nil {
		return fmt.Errorf("load font face %w", err)
	}
	textMargin := x * 1.5
	textTopMargin := imageHeight / 2.5
	x = textMargin
	y = textTopMargin
	maxWidth := imageWidth - textMargin - textMargin

	dc.SetColor(textShadowColor)
	dc.DrawStringWrapped(*title, x+1, y+1, 0, 0, float64(maxWidth), 1.5, gg.AlignCenter)
	dc.SetHexColor(*titleColor)
	dc.DrawStringWrapped(*title, x, y, 0, 0, float64(maxWidth), 1.5, gg.AlignCenter)

	// save the new image
	if err := dc.SavePNG("out.png"); err != nil {
		return fmt.Errorf("save png %w", err)
	}
	return nil
}
