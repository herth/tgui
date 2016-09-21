package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/herth/tgui"
	"github.com/herth/tgui/games"

	"github.com/gdamore/tcell"
)

// sorting

type byName []os.FileInfo

func (s byName) Len() int {
	return len(s)
}
func (s byName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byName) Less(i, j int) bool {
	return s[i].Name() < s[j].Name()
}

type bySize []os.FileInfo

func (s bySize) Len() int {
	return len(s)
}
func (s bySize) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s bySize) Less(i, j int) bool {
	return s[i].Size() < s[j].Size()
}

type byTime []os.FileInfo

func (s byTime) Len() int {
	return len(s)
}
func (s byTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byTime) Less(i, j int) bool {
	return s[i].ModTime().Before(s[j].ModTime())
}

type Info struct {
	name, size, time string
}

// print a file info as an one line output
func prInfo(info os.FileInfo) Info {
	s := prdec(info.Size())
	return Info{name: info.Name(), size: s, time: info.ModTime().Format("Jan 02 15:04")}
}

// print i as a dotted decimal 23.456.789
func prdec(i int64) string {
	txt := fmt.Sprintf("%d", i)
	l := len(txt)
	dots := (l - 1) / 3

	bytes := make([]byte, l+dots)
	for pos := 0; pos < l; pos++ {
		bytes[l+dots-pos-1] = txt[l-pos-1]
		if (pos+1)%3 == 0 && dots > 0 {
			dots--
			bytes[l+dots-pos-1] = byte("."[0])
		}
	}
	return string(bytes)
}

func lsPath(path string, multiple, size, time, reverse bool) (result []Info) {
	f, err := os.Lstat(path)
	if err != nil {
		msg := err.Error()
		// strip leading "lstat "
		if len(msg) > 6 {
			if msg[:5] == "lstat" {
				msg = msg[6:]
			}
		}
		fmt.Printf("%s\n", msg)
	} else {
		if f.IsDir() {
			// if multiple {
			// 	fmt.Printf("%s:\n", path)
			// }
			info, _ := ioutil.ReadDir(path)
			switch {
			case size:
				if reverse {
					sort.Sort(sort.Reverse(bySize(info)))
				} else {
					sort.Sort(bySize(info))
				}
			case time:
				if reverse {
					sort.Sort(sort.Reverse(byTime(info)))
				} else {
					sort.Sort(byTime(info))
				}
			default:
				if reverse {
					sort.Sort(sort.Reverse(byName(info)))
				} else {
					sort.Sort(byName(info))
				}
			}
			for _, i := range info {
				result = append(result, prInfo(i))
			}
		} else {
			result = append(result, prInfo(f))
		}
	}
	return
}

type MC struct {
	tgui.DecoratedWin
	Dir      string
	CW2, CW3 int
	Files    []Info
	bySize   bool
	byTime   bool
	reverse  bool
}

func NewMC(a *tgui.App, dir string, x, y int) *MC {
	w := 80
	h := 30
	st := tcell.StyleDefault.Background(tcell.ColorNavy).Foreground(tcell.ColorYellow)
	win := &MC{Dir: dir,
		Files: lsPath(dir, true, false, false, false),
		CW2:   16,
		CW3:   15,
		DecoratedWin: tgui.DecoratedWin{Title: " " + dir + " ",
			SimpleWin: tgui.SimpleWin{App: a, Box: tgui.Box{x, y, x + w, y + h}, Fill: ' ', Style: st}}}
	a.AddWindow(win)
	return win
}

func (w *MC) MMove(x, y int) {
	s := w.Screen
	width := w.X1 - w.X0
	x0 := w.X0
	y0 := w.Y0
	cw2 := w.CW2
	cw3 := w.CW3
	st := tcell.StyleDefault.Background(tcell.ColorYellow).Foreground(tcell.ColorNavy)
	if y == w.Y0+1 {
		if x > w.X0+1 && x < w.X1-w.CW2-w.CW3 { // name column
			tgui.Print(s, x0+1+(width-cw2-cw3)/2-2, y0+1, "Name", st)
			tgui.Print(s, w.X1-cw2-cw3+cw2/2-2, y0+1, "Size", w.Style)
			tgui.Print(s, w.X1-cw3+cw3/2-2, y0+1, "MTime", w.Style)
		} else if x > w.X1-w.CW2-w.CW3 && x < w.X1-w.CW3 {
			tgui.Print(s, x0+1+(width-cw2-cw3)/2-2, y0+1, "Name", w.Style)
			tgui.Print(s, w.X1-cw2-cw3+cw2/2-2, y0+1, "Size", st)
			tgui.Print(s, w.X1-cw3+cw3/2-2, y0+1, "MTime", w.Style)
		} else if x > w.X1-w.CW3 && x < w.X1 {
			tgui.Print(s, x0+1+(width-cw2-cw3)/2-2, y0+1, "Name", w.Style)
			tgui.Print(s, w.X1-cw2-cw3+cw2/2-2, y0+1, "Size", w.Style)
			tgui.Print(s, w.X1-cw3+cw3/2-2, y0+1, "MTime", st)
		}
	} else {
		tgui.Print(s, x0+1+(width-cw2-cw3)/2-2, y0+1, "Name", w.Style)
		tgui.Print(s, w.X1-cw2-cw3+cw2/2-2, y0+1, "Size", w.Style)
		tgui.Print(s, w.X1-cw3+cw3/2-2, y0+1, "MTime", w.Style)
	}
	s.Show()
}

