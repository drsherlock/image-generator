package main

import (
	"fmt"
	"github.com/fogleman/gg"
	"os"
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

	if err := dc.SavePNG("out.png"); err != nil {
		return fmt.Errorf("save png %w", err)
	}
	return nil
}
