package view

import (
    "fmt"
    "image"
    "image/draw"
    "image/color"
    //"fmt"
    "github.com/usedbytes/ui"
    //"github.com/usedbytes/fonts"
)


type View struct {
    *ui.Widget

    Name string
    
    children []ui.Drawable
    
    canvas *image.Paletted
}

func NewView(p *ui.Widget, name string) *View {
    view := new(View)
    view.Widget = ui.NewWidget(p)
    view.Name = name

    fmt.Printf("View %s created\n", view.Name)
    // make zero image
    return view
}

func (v *View) Update() {
    
    if (v.canvas == nil) || (v.Bounds().Size() != v.canvas.Bounds().Size()) {
        p := color.Palette{v.Widget.Background, v.Widget.Foreground}
        v.canvas = image.NewPaletted(image.Rectangle{image.Point{0,0}, v.Bounds().Size()}, p)
    }
}

func (v *View) Draw(to draw.Image) {

    if (v.IsDirty()) {
        v.Update();
    }

    draw.Draw(v.canvas, v.Bounds(), &image.Uniform{v.Widget.Background}, image.ZP, draw.Src)
    for _, c := range v.children {
        c.Draw(v.canvas)
    }

    draw.Draw(to, v.Bounds(), v.canvas, image.ZP, draw.Src)
}

func (v *View) IsDirty() bool {
    
    if v.Dirty {
        return true
    }
    
    for _, c := range v.children {
        if (c.IsDirty()) {
            return true
        }
    }
    
    return false
}

func (v *View) AddChild(c ui.Drawable) {
    v.children = append(v.children, c)
}
