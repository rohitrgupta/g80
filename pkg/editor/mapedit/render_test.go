package mapedit

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"g80/pkg/cart"
	"g80/pkg/ui"

	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
)

type renderGame struct {
	t      *testing.T
	editor *Editor
	ui     *ebitenui.UI
	stage  int
	frames int
	dir    string
}

func (g *renderGame) Update() error {
	if g.stage == 4 {
		return ebiten.Termination
	}
	g.ui.Update()
	return nil
}

func (g *renderGame) Draw(screen *ebiten.Image) {
	if g.stage >= 4 {
		return
	}
	g.ui.Draw(screen)
	g.frames++
	if g.frames < 2 {
		return
	}
	names := []string{"map-editor.png", "repeat-paste.png", "sprite-picker.png", "map-selector.png"}
	pixels := make([]byte, ui.ScreenWidth*ui.ScreenHeight*4)
	screen.ReadPixels(pixels)
	colors := map[[3]byte]bool{}
	for index := 0; index < len(pixels); index += 4 {
		colors[[3]byte{pixels[index], pixels[index+1], pixels[index+2]}] = true
	}
	if len(colors) < 3 {
		g.t.Error("rendered view is blank")
	}
	file, err := os.Create(filepath.Join(g.dir, names[g.stage]))
	if err != nil {
		g.t.Fatal(err)
	}
	img := &image.RGBA{Pix: pixels, Stride: ui.ScreenWidth * 4, Rect: image.Rect(0, 0, ui.ScreenWidth, ui.ScreenHeight)}
	if err := png.Encode(file, img); err != nil {
		g.t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		g.t.Fatal(err)
	}
	g.stage++
	g.frames = 0
	if g.stage == 1 {
		g.editor.st.selection = image.Rect(0, 0, 2, 2)
		g.editor.st.copy()
		g.editor.st.paste()
	}
	if g.stage == 2 {
		g.editor.st.pasting = false
		g.editor.st.choosing = true
	}
	if g.stage == 3 {
		g.editor.st.overview = true
	}
	if g.stage < 4 {
		g.editor.build()
	}
}

func (g *renderGame) Layout(int, int) (int, int) { return ui.ScreenWidth, ui.ScreenHeight }

func TestRenderViews(t *testing.T) {
	dir := os.Getenv("G80_RENDER_DIR")
	if dir == "" {
		t.Skip("set G80_RENDER_DIR to capture map editor views")
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	c := cart.New()
	for index := 0; index < 32; index++ {
		region := c.Sprites.TileRect(index)
		for row := region.Min.Y; row < region.Max.Y; row++ {
			for col := region.Min.X; col < region.Max.X; col++ {
				c.Sprites.Set(col, row, uint8(index%15+1))
			}
		}
	}
	e := New(c).(*Editor)
	for col := 0; col < 32; col++ {
		c.Maps[0].Set(col, col/2, uint8(col))
	}
	e.st.selectBrush(image.Pt(0, 0), image.Pt(1, 1))
	game := &renderGame{t: t, editor: e, ui: &ebitenui.UI{Container: e.root}, dir: dir}
	ebiten.SetWindowSize(960, 544)
	if err := ebiten.RunGame(game); err != nil {
		t.Fatal(err)
	}
	t.Logf("screenshots: %s", dir)
}
