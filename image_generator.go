package imagegen

import (
	"flag"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"sync"
)

type Image struct {
	baseFileName string
	title        string
	titleColor   string
	fontName     string
	width        float64
	height       float64
}

type DrawingContext interface {
	DrawImage(im image.Image, x, y int)
	Image() image.Image
	SetColor(c color.Color)
	DrawRectangle(x, y, w, h float64)
	Fill()
	LoadFontFace(path string, points float64) error
	DrawStringWrapped(s string, x, y, ax, ay, width, lineSpacing float64, align gg.Align)
	SetHexColor(x string)
	SavePNG(path string) error
}

type FillFunc func(img image.Image, width, height int, anchor imaging.Anchor, filter imaging.ResampleFilter) *image.NRGBA

func Create(imFile *os.File, title string, titleColor string, fonts []string) error {
	var wg sync.WaitGroup
	err := filepath.Walk("fonts", func(fontPath string, fontInfo os.FileInfo, err error) error {
		for _, f := range fonts {
			if f+".ttf" == fontInfo.Name() {
				im := &Image{
					baseFileName: imFile.Name(),
					title:        title,
					titleColor:   titleColor,
					fontName:     fontInfo.Name(),
				}

				if !fontInfo.IsDir() {
					wg.Add(1)
					go generate(im, fontPath, &wg)
				}
			}
		}
		return nil
	})

	wg.Wait()

	if err != nil {
		return fmt.Errorf("load font face %w", err)
	}

	return nil
}

func main() {
	// get file name, title from arguments
	baseFileName := flag.String("file-name", "input.jpg", "name of image")
	title := flag.String("title", "Yello!", "title of image")
	titleColor := flag.String("title-color", "#ffffff", "color of title")
	flag.Parse()

	if err := run(baseFileName, title, titleColor); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(baseFileName *string, title *string, titleColor *string) error {
	var wg sync.WaitGroup
	err := filepath.Walk("fonts", func(fontPath string, fontInfo os.FileInfo, err error) error {
		im := &Image{
			baseFileName: *baseFileName,
			title:        *title,
			titleColor:   *titleColor,
			fontName:     fontInfo.Name(),
		}

		if !fontInfo.IsDir() {
			wg.Add(1)
			go generate(im, fontPath, &wg)
		}
		return nil
	})

	wg.Wait()

	if err != nil {
		return fmt.Errorf("load font face %w", err)
	}

	return nil
}

func generate(im *Image, fontPath string, wg *sync.WaitGroup) error {
	defer wg.Done()

	backgroundImage, err := gg.LoadImage(im.baseFileName)
	if err != nil {
		return fmt.Errorf("load background image %w", err)
	}
	dc := gg.NewContextForImage(backgroundImage)

	im.width = float64(dc.Width())
	im.height = float64(dc.Height())

	// draw the image
	drawImage(im, dc, imaging.Fill)

	// draw the overlay
	drawOverlay(im, dc)

	// add text
	err = addText(im, dc, fontPath)
	if err != nil {
		return fmt.Errorf("add text %w", err)
	}

	// save image
	err = saveImage(im, dc)
	if err != nil {
		return fmt.Errorf("save image %w", err)
	}

	return nil
}

func drawImage(im *Image, dc DrawingContext, fill FillFunc) {
	backgroundImage := fill(dc.Image(), int(im.width), int(im.height), imaging.Center, imaging.Lanczos)
	dc.DrawImage(backgroundImage, 0, 0)
}

func drawOverlay(im *Image, dc DrawingContext) {
	x := im.width / 50
	y := im.height / 50
	w := im.width - (2.0 * x)
	h := im.height - (2.0 * y)
	dc.SetColor(color.RGBA{150, 150, 150, 150})
	dc.DrawRectangle(x, y, w, h)
	dc.Fill()
}

func addText(im *Image, dc DrawingContext, fontPath string) error {
	textShadowColor := color.Black
	if err := dc.LoadFontFace(fontPath, 200); err != nil {
		return err
	}
	textMargin := im.width * 0.03
	textTopMargin := im.height / 2.5
	x := textMargin
	y := textTopMargin
	maxWidth := im.width - textMargin - textMargin

	dc.SetColor(textShadowColor)
	dc.DrawStringWrapped(im.title, x+1, y+1, 0, 0, float64(maxWidth), 1.5, gg.AlignCenter)
	dc.SetHexColor(im.titleColor)
	dc.DrawStringWrapped(im.title, x, y, 0, 0, float64(maxWidth), 1.5, gg.AlignCenter)

	return nil
}

func saveImage(im *Image, dc DrawingContext) error {
	err := os.MkdirAll("./output/"+im.baseFileName, os.ModePerm)
	if err != nil {
		return err
	}

	outputPath := "output/" + im.baseFileName + "/" + im.fontName + "_" + im.titleColor + ".png"
	if err := dc.SavePNG(outputPath); err != nil {
		return err
	}

	return nil
}
