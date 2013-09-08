package view

import (
    "fmt"
    "image"
    //"image/draw"
    //"image/color"
    //"github.com/usedbytes/ui"
    "github.com/usedbytes/fonts"
    "github.com/usedbytes/ui/label"
    "github.com/usedbytes/ui/progbar"
    "github.com/usedbytes/s4548"
    "testing"
    "time"
    "math/rand"
)

func TestView(t* testing.T) {
    screen := s4548.NewS4548(s4548.GetS4548EnvPath())
    font := fonts.NewFontFromFile("/home/kernelcode/tiny_font.fnt")
    fmt.Println("Font Loaded")

    r := rand.New(rand.NewSource(99))

    view := NewView(nil, "Test View");
    view.SetWidth(101)
    view.SetHeight(40)
    
    lbl1 := label.NewLabel(nil, font)
    lbl1.AutoWidth = false
    lbl1.AutoHeight = true
    lbl1.SetWidth(101)
    lbl1.SetHeight(40)
    lbl1.VAlign = label.Middle
    lbl1.HAlign = label.Centre
    lbl1.Text = "This is a big long label test which should be scrolling with any luck"
    lbl1.Scroll = true
    
    view.AddChild(lbl1)
    
    bar := progbar.NewProgressBar(nil)
    bar.SetWidth(101)
    bar.SetHeight(8)
    bar.SetPos(image.Point{0, 32})
    bar.Progress = 33
    
    view.AddChild(bar)
    
    for {
        time.Sleep(500 * time.Millisecond)
        bar.Progress = r.Intn(101)
        view.Draw(screen)
        screen.Scanout()
    }
}
