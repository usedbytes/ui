package ui

import (
    "image"
    "image/color"
    "image/draw"
    //"github.com/gvalkov/golang-evdev"
)


type Interactable interface {
    // Handle(e *evdev.InputEvent)
    GiveFocus()
    RemoveFocus()
}


type Drawable interface {
    Draw(to draw.Image)
    IsDirty() bool
    Update()
}

type Widget struct {
    *image.Rectangle
    AutoHeight, AutoWidth bool
    Dirty bool
    
    Parent *Widget
    //Children []*Widget

    Foreground, Background color.Color
}

func NewWidget(p *Widget) *Widget {
    w := new(Widget)
    r := image.Rect(0,0,0,0)
    w.Rectangle = &r
    w.AutoHeight = true
    w.AutoWidth = true
    w.Dirty = true
    w.Parent = p
    w.Foreground = color.Black
    w.Background = color.White
    return w
}

func (w *Widget) SetWidth(width int) {
    w.Rectangle.Max.X = w.Rectangle.Min.X + width
    w.Dirty = true
}

func (w *Widget) SetHeight(height int) {
    w.Rectangle.Max.Y = w.Rectangle.Min.Y + height
    w.Dirty = true
}

func (w *Widget) Bounds() (image.Rectangle) {
    return *w.Rectangle
}

func (w *Widget) SetPos(topLeft image.Point) {
    r := image.Rectangle{topLeft, topLeft.Add(w.Size())}
    w.Rectangle = &r
}
