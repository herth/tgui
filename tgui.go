package tgui

import (
	"math/rand"
	"strings"
	"time"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
)

// not used
func makebox(s tcell.Screen) {
	w, h := s.Size()

	if w == 0 || h == 0 {
		return
	}

	glyphs := []rune{'@', '#', '&', '*', '=', '%', 'Z', 'A'}

	lx := rand.Int() % w
	ly := rand.Int() % h
	lw := rand.Int() % (w - lx)
	lh := rand.Int() % (h - ly)
	st := tcell.StyleDefault
	gl := ' '
	if s.Colors() > 256 {
		rgb := tcell.NewHexColor(int32(rand.Int() & 0xffffff))
		st = st.Background(rgb)
	} else if s.Colors() > 1 {
		st = st.Background(tcell.Color(rand.Int() % s.Colors()))
	} else {
		st = st.Reverse(rand.Int()%2 == 0)
		gl = glyphs[rand.Int()%len(glyphs)]
	}

	for row := 0; row < lh; row++ {
		for col := 0; col < lw; col++ {
			s.SetCell(lx+col, ly+row, st, gl)
		}
	}
	s.Show()
}

// print the text to the (x,y) coordinates of the screen.
func Print(s tcell.Screen, x, y int, text string, st tcell.Style) {
	for _, c := range text {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, st)
		x += w
	}
}

// Print the text up to a maximum X coordinate.
func PrintM(s tcell.Screen, x, y int, maxX int, text string, st tcell.Style) {
	for _, c := range text {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		if maxX >= 0 && (x+w) <= maxX {
			s.SetContent(x, y, c, comb, st)
			x += w
		} else {
			return
		}
	}
}

type MenuBar struct {
	menus  [](*Menu)
	Screen tcell.Screen
}

type Menu struct {
	Name string
	X    int
	Y    int
	W    int
	hl   bool
}

func (m *MenuBar) AddMenu(menu Menu) {
	m.menus = append(m.menus, &menu)
	x := 0
	y := 0
	for _, menu := range m.menus {
		menu.X = x
		menu.Y = y
		menu.W = runewidth.StringWidth(menu.Name) + 2
		x += menu.W
	}
}

type App struct {
	MB          *MenuBar
	Screen      tcell.Screen
	quit        bool
	windows     []Window
	Background  tcell.Style
	RuneHandler func(r rune)
	KeyHandler  func(k tcell.Key)
}

func (a *App) AddWindow(w Window) {
	a.windows = append(a.windows, w)
}

type Box struct {
	X0, Y0, X1, Y1 int
}

func (b Box) Inside(x, y int) bool {
	if x >= b.X0 && x <= b.X1 &&
		y >= b.Y0 && y <= b.Y1 {
		return true
	} else {
		return false
	}
}

type Window interface {
	Draw()
	Size() (w int, h int)
	GetBox() Box
	SetBox(b Box)
	Click(x, y int)
	MMove(x, y int)
}

func DrawBorder(s tcell.Screen, b Box, style tcell.Style) {
	x0 := b.X0
	x1 := b.X1
	y0 := b.Y0
	y1 := b.Y1
	s.SetContent(x0, y0, tcell.RuneULCorner, nil, style)
	s.SetContent(x0, y1, tcell.RuneLLCorner, nil, style)
	s.SetContent(x1, y1, tcell.RuneLRCorner, nil, style)
	s.SetContent(x1, y0, tcell.RuneURCorner, nil, style)
	for x := x0 + 1; x < x1; x++ {
		s.SetContent(x, y0, tcell.RuneHLine, nil, style)
		s.SetContent(x, y1, tcell.RuneHLine, nil, style)
	}
	for y := y0 + 1; y < y1; y++ {
		s.SetContent(x0, y, tcell.RuneVLine, nil, style)
		s.SetContent(x1, y, tcell.RuneVLine, nil, style)
	}
}

func DrawBox(s tcell.Screen, b Box, r rune, style tcell.Style) {
	for y := b.Y0; y <= b.Y1; y++ {
		for x := b.X0; x <= b.X1; x++ {
			s.SetContent(x, y, r, nil, style)
		}
	}
}

type SimpleWin struct {
	*App
	Box
	Fill  rune
	Style tcell.Style
}

