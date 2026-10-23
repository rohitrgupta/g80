package spriteedit

import (
	"image"

	"g80/pkg/gfx"
	"g80/pkg/ui"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	sliderStep  = 6
	sliderKnob  = 5
	sliderWidth = 10
)

// SizeSlider is a vertical slider choosing the edit region size; the smallest size is at the bottom.
type SizeSlider struct {
	ui.Base
	st *State
}

func newSizeSlider(st *State, opts ...widget.WidgetOpt) *SizeSlider {
	s := &SizeSlider{st: st}
	n := len(Sizes)
	s.Base = ui.NewBase(sliderWidth, (n-1)*sliderStep+sliderKnob, append(opts,
		widget.WidgetOpts.MouseButtonPressedHandler(func(a *widget.WidgetMouseButtonPressedEventArgs) {
			if a.Button != ebiten.MouseButtonLeft {
				return
			}
			stop := (a.OffsetY - sliderKnob/2 + sliderStep/2) / sliderStep
			i := n - 1 - min(max(stop, 0), n-1)
			st.SetSize(Sizes[i])
		}),
	)...)
	return s
}

// knobTop returns the y offset of the knob for size index i.
func (s *SizeSlider) knobTop(i int) int { return (len(Sizes) - 1 - i) * sliderStep }

func (s *SizeSlider) Render(screen *ebiten.Image) {
	s.Base.Render(screen)
	r := s.Rect()
	x, y := r.Min.X, r.Min.Y
	mid := sliderKnob / 2

	gfx.FillRect(screen, image.Rect(x+mid, y+mid, x+mid+1, r.Max.Y-mid), ui.ColorBackground)
	for i, size := range Sizes {
		ty := y + s.knobTop(i) + mid
		gfx.FillRect(screen, image.Rect(x+sliderKnob+1, ty, x+sliderKnob+2+i*2, ty+1), ui.ColorBackground)
		if size == s.st.Size {
			knob := image.Rect(x, y+s.knobTop(i), x+sliderKnob, y+s.knobTop(i)+sliderKnob)
			gfx.FillRect(screen, knob, ui.ColorBackground)
			gfx.FillRect(screen, knob.Inset(1), ui.ColorMuted)
		}
	}
}
