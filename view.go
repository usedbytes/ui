package ui

import (
    "image"
    "image/draw"
    "image/color"
)


type View struct {
    *Widget

    Name string
    
    children []*Drawable
    Damage image.Rectangle
    
    canvas *image.Paletted
}

func NewView(p *Widget, name string) *View {
    view := new(View)
    view.Widget = NewWidget(p)
    view.Name = name

    return view
}

func (v *View) Update() {
    
    if (v.canvas == nil) || (v.Bounds().Size() != v.canvas.Bounds().Size()) {
        p := color.Palette{v.Widget.Background, v.Widget.Foreground}
        v.canvas = image.NewPaletted(image.Rectangle{image.Point{0,0}, v.Bounds().Size()}, p)
    }
}

func (v *View) Draw(to draw.Image) image.Rectangle {
    damaged := image.ZR
    if (v.IsDirty()) {
        v.Update();
    }

    //draw.Draw(v.canvas, image.Rectangle{image.ZP, v.Bounds().Size()}, &image.Uniform{v.Widget.Background}, image.ZP, draw.Src)
    
    for _, c := range v.children {
        if (*c).IsVisible() {
            hasDamage := v.Damage != image.ZR
            if ( (hasDamage && v.Damage.Overlaps((*c).Bounds())) || 
                 (*c).IsDirty()) {
                (*c).MakeDirty()
                d := (*c).Draw(v.canvas)
                if (damaged == image.ZR) {
                    damaged = d
                } else {
                    damaged = damaged.Union(d)
                }
            }
        }
    }

    v.Damage = image.ZR
    if (v.IsVisible()) {
        draw.Draw(to, v.Bounds(), v.canvas, image.ZP, draw.Src)
        return damaged
    }
    return image.ZR
}

func (v *View) IsDirty() bool {
    
    if v.Dirty {
        return true
    }
    
    for _, c := range v.children {
        if ((*c).IsDirty()) {
            return true
        }
    }
    
    return false
}

func (v *View) AddChild(c Drawable) {
    v.children = append(v.children, &c)
    v.Damage = v.Damage.Union(c.Bounds())
}

func (v *View) MakeDirty() {
    v.Dirty = true
}
