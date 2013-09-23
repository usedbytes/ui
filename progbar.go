package ui

import (
    "github.com/usedbytes/basic2d"
    "image"
    "image/draw"
    "image/color"
    "fmt"
)

type Orientation int

const (
    Vertical Orientation = iota
    Horizontal
)

const hIndent = 1
const vIndent = 1
const borderWidth = 1

type ProgressBar struct {
    *Widget
    Min, Max int
    Progress int
    Direction Orientation
    graphics []*image.Paletted
    canvas *image.Paletted
}

func NewProgressBar(p *Widget) (*ProgressBar) {
    bar := new(ProgressBar)
    bar.Widget = NewWidget(p)
    bar.Min = 0
    bar.Max = 100
    bar.Progress = 50
    bar.Direction = Horizontal

    return bar
}

func (b *ProgressBar) MakeGraphics() {
    p := color.Palette{b.Widget.Background, b.Widget.Foreground}
    b.graphics = make([]*image.Paletted, 1)
    b.graphics[0] = image.NewPaletted(image.Rectangle{image.ZP, b.Bounds().Size()}, p)

    r := b.graphics[0].Bounds()
    border := image.Rectangle{r.Min.Add(image.Point{hIndent, vIndent}), r.Max.Sub(image.Point{hIndent, vIndent})}
    basic2d.Box(b.graphics[0], border, borderWidth, b.Widget.Foreground)
    
    if (b.canvas == nil) || (b.Bounds().Size() != b.canvas.Bounds().Size()) {    
        b.canvas = image.NewPaletted(image.Rectangle{image.ZP, b.Bounds().Size()}, p)
    }
}

func (b *ProgressBar) IsDirty() bool {
    return b.Dirty
}

func (b *ProgressBar) Update() {
    b.MakeGraphics();
}

func (b *ProgressBar) Draw(to draw.Image) {
    var fillWidth, fillHeight int
    var dp image.Point
    var dr image.Rectangle
    
    if (b.IsDirty()) {
        b.Update()
    }
    
    g := b.graphics[0]
    if (debug) {
        fmt.Printf("b.Bounds(): %v, canvas.Bounds(): %v, g.Bounds(): %v\n", b.Bounds(), b.canvas.Bounds(), g.Bounds())
    }
    draw.Draw(b.canvas, b.canvas.Bounds(), g, g.Bounds().Min, draw.Src)

    mrange := b.Max - b.Min
    percent := (b.Progress  * 100) / mrange
    if (percent > 100) {
        percent = 100
    } else if (percent < 0) {
        percent = 0
    }
    activeWidth := b.Widget.Dx() - (2 * hIndent) - (2 * borderWidth)
    activeHeight := b.Widget.Dy() - (2 * vIndent) - (2 * borderWidth)
    
    
    if (b.Direction == Horizontal) {
        fillWidth = (activeWidth * percent) / 100
        fillHeight = activeHeight
        dp = image.Point{hIndent + borderWidth, vIndent + borderWidth}
        dr = image.Rectangle{dp, dp.Add(image.Point{fillWidth, fillHeight})}
    } else {
        fillHeight = (activeHeight * percent) / 100
        fillWidth = activeWidth
        dp = image.Point{hIndent + borderWidth + (activeWidth - fillWidth), vIndent + borderWidth + (activeHeight - fillHeight)}
        dr = image.Rectangle{dp, dp.Add(image.Point{fillWidth, fillHeight})}
    }
    
    if (debug) {
        fmt.Printf("percent: %v, activeWidth: %v, activeHeight: %v, fillWidth: %v, fillHeight: %v, dr: %v\n", percent, activeWidth, activeHeight, fillWidth, fillHeight, dr)
    }
    draw.Draw(b.canvas, dr, &image.Uniform{b.Widget.Foreground}, image.ZP, draw.Src)

    if (b.IsVisible()) {
        draw.Draw(to, b.Bounds(), b.canvas, image.ZP, draw.Src)
    }
}
