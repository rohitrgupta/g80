package editor

import (
	"g80/pkg/cart"
	"g80/pkg/gfx"

	"github.com/ebitenui/ebitenui/widget"
)

// Editor is a full-screen tool (sprites, map, code...) hosted below the app's top bar.
type Editor interface {
	Name() string
	Icon() *gfx.Bitmap
	// Root is the editor's widget tree, sized ui.ContentWidth x ui.ContentHeight.
	Root() widget.PreferredSizeLocateableWidget
}

// Factory creates an editor working on a cart.
type Factory func(c *cart.Cart) Editor

// Commands is optionally implemented by editors that support history and clipboard;
// the app routes top bar buttons and keyboard shortcuts to it.
type Commands interface {
	Undo()
	Redo()
	Cut()
	Copy()
	Paste()
}