func (w *SimpleWin) Draw() {
	DrawBox(w.App.Screen, w.Box, w.Fill, w.Style)
}

func (w *SimpleWin) Size() (int, int) {
	return w.X1 - w.X0, w.Y1 - w.Y0
}

func (w *SimpleWin) GetBox() Box {
	return w.Box
}

func (w *SimpleWin) SetBox(b Box) {
	w.Box = b
}

func (w *SimpleWin) Click(x, y int) {
}

func (w *SimpleWin) MMove(x, y int) {
}

type DecoratedWin struct {
	SimpleWin
	Title string
}

func (w *DecoratedWin) Draw() {
	w.SimpleWin.Draw()
	style := w.Style.Foreground(tcell.ColorYellow)
	s := w.Screen
	DrawBorder(s, w.Box, style)
	PrintM(s, w.X0+2, w.Y0, w.X1, w.Title, style)
	for x := w.X0 + 1; x <= w.X1+1; x++ {
		y := w.Y1 + 1
		r1, r2, st, _ := s.GetContent(x, y)
		s.SetContent(x, y, r1, r2, st.Background(tcell.ColorBlack))
	}
	for y := w.Y0 + 1; y <= w.Y1+1; y++ {
		x := w.X1 + 1
		r1, r2, st, _ := s.GetContent(x, y)
		s.SetContent(x, y, r1, r2, st.Background(tcell.ColorBlack))
	}

}

type TextWin struct {
	DecoratedWin
	Text string
	Tail bool
}

func (w *TextWin) Draw() {
	w.DecoratedWin.Draw()
	x := w.X0 + 1
	h := w.Y1 - w.Y0 - 1
	text := strings.Split(w.Text, "\n")
	if w.Tail && len(text) > h {
		text = text[len(text)-h-1:]
	}
	for i, line := range text {
		y := w.Y0 + 1 + i
		if x < w.X1-1 && y < w.Y1 {
			PrintM(w.Screen, x, y, w.X1, line, w.Style)
		}
	}
}

func (w *TextWin) Write(b []byte) (n int, err error) {
	w.Text = w.Text + string(b)
	w.App.Draw()
	return len(b), nil

}

func (a *App) NewTextWin(title string, x, y, w, h int, text string, st tcell.Style) *TextWin {
	win := &TextWin{Text: text,
		DecoratedWin: DecoratedWin{Title: title,
			SimpleWin: SimpleWin{App: a, Box: Box{x, y, x + w, y + h}, Fill: ' ', Style: st}}}
	a.AddWindow(win)
	return win
}

func (a *App) Quit() {
	a.quit = true
}

func (m *MenuBar) Draw() {
	s := m.Screen
	w, _ := s.Size()
	st := tcell.StyleDefault.
		Background(tcell.ColorLightGray).
		Foreground(tcell.ColorBlack)
	r := st.Foreground(tcell.ColorRed)
	for i := 0; i < w; i++ {
		s.SetContent(i, 0, ' ', nil, st)
	}
	for _, item := range m.menus {
		name := item.Name
		st := st
		r := r
		if item.hl {
			st = st.Background(tcell.ColorLime)
			r = r.Background(tcell.ColorLime)
		}
		Print(s, item.X, item.Y, " ", st)
		Print(s, item.X+1, item.Y, string(name[0]), r)
		Print(s, item.X+2, 0, name[1:], st)
		Print(s, item.X+item.W-1, item.Y, " ", st)
	}
}

func (a *App) Draw() {

	w, h := a.Screen.Size()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			a.Screen.SetContent(x, y, ' ', nil, a.Background)
		}
	}

	a.MB.Draw()

	for _, win := range a.windows {
		(win).Draw()
	}
}

func highlightMenus(app *App, x, y int) {
	for _, menu := range app.MB.menus {
		if x >= menu.X &&
			x < menu.X+menu.W &&
			y == menu.Y {
			menu.hl = true
		} else {
			menu.hl = false
		}
	}
	app.MB.Draw()
	showTime(app.Screen)
	app.Screen.Show()
}

func (app *App) Raise(w Window) {
	var temp Window
	for i, win := range app.windows {
		if win == w {
			temp = w
		} else if temp != nil {
			app.windows[i-1] = app.windows[i]
		}
	}
	if temp != nil {
		app.windows[len(app.windows)-1] = temp
		temp.Draw()
		app.Screen.Show()
	}
}

