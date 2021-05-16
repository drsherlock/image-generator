package imagegen

import (
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"testing"
)

type ContextSpy struct {
	calledDrawImage         bool
	calledImage             bool
	calledSetColor          bool
	calledDrawRectangle     bool
	calledFill              bool
	calledLoadFontFace      bool
	calledDrawStringWrapped bool
	calledSetHexColor       bool
	calledSavedPNG          bool
}

func (dc *ContextSpy) DrawImage(im image.Image, x, y int) {
	dc.calledDrawImage = true
}

func (dc *ContextSpy) Image() image.Image {
	dc.calledImage = true
	return nil
}

func (dc *ContextSpy) SetColor(c color.Color) {
	dc.calledSetColor = true
}

func (dc *ContextSpy) DrawRectangle(x, y, w, h float64) {
	dc.calledDrawRectangle = true
}

func (dc *ContextSpy) Fill() {
	dc.calledFill = true
}

func (dc *ContextSpy) LoadFontFace(path string, points float64) error {
	dc.calledLoadFontFace = true
	return nil
}

func (dc *ContextSpy) DrawStringWrapped(s string, x, y, ax, ay, width, lineSpacing float64, align gg.Align) {
	dc.calledDrawStringWrapped = true
}

func (dc *ContextSpy) SetHexColor(x string) {
	dc.calledSetHexColor = true
}

func (dc *ContextSpy) SavePNG(path string) error {
	dc.calledSavedPNG = true
	return nil
}

var mockFill FillFunc = func(img image.Image, width, height int, anchor imaging.Anchor, filter imaging.ResampleFilter) *image.NRGBA {
	return nil
}

func TestDrawImage(t *testing.T) {
	t.Run("test draw image", func(t *testing.T) {
		im := &Image{}
		spyDC := &ContextSpy{}

		drawImage(im, spyDC, mockFill)

		if !spyDC.calledDrawImage {
			t.Errorf("should have called DrawImage")
		}
		if !spyDC.calledImage {
			t.Errorf("should have called Image")
		}
	})
}

func TestDrawOverlay(t *testing.T) {
	t.Run("test draw overlay", func(t *testing.T) {
		im := &Image{}
		spyDC := &ContextSpy{}

		drawOverlay(im, spyDC)

		if !spyDC.calledSetColor {
			t.Errorf("should have called SetColor")
		}
		if !spyDC.calledDrawRectangle {
			t.Errorf("should have called DrawRectangle")
		}
		if !spyDC.calledFill {
			t.Errorf("should have called Fill")
		}
	})
}

func TestAddText(t *testing.T) {
	t.Run("test add text", func(t *testing.T) {
		im := &Image{}
		spyDC := &ContextSpy{}
		fontPath := ""

		addText(im, spyDC, fontPath)

		if !spyDC.calledLoadFontFace {
			t.Errorf("should have called LoadFontFace")
		}
		if !spyDC.calledDrawStringWrapped {
			t.Errorf("should have called DrawStringWrapped")
		}
		if !spyDC.calledSetHexColor {
			t.Errorf("should have called SetHexColor")
		}
	})
}

func TestSaveImage(t *testing.T) {
	t.Run("test save image", func(t *testing.T) {
		im := &Image{}
		spyDC := &ContextSpy{}

		saveImage(im, spyDC)

		if !spyDC.calledSavedPNG {
			t.Errorf("should have called SaveImage")
		}
	})
}

func BenchmarkDrawImage(b *testing.B) {
	im := &Image{}
	spyDC := &ContextSpy{}

	for i := 0; i < b.N; i++ {
		drawImage(im, spyDC, mockFill)
	}
}

func BenchmarkDrawOverlay(b *testing.B) {
	im := &Image{}
	spyDC := &ContextSpy{}

	for i := 0; i < b.N; i++ {
		drawOverlay(im, spyDC)
	}
}

func BenchmarkAddText(b *testing.B) {
	im := &Image{}
	spyDC := &ContextSpy{}
	fontPath := ""

	for i := 0; i < b.N; i++ {
		addText(im, spyDC, fontPath)
	}
}

func BenchmarkSaveImage(b *testing.B) {
	im := &Image{}
	spyDC := &ContextSpy{}

	for i := 0; i < b.N; i++ {
		saveImage(im, spyDC)
	}
}
