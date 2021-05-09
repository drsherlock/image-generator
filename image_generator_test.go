package main

import (
	// "github.com/fogleman/gg"
	"github.com/disintegration/imaging"
	"image"
	"testing"
)

type ContextSpy struct {
	CalledDrawImage bool
	CalledImage     bool
}

func (dc *ContextSpy) DrawImage(im image.Image, x, y int) {
	dc.CalledDrawImage = true
}

func (dc *ContextSpy) Image() image.Image {
	dc.CalledImage = true
	return nil
}

var mockFill FillFunc = func(img image.Image, width, height int, anchor imaging.Anchor, filter imaging.ResampleFilter) *image.NRGBA {
	return nil
}

func TestDrawImage(t *testing.T) {
	t.Run("test draw image", func(t *testing.T) {
		im := &Image{}

		spyDc := &ContextSpy{}

		DrawImage(im, spyDc, mockFill)

		if !spyDc.CalledDrawImage {
			t.Errorf("should have called DrawImage")
		}
		if !spyDc.CalledImage {
			t.Errorf("should have called Image")
		}

	})
}

func BenchmarkDrawImage(b *testing.B) {
	im := &Image{}

	spyDc := &ContextSpy{}
	for i := 0; i < b.N; i++ {
		DrawImage(im, spyDc, mockFill)
	}
}
