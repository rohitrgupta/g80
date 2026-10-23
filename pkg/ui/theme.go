package ui

import "g80/pkg/gfx"

// Logical screen layout (scaled up to the window).
const (
	ScreenWidth   = 240
	ScreenHeight  = 136
	BarHeight     = 7
	ContentWidth  = ScreenWidth
	ContentHeight = ScreenHeight - BarHeight
)

// UI colors, taken from the default palette so the UI is unaffected by cart palette edits.
var (
	ColorBackground = gfx.Sweetie16[0]
	ColorWhite      = gfx.Sweetie16[12]
	ColorMuted      = gfx.Sweetie16[13]
	ColorPanel      = gfx.Sweetie16[14]
	ColorDark       = gfx.Sweetie16[15]
)
