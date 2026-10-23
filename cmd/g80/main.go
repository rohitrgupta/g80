package main

import (
	"errors"
	"flag"
	"io/fs"
	"log"

	"g80/pkg/app"
	"g80/pkg/cart"
	"g80/pkg/editor/spriteedit"
	"g80/pkg/ui"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	path := flag.String("cart", "cart.g80", "plain-text cart file to load and save")
	scale := flag.Int("scale", 4, "window scale")
	flag.Parse()

	c, err := cart.Load(*path)
	if errors.Is(err, fs.ErrNotExist) {
		c = cart.New()
	} else if err != nil {
		log.Fatalf("load %s: %v", *path, err)
	}

	game := app.New(c, *path, spriteedit.New)

	ebiten.SetWindowTitle("g80 - " + *path)
	ebiten.SetWindowSize(ui.ScreenWidth**scale, ui.ScreenHeight**scale)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
