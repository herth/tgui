package tgui

import (
	"math/rand"
	"strings"
	"time"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
)

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

// func print(s tcell.Screen, x, y int, text string, st tcell.Style) {
// 	for i, c := range text {
// 		s.SetContent(x+i, y, c, nil, st)
// 	}
// }

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
		menu.W = len(menu.Name) + 2 // fixme
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

type SimpleWin struct {
	*App
	Box
	Fill  rune
	Style tcell.Style
}

func (w *SimpleWin) Draw() {
	// st := tcell.StyleDefault.
	// 	Background(tcell.ColorLightGray).
	// 	Foreground(tcell.ColorBlack)
	for y := w.Y0; y <= w.Y1; y++ {
		for x := w.X0; x <= w.X1; x++ {
			w.App.Screen.SetContent(x, y, w.Fill, nil, w.Style)
		}
	}
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

func (w *SimpleWin) DrawRect(x0, y0, x1, y1 int, style tcell.Style) {
	s := w.Screen
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

func (w *DecoratedWin) Draw() {
	w.SimpleWin.Draw()
	style := w.Style.Foreground(tcell.ColorYellow)
	s := w.Screen
	w.DrawRect(w.X0, w.Y0, w.X1, w.Y1, style)
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
		//s.SetContent(i, h-2, ' ', nil, st)
	}
	// Print(s, 1, 0, "≡", r)
	//x := 3
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
	//Print(s, 1, 0, "= File  Edit  Search  Run ", st)
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
	// ok := false
	// defer func() {
	// 	if !ok {
	// 		screen.Fini()
	// 	}
	// }()
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
				// case 'a':
				// 	makebox(screen)
				// 	screen.Show()
				// case 'c':
				// 	screen.Clear()
				// 	screen.Show()
				// 	//case 'q', 'Q':
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

			// for i := uint(0); i < 8; i++ {
			// 	if int(button)&(1<<i) != 0 {
			// 		bstr += fmt.Sprintf(" Button%d", i+1)
			// 	}
			// }
			// if button&tcell.WheelUp != 0 {
			// 	bstr += " WheelUp"
			// }
			// if button&tcell.WheelDown != 0 {
			// 	bstr += " WheelDown"
			// }
			// if button&tcell.WheelLeft != 0 {
			// 	bstr += " WheelLeft"
			// }
			// if button&tcell.WheelRight != 0 {
			// 	bstr += " WheelRight"
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
			//ch := '*'
			// screen.ShowCursor(x, y)
			// screen.Show()
			// if button != tcell.ButtonNone && ox < 0 {
			// 	ox, oy = x, y
			// }
			// switch ev.Buttons() {
			// case tcell.ButtonNone:
			// 	if ox >= 0 {
			// 		bg := tcell.Color((lchar - '0') * 2)
			// 		drawBox(s, ox, oy, x, y,
			// 			up.Background(bg),
			// 			lchar)
			// 		ox, oy = -1, -1
			// 		bx, by = -1, -1
			// 	}
			// case tcell.Button1:
			// 	ch = '1'
			// case tcell.Button2:
			// 	ch = '2'
			// case tcell.Button3:
			// 	ch = '3'
			// case tcell.Button4:
			// 	ch = '4'
			// case tcell.Button5:
			// 	ch = '5'
			// case tcell.Button6:
			// 	ch = '6'
			// case tcell.Button7:
			// 	ch = '7'
			// case tcell.Button8:
			// 	ch = '8'
			// default:
			// 	ch = '*'

			// }
			// if button != tcell.ButtonNone {
			// 	bx, by = x, y
			// }
			// lchar = ch

		}
	}
}

// func main() {
// 	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
// 	screen, e := tcell.NewScreen()
// 	if e != nil {
// 		fmt.Fprintf(os.Stderr, "%v\n", e)
// 		os.Exit(1)
// 	}
// 	if e = screen.Init(); e != nil {
// 		fmt.Fprintf(os.Stderr, "%v\n", e)
// 		os.Exit(1)
// 	}

// 	ncolor := screen.Colors()
// 	defer func() {
// 		screen.Fini()
// 		fmt.Printf("%d Colors\n", ncolor)
// 		fmt.Printf("Goodbye\n")
// 	}()

// 	screen.SetStyle(tcell.StyleDefault.
// 		Background(tcell.ColorBlack).
// 		Foreground(tcell.ColorBlue))

// 	_, h := screen.Size()
// 	screen.ShowCursor(0, h-1)
// 	screen.EnableMouse()
// 	app := &App{Screen: screen, Background: tcell.StyleDefault.Background(tcell.NewRGBColor(80, 80, 80))}

// 	mb := &MenuBar{Screen: screen}
// 	//mb.AddMenu(Menu{Name: "≡"})
// 	mb.AddMenu(Menu{Name: "="})
// 	mb.AddMenu(Menu{Name: "File"})
// 	mb.AddMenu(Menu{Name: "Edit"})
// 	mb.AddMenu(Menu{Name: "Search"})
// 	mb.AddMenu(Menu{Name: "Run"})
// 	mb.AddMenu(Menu{Name: "Exit"})
// 	app.MB = mb

// 	app.AddWindow(&SimpleWin{App: app, Box: Box{2, 2, 40, 30}, Style: tcell.StyleDefault.Background(tcell.ColorLightGray)})
// 	app.AddWindow(&SimpleWin{App: app, Box: Box{10, 12, 30, 22}, Fill: ' ', Style: tcell.StyleDefault.Background(tcell.ColorBlue)})
// 	app.AddWindow(&DecoratedWin{Title: "foo", SimpleWin: SimpleWin{App: app, Box: Box{31, 12, 50, 22}, Fill: ' ', Style: tcell.StyleDefault.Background(tcell.ColorNavy)}})

// 	app.AddWindow(&TextWin{Text: "foo\n\nbar\n\nbaz", DecoratedWin: DecoratedWin{Title: "foo", SimpleWin: SimpleWin{App: app, Box: Box{31, 12, 50, 22}, Fill: ' ', Style: tcell.StyleDefault.Background(tcell.ColorNavy).Foreground(tcell.ColorYellow)}}})

// 	app.NewTextWin("foo", 31, 12, 20, 10, "foo\n\nbar\n\nbaz", tcell.StyleDefault.Background(tcell.ColorNavy).Foreground(tcell.ColorWhite))
// 	app.NewTextWin("hallodri", 40, 40, 20, 10, "foo\n\nbar\n\nbaz", tcell.StyleDefault.Background(tcell.ColorMaroon).Foreground(tcell.ColorYellow))

// 	for i := 0; i < 16; i++ {
// 		app.AddWindow(&SimpleWin{App: app, Box: Box{1, i + 20, 30, i + 20}, Style: tcell.StyleDefault.Background(tcell.Color(i))})
// 	}
// 	app.AddWindow(&SimpleWin{App: app, Box: Box{51, 12, 70, 22}, Fill: ' ', Style: tcell.StyleDefault.Background(tcell.NewRGBColor(60, 60, 255))})
// 	app.Draw()

// 	go clock(screen)
// 	screen.Show()

// 	eventHandle(app)
// 	//screen.Fini()

// }
