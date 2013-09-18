package list

import (
    "fmt"
    //"image/draw"
    //"image/color"
    //"github.com/usedbytes/ui"
    "github.com/usedbytes/fonts"
    "github.com/usedbytes/s4548"
    "testing"
    "time"
)

func TestList(t* testing.T) {
    screen := s4548.NewS4548(s4548.GetS4548EnvPath())
    font := fonts.NewFontFromFile("/home/kernelcode/tiny_font.fnt")
    fmt.Println("Font Loaded")

    list := NewList(nil, font);
    list.Title = "A test list:"
    list.SetWidth(101)
    list.SetHeight(40)
    fmt.Println("Set size:", list.Bounds())
    list.AutoHeight = false
    list.selected = 1
    
    list.AddItem("Zero")
    list.AddItem("One")
    list.AddItem("Two")
    list.AddItem("Three")
    list.AddItem("Four")
    list.AddItem("Five")
    list.AddItem("Six")
    time.Sleep(500 * time.Millisecond)
    
    for {
        for i := 0; i < 7; i++ {
            time.Sleep(500 * time.Millisecond)
            list.selected = i
            list.Draw(screen)
            screen.Scanout()
            
        }
    }
    
    
    
}