func (w *MC) Click(x, y int) {
	if y == w.Y0+1 {
		if x > w.X0+1 && x < w.X1-w.CW2-w.CW3 { // name column
			if w.bySize || w.byTime {
				w.bySize = false
				w.byTime = false
				w.reverse = false
			} else {
				w.reverse = !w.reverse
			}
		} else if x > w.X1-w.CW2-w.CW3 && x < w.X1-w.CW3 {
			if w.bySize {
				w.reverse = !w.reverse
			} else {
				w.bySize = true
				w.byTime = false
				w.reverse = false
			}
		} else if x > w.X1-w.CW3 && x < w.X1 {
			if w.byTime {
				w.reverse = !w.reverse
			} else {
				w.bySize = false
				w.byTime = true
				w.reverse = false
			}
		}
		w.Files = lsPath(w.Dir, true, w.bySize, w.byTime, w.reverse)
		w.Draw()
		w.App.Screen.Show()
	}
}

func (w *MC) Draw() {
	s := w.Screen
	x0 := w.X0
	y0 := w.Y0
	cw2 := w.CW2
	cw3 := w.CW3
	width := w.X1 - w.X0
	height := w.Y1 - w.Y0
	// if w.reverse {
	// 	w.Style = tcell.StyleDefault.Background(tcell.ColorYellow).Foreground(tcell.ColorNavy)
	// } else {
	// 	w.Style = tcell.StyleDefault.Background(tcell.ColorNavy).Foreground(tcell.ColorYellow)
	// }

	w.DecoratedWin.Draw()

	tgui.Print(s, x0+1+(width-cw2-cw3)/2-2, y0+1, "Name", w.Style)
	tgui.Print(s, w.X1-cw2-cw3+cw2/2-2, y0+1, "Size", w.Style)
	tgui.Print(s, w.X1-cw3+cw3/2-2, y0+1, "MTime", w.Style)

	tgui.Print(s, w.X1-cw2-cw3, y0, string(tcell.RuneTTee), w.Style)
	tgui.Print(s, w.X1-cw2-cw3, w.Y1, string(tcell.RuneBTee), w.Style)
	for y := 1; y < height; y++ {
		tgui.Print(s, w.X1-cw2-cw3, y0+y, string(tcell.RuneVLine), w.Style)
	}
	tgui.Print(s, w.X1-cw3, y0, string(tcell.RuneTTee), w.Style)
	tgui.Print(s, w.X1-cw3, w.Y1, string(tcell.RuneBTee), w.Style)
	for y := 1; y < height; y++ {
		tgui.Print(s, w.X1-cw3, y0+y, string(tcell.RuneVLine), w.Style)
	}

	n := len(w.Files)
	if n > height {
		n = height
	}

	for y := 2; y < n; y++ {
		tgui.PrintM(s, x0+2, y0+y, x0+width-cw2-cw3, w.Files[y].name, w.Style)
		si := w.Files[y].size
		tgui.PrintM(s, w.X1-cw3-len(si)-1, y0+y, x0+width-cw3, si, w.Style)
		tgui.PrintM(s, w.X1-cw3+2, y0+y, x0+width, w.Files[y].time, w.Style)
	}

	if w.reverse {
		tgui.Print(s, x0+1, y0+1, "", w.Style)
	} else {
		tgui.Print(s, x0+1, y0+1, "", w.Style)
	}
}

func main() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	screen, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = screen.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	ncolor := screen.Colors()
	defer func() {
		screen.Fini()
		fmt.Printf("%d Colors\n", ncolor)
		fmt.Printf("Goodbye\n")
	}()

	screen.SetStyle(tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorBlue))

	_, h := screen.Size()
	screen.ShowCursor(0, h-1)
	screen.EnableMouse()
	app := &tgui.App{Screen: screen, Background: tcell.StyleDefault.Background(tcell.NewRGBColor(80, 80, 80))}

	mb := &tgui.MenuBar{Screen: screen}
	//mb.AddMenu(Menu{Name: "≡"})
	mb.AddMenu(tgui.Menu{Name: "="})
	mb.AddMenu(tgui.Menu{Name: "File"})
	mb.AddMenu(tgui.Menu{Name: "Edit"})
	mb.AddMenu(tgui.Menu{Name: "Search"})
	mb.AddMenu(tgui.Menu{Name: "Run"})
	mb.AddMenu(tgui.Menu{Name: "Exit"})
	app.MB = mb

	w1 := app.NewTextWin("ls", 31, 12, 20, 10, "", tcell.StyleDefault.Background(tcell.ColorNavy).Foreground(tcell.ColorWhite))
	w1.Tail = true
	go tgui.Clock(screen)
	screen.Show()

	wt := games.NewTetris(app)
	NewMC(app, "/etc/", 30, 10)
	ls := NewMC(app, "/Users/herth/", 0, 1)

	app.Draw()

	//wt.Draw()
	_ = wt
	_ = ls
	app.RuneHandler = func(r rune) {
		fmt.Fprintf(w1, "%c", r)
	}

	app.KeyHandler = func(k tcell.Key) {
		switch k {
		case tcell.KeyEnter:
			fmt.Fprintf(w1, "\n")
		}
	}

	tgui.EventHandle(app)
	//screen.Fini()

}
