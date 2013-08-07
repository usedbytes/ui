package progbar

import (
    "github.com/usedbytes/ui"
    "github.com/usedbytes/ui/basic2d"
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

type ProgressBar struct {
    *ui.Widget
    Min, Max int
    Progress int
    Direction Orientation
    graphics []*image.Paletted
    canvas *image.Paletted
}

func NewProgressBar(p *ui.Widget) (*ProgressBar) {
    bar := new(ProgressBar)
    bar.Widget = ui.NewWidget(p)
    bar.Min = 0
    bar.Max = 100
    bar.Progress = 50
    bar.Direction = Horizontal

    return bar
}

func (b *ProgressBar) MakeGraphics() {
    p := color.Palette{b.Widget.Background, b.Widget.Foreground}
    b.graphics = make([]*image.Paletted, 1)
    b.graphics[0] = image.NewPaletted(image.Rectangle{image.Point{0,0}, b.Bounds().Size()}, p)

    r := b.graphics[0].Bounds()
    border := image.Rectangle{r.Min.Add(image.Point{1, 1}), r.Max.Sub(image.Point{1, 1})}
    basic2d.Box(b.graphics[0], border, 1, b.Widget.Foreground)

    b.canvas = image.NewPaletted(image.Rectangle{image.Point{0,0}, b.Bounds().Size()}, p)
}

func (b *ProgressBar) Draw(to draw.Image) {
    g := b.graphics[0]
    fmt.Printf("b.Bounds(): %v, canvas.Bounds(): %v, g.Bounds(): %v\n", b.Bounds(), b.canvas.Bounds(), g.Bounds())
    draw.Draw(b.canvas, b.canvas.Bounds(), g, g.Bounds().Min, draw.Src)

    //mrange := b.Max - b.Min

    draw.Draw(to, b.Bounds(), b.canvas, image.ZP, draw.Src)
}
