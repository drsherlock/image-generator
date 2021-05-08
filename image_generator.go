package main

import (
	"flag"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"image/color"
	"os"
	"path/filepath"
	"sync"
)

type Image struct {
	fileName   string
	title      string
	titleColor string
	fontName   string
	width      float64
	height     float64
}

func main() {
	// get file name, title from arguments
	fileName := flag.String("file-name", "input.jpg", "name of image")
	title := flag.String("title", "Yello!", "title of image")
	titleColor := flag.String("title-color", "#ffffff", "color of title")
	flag.Parse()

	image := Image{
		fileName:   *fileName,
		title:      *title,
		titleColor: *titleColor,
	}

	if err := run(image); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(image Image) error {
	var wg sync.WaitGroup
	err := filepath.Walk("fonts", func(fontPath string, fontInfo os.FileInfo, err error) error {
		image.fontName = fontInfo.Name()

		if !fontInfo.IsDir() {
			wg.Add(1)
			go generate(image, fontPath, &wg)
		}
		return nil
	})

	wg.Wait()

	if err != nil {
		return fmt.Errorf("load font face %w", err)
	}

	return nil
}

func generate(image Image, fontPath string, wg *sync.WaitGroup) error {
	defer wg.Done()

	backgroundImage, err := gg.LoadImage(image.fileName)
	if err != nil {
		return fmt.Errorf("load background image %w", err)
	}
	dc := gg.NewContextForImage(backgroundImage)

	image.width = float64(dc.Width())
	image.height = float64(dc.Height())

	// draw the image
	drawImage(image, dc)

	// draw the overlay
	drawOverlay(image, dc)

	// add text
	err = addText(image, dc, fontPath)
	if err != nil {
		return fmt.Errorf("add text %w", err)
	}

	// save image
	err = saveImage(image, dc)
	if err != nil {
		return fmt.Errorf("save image %w", err)
	}

	return nil
}

func drawImage(image Image, dc *gg.Context) {
	backgroundImage := imaging.Fill(dc.Image(), dc.Width(), dc.Height(), imaging.Center, imaging.Lanczos)
	dc.DrawImage(backgroundImage, 0, 0)
}

func drawOverlay(image Image, dc *gg.Context) {
	x := image.width / 50
	y := image.height / 50
	w := image.width - (2.0 * x)
	h := image.height - (2.0 * y)
	dc.SetColor(color.RGBA{0, 0, 0, 150})
	dc.DrawRectangle(x, y, w, h)
	dc.Fill()
}

func addText(image Image, dc *gg.Context, fontPath string) error {
	textShadowColor := color.Black
	if err := dc.LoadFontFace(fontPath, 200); err != nil {
		return err
	}
	textMargin := image.width * 0.03
	textTopMargin := image.height / 2.5
	x := textMargin
	y := textTopMargin
	maxWidth := image.width - textMargin - textMargin

	dc.SetColor(textShadowColor)
	dc.DrawStringWrapped(image.title, x+1, y+1, 0, 0, float64(maxWidth), 1.5, gg.AlignCenter)
	dc.SetHexColor(image.titleColor)
	dc.DrawStringWrapped(image.title, x, y, 0, 0, float64(maxWidth), 1.5, gg.AlignCenter)

	return nil
}

func saveImage(image Image, dc *gg.Context) error {
	if err := dc.SavePNG("output/" + image.fontName + ".png"); err != nil {
		return err
	}

	return nil
}
