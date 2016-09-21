package games

import (
	"herth/tgui"

	"github.com/gdamore/tcell"
)

//(defparameter *colors* #(:black :blue :red :green :brown :magenta :cyan :grey :yellow))
var color = []tcell.Color{tcell.ColorBlack, tcell.ColorBlue, tcell.ColorRed, tcell.ColorGreen, tcell.ColorMaroon, tcell.ColorFuchsia, tcell.ColorAqua, tcell.ColorGrey, tcell.ColorYellow}

type Stone [][]int

var Stones = []Stone{
	Stone{
		{1, 1},
		{1, 1}},
	Stone{
		{2},
		{2},
		{2},
		{2}},
	Stone{
		{3, 3, 0},
		{0, 3, 3}},
	Stone{
		{0, 4, 4},
		{4, 4, 0}},
	Stone{
		{5, 5},
		{5, 0},
		{5, 0}},
	Stone{
		{6, 6},
		{0, 6},
		{0, 6}},
	Stone{
		{0, 7, 0},
		{7, 7, 7}}}

type Tetris struct {
	tgui.SimpleWin
	Width  int
	Height int
	speed  int
	field  [][]int
}

func NewTetris(app *tgui.App) *Tetris {
	u := 10
	l := 68
	w := &Tetris{
		SimpleWin: tgui.SimpleWin{
			App:   app,
			Box:   tgui.Box{X0: l, Y0: u, X1: l + 20 + 11, Y1: u + 20 + 2},
			Style: tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorYellow)},
		Width:  10,
		Height: 20,
	}
	w.field = make([][]int, w.Height)
	for i := 0; i < w.Height; i++ {
		w.field[i] = make([]int, w.Width)
	}
	for j := 10; j < w.Height; j++ {
		for i := 0; i < w.Width; i++ {
			w.field[j][i] = (i+j*12)%7 + 1
		}
	}
	app.AddWindow(w)
	return w
}

func (w *Tetris) DrawBox(col, px, py int) {
	x0 := w.X0 + 1
	y0 := w.Y0 + 2
	st := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(color[col])
	w.Screen.SetContent(x0+2*px, y0+py, '█', nil, st)
	w.Screen.SetContent(x0+2*px+1, y0+py, '▌', nil, st)
}

func (w *Tetris) DrawStone(stone Stone, px, py int) {
	for y := 0; y < len(stone); y++ {
		for x := 0; x < len(stone[y]); x++ {
			w.DrawBox(stone[y][x], px+x, py+y)
		}
	}
}

func (w *Tetris) Draw() {
	w.SimpleWin.Draw()
	st := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorNavy)
	w.DrawRect(w.X0, w.Y0+1, w.X0+w.Width*2+1, w.Y0+1+w.Height+1, st)
	w.DrawRect(w.X0+w.Width*2+2, w.Y0+1, w.X0+w.Width*2+11, w.Y0+6, st)

	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			w.DrawBox(w.field[y][x], x, y)
		}
	}
	w.DrawStone(Stones[3], 11, 1)
	// for i := 0; i < len(Stones); i++ {
	// 	w.DrawStone(Stones[i], 0, i*4)
	// }
}