func (app *App) ButtonEvent(nr, x, y int) {
	for _, menu := range app.MB.menus {
		if x >= menu.X &&
			x < menu.X+menu.W &&
			y == menu.Y {
			if menu.Name == "Exit" {
				app.Quit()
			}
		}
	}

	raised := false
	var clicked Window
wloop:
	for i := len(app.windows) - 1; i >= 0; i-- {
		win := app.windows[i]
		if win.GetBox().Inside(x, y) {
			clicked = win
			if i < len(app.windows)-1 {
				app.Raise(win)
				raised = true
			}
			break wloop
		}
	}
	if clicked != nil && !raised {
		clicked.Click(x, y)
	}
}

func (app *App) FindWin(x, y int) Window {
	for i := len(app.windows) - 1; i >= 0; i-- {
		win := app.windows[i]
		if win.GetBox().Inside(x, y) {
			return win
		}
	}
	return nil
}

func (app *App) MMove(x, y int) {
	for i := len(app.windows) - 1; i >= 0; i-- {
		win := app.windows[i]
		if win.GetBox().Inside(x, y) {
			win.MMove(x, y)
			return
		}
	}
}

// handle a pressed mouse move from pxy to xy
func (app *App) Drag(px, py, x, y int) {
	w := app.FindWin(px, py)
	if w != nil {
		b := w.GetBox()
		dx := x - px
		dy := y - py
		if b.Y0 == py { //  title bar drag
			b.X0 += dx
			b.X1 += dx
			b.Y0 += dy
			b.Y1 += dy
			w.SetBox(b)
			app.Draw()
			return
		} else if b.X1 == px && b.Y1 == py { // lower right corner -> resize
			b.X1 += dx
			b.Y1 += dy
			if b.X1 > b.X0+2 && b.Y1 > b.Y0+1 {
				w.SetBox(b)
				app.Draw()
				return
			}
		}
	}
}

func showTime(s tcell.Screen) {
	t := time.Now()
	w, _ := s.Size()
	st := tcell.StyleDefault.
		Background(tcell.ColorLightGray).
		Foreground(tcell.ColorBlack)
	Print(s, w-9, 0, t.Format("15:04:05"), st)

	s.Show()
}

func Clock(s tcell.Screen) {
	showTime(s)
	for {
		select {
		case <-time.After(time.Millisecond * 500):
			showTime(s)
		}
	}
}

func EventHandle(app *App) {
	screen := app.Screen
	pressed := false
	px := 0
	py := 0
	for {
		if app.quit {
			return
		}
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape: // , tcell.KeyEnter:
				//close(quit)
				//	ok = true
				return
			case tcell.KeyCtrlL:
				screen.Sync()
			case tcell.KeyCtrlQ:
				return
			case tcell.KeyRune:
				rune := ev.Rune()
				switch rune {
				default:
					if app.RuneHandler != nil {
						app.RuneHandler(rune)
					}
				}
			default:
				if app.KeyHandler != nil {
					app.KeyHandler(ev.Key())
				}
			}
		case *tcell.EventResize:
			_, h := screen.Size()
			screen.Clear()
			screen.ShowCursor(0, h-1)
			app.Draw()
			screen.Sync()
		case *tcell.EventMouse:
			x, y := ev.Position()
			button := ev.Buttons()
			highlightMenus(app, x, y)

			// if button&tcell.WheelUp != 0 {
			// }
			// if button&tcell.WheelDown != 0 {
			// }
			// if button&tcell.WheelLeft != 0 {
			// }
			// if button&tcell.WheelRight != 0 {
			// }
			// // Only buttons, not wheel events
			button &= tcell.ButtonMask(0xff)
			justPressed := false
			if button != tcell.ButtonNone && !pressed { // click press
				pressed = true
				justPressed = true
				px, py = x, y
			}
			switch ev.Buttons() {
			case tcell.ButtonNone:
				app.MMove(x, y)
				if pressed { // click release
					pressed = false
				}
			case tcell.Button1:
				if justPressed {
					app.ButtonEvent(1, x, y)
				} else {
					app.Drag(px, py, x, y)
					px, py = x, y
				}
			}
		}
	}
}
