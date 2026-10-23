package sprite

import "image"

// Line draws a line from (x0, y0) to (x1, y1), only touching pixels inside clip.
func Line(s *Sheet, clip image.Rectangle, x0, y0, x1, y1 int, c uint8) {
	dx, dy := abs(x1-x0), -abs(y1-y0)
	sx, sy := 1, 1
	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}
	err := dx + dy
	for {
		if image.Pt(x0, y0).In(clip) {
			s.Set(x0, y0, c)
		}
		if x0 == x1 && y0 == y1 {
			return
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

// FloodFill replaces the 4-connected area of same-colored pixels at (x, y) within clip.
func FloodFill(s *Sheet, clip image.Rectangle, x, y int, c uint8) {
	if !image.Pt(x, y).In(clip) {
		return
	}
	target := s.At(x, y)
	if target == c {
		return
	}
	stack := []image.Point{{x, y}}
	for len(stack) > 0 {
		p := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if !p.In(clip) || s.At(p.X, p.Y) != target {
			continue
		}
		s.Set(p.X, p.Y, c)
		stack = append(stack,
			image.Pt(p.X+1, p.Y), image.Pt(p.X-1, p.Y),
			image.Pt(p.X, p.Y+1), image.Pt(p.X, p.Y-1))
	}
}

// Fill sets every pixel in r to c.
func Fill(s *Sheet, r image.Rectangle, c uint8) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			s.Set(x, y, c)
		}
	}
}

// FlipH mirrors r horizontally.
func FlipH(s *Sheet, r image.Rectangle) {
	transform(s, r, func(src []uint8, w, h, x, y int) uint8 { return src[y*w+(w-1-x)] })
}

// FlipV mirrors r vertically.
func FlipV(s *Sheet, r image.Rectangle) {
	transform(s, r, func(src []uint8, w, h, x, y int) uint8 { return src[(h-1-y)*w+x] })
}

// Rotate rotates a square r 90 degrees clockwise.
func Rotate(s *Sheet, r image.Rectangle) {
	transform(s, r, func(src []uint8, w, h, x, y int) uint8 { return src[(h-1-x)*w+y] })
}

func transform(s *Sheet, r image.Rectangle, f func(src []uint8, w, h, x, y int) uint8) {
	src := s.Read(r)
	w, h := r.Dx(), r.Dy()
	dst := make([]uint8, len(src))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dst[y*w+x] = f(src, w, h, x, y)
		}
	}
	s.Write(r, dst)
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
