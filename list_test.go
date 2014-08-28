package ui

import (
    "fmt"
    //"image/draw"
    //"image/color"
    "github.com/usedbytes/fonts"
    "github.com/usedbytes/s4548"
    "testing"
    "time"
)

func TestList(t* testing.T) {
    screen := s4548.NewS4548(s4548.GetS4548EnvPath())
    font := fonts.NewFontFromFile("/home/kernelcode/tiny_font.fnt")
    iconFont := fonts.NewFontFromFile("/home/kernelcode/icon_font.fnt")
    fmt.Println("Font Loaded")

    list := NewList(nil, font, iconFont);
    list.Title = "A test list:"
    list.SetWidth(101)
    list.SetHeight(40)
    fmt.Println("Set size:", list.Bounds())
    list.AutoHeight = false
    list.Selected = 0
    
    list.AddItem("Zero", 0, nil, nil)
    /*
    list.AddItem("One")
    list.AddItem("Two")
    list.AddItem("Three")
    list.AddItem("Four")
    list.AddItem("Five")
    list.AddItem("Six")
    */
    time.Sleep(500 * time.Millisecond)
    
    for j := 0; j < 5; j++ {
        for i := 0; i < 1; i++ {
            time.Sleep(500 * time.Millisecond)
            list.Selected = i
            list.Draw(screen)
            screen.Scanout()
            
        }
    }
    
    
    
}
