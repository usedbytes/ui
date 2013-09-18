package label

import (
    "fmt"
    "image"
    "image/color"
    "testing"
    "time"
    "github.com/usedbytes/s4548"
    "github.com/usedbytes/fonts"
)

//import _ "image/png"

func TestAlign(t *testing.T) {
    //screen := s4548.NewS4548("/dev/s4548-0")
    screen := s4548.NewS4548(s4548.GetS4548EnvPath())
    font := fonts.NewFontFromFile("/home/kernelcode/tiny_font.fnt")
    fmt.Println("Font Loaded")

    short:= NewLabel(nil, font)
    short.AutoWidth = false
    short.AutoHeight = false
    short.SetWidth(101)
    short.SetHeight(40)
    short.VAlign = Middle
    short.HAlign = Centre
    short.Text = "Short"
    short.Background = color.Black
    short.Foreground = color.White
    
    for i := 0; i < 3; i++ {
        switch i {
        case 0:
            short.VAlign = Top
        case 1:
            short.VAlign = Middle
        case 2:
            short.VAlign = Bottom
        } 
        for j := 0; j < 3; j++ {
            switch j {
            case 0:
                short.HAlign = Left
            case 1:
                short.HAlign = Centre
            case 2:
                short.HAlign = Right
            } 
            time.Sleep(500 * time.Millisecond)
            
            short.Update()
            short.Draw(screen)
            screen.Scanout()
        }
    }

    screen.Close()
    time.Sleep(1000 * time.Millisecond)

/*
    for {v
        l.Update()
        l.Draw(screen)
    
        l.Text = "Today I learned that writing word wrapping for the second time was still around 100,000,000 times harder than I was expecting"
        l.Widget.AutoWidth = true
        //l.SetWidth(101)
        l.SetPos(image.Point{0,0})
        l.Update()
        l.Draw(screen)
    
        screen.Scanout()
        time.Sleep(300 * time.Millisecond)
    }
*/
}

func TestHorizontal(t *testing.T) {
    //screen := s4548.NewS4548("/dev/s4548-0")
    screen := s4548.NewS4548(s4548.GetS4548EnvPath())
    font := fonts.NewFontFromFile("/home/kernelcode/tiny_font.fnt")
    fmt.Println("Font Loaded")


    longLine := NewLabel(nil, font)
    longLine.AutoWidth = false
    longLine.AutoHeight = true
    longLine.Wrap = false
    longLine.Scroll = false
    longLine.SetWidth(80)
    longLine.SetPos(image.Point{10, 8})
    longLine.VAlign = Middle
    longLine.HAlign = Centre
    longLine.Text = "Lazee ft. Neverstore - Hold On (Matrix & Futurebound Remix)"

    longLine.Update()
    longLine.Draw(screen)
    screen.Scanout()

    time.Sleep(1000 * time.Millisecond)

    longLine.Scroll = true
    longLine.Active = true
    for i := 0; i < 30; i++ {
        longLine.Update()
        longLine.Draw(screen)
        screen.Scanout()

        time.Sleep(200 * time.Millisecond)
    }
    screen.Close()

    time.Sleep(1000 * time.Millisecond)
}

func TestVertical(t *testing.T) {
    //screen := s4548.NewS4548("/dev/s4548-0")
    screen := s4548.NewS4548(s4548.GetS4548EnvPath())
    font := fonts.NewFontFromFile("/home/kernelcode/tiny_font.fnt")
    fmt.Println("Font Loaded")


    multiLine := NewLabel(nil, font)
    multiLine.AutoWidth = false
    multiLine.AutoHeight = false
    multiLine.Wrap = true
    multiLine.Scroll = false
    multiLine.SetWidth(101)
    multiLine.SetHeight(40)
    multiLine.SetPos(image.Point{0, 0})
    multiLine.VAlign = Middle
    multiLine.HAlign = Centre
    multiLine.Text = "Once upon a time in a galaxy far away... STAR WARS The evil galactic empire have taken over hundreds of systems"

    multiLine.Update()
    multiLine.Draw(screen)
    screen.Scanout()

    time.Sleep(5000 * time.Millisecond)

    multiLine.Scroll = true
    multiLine.Active = true
    for i := 0; i < 20; i++ {
        multiLine.Update()
        multiLine.Draw(screen)
        screen.Scanout()

        time.Sleep(300 * time.Millisecond)
    }

    screen.Close()

    time.Sleep(5000 * time.Millisecond)
}
