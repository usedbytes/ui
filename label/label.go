package label

import (
    "fmt"
    "image"
    "image/draw"
    "image/color"
    //"fmt"
    "github.com/usedbytes/ui"
    "github.com/usedbytes/fonts"
)

type Alignment int

const (
    Top Alignment = iota
    Middle
    Bottom
    Centre
    Left
    Right
)


type Label struct {
    *ui.Widget
    font *fonts.Font

    Text string
    Scroll, Wrap bool

    text string
    lines []string
    scroll, wrap bool
    scrollSpeed, scrollPos int
    drawMode int
    VAlign Alignment
    HAlign Alignment
    graphics []*image.Paletted
    canvas *image.Paletted
}

const WRAP_GAP = 12

func NewLabel(p *ui.Widget, f *fonts.Font) *Label {
    label := new(Label)
    label.Widget = ui.NewWidget(p)
    label.font = f
    label.VAlign = Top
    label.HAlign = Left
    label.scrollSpeed = -5
    // make zero image
    return label
}

//TODO: Could do optimisation so this can be called automatically in Draw.
// figure out what has changed and thus what needs updating
func (l *Label) Update() {

    l.text = l.Text
    l.scroll = l.Scroll
    l.wrap = l.Wrap

    if (l.Widget.AutoWidth) {
        l.SetWidth(l.font.Width(l.text));
    }

    if (l.wrap) {
        l.lines = l.font.WrapText(l.text, l.Dx());
    } else {
        l.lines = []string{l.text}
    }

    if (l.Widget.AutoHeight) {
        l.SetHeight(len(l.lines) * l.font.Height());
    }

    if (!l.scroll) {
        l.scrollPos = 0;
    }

    l.makeGraphics();
}

func (l *Label) makeGraphics() {

    l.graphics = make([]*image.Paletted, len(l.lines)*2)
    for i, line := range l.lines {
        l.graphics[i] = l.font.MakeWordColor(line, l.Widget.Background, l.Widget.Foreground)
        l.graphics[i + len(l.lines)] = l.font.MakeWordColor(line, l.Widget.Background, l.Widget.Foreground)
    }

    //p := color.Palette{color.White, color.Black}
    p := color.Palette{l.Widget.Background, l.Widget.Foreground}
    l.canvas = image.NewPaletted(image.Rectangle{image.Point{0,0}, l.Bounds().Size()}, p)
}

func (l *Label) Draw(to draw.Image) {
    fmt.Printf("Label bounds: %v\n", l.Bounds())
    for i, _ := range l.lines {
        l.verticalPosition(i)
        l.horizontalPosition(i)
    }

    for _, g := range l.graphics {
        //topLeft := l.Bounds().Min.Add(image.Point{0, l.font.Height() * i})
        //dest := image.Rectangle{topLeft, topLeft.Add(g.Bounds().Max)}
        fmt.Printf("\tSprite bounds: %v\n", g.Bounds())
        draw.Draw(l.canvas, g.Bounds(), g, g.Bounds().Min, draw.Src)
    }

    draw.Draw(to, l.Bounds(), l.canvas, image.ZP, draw.Src)
    l.scrollPos += l.scrollSpeed;
}

func (l *Label) verticalPosition(line int) {
    textHeight := len(l.lines) * l.font.Height()
    var topOffset, tail int

    switch l.VAlign {
    case Top:
        topOffset = 0
    case Middle:
        topOffset = (l.Bounds().Dy() - textHeight) / 2
    case Bottom:
        topOffset = (l.Bounds().Dy() - textHeight)
    }

    totalHeight := textHeight + WRAP_GAP

    if ((textHeight > l.Dy()) && l.scroll) {
        if ( l.scrollPos > 0 ) {
            l.scrollPos -= totalHeight 
            //l.scrollPos %= ( totalHeight + WRAP_GAP )
        }
        tail = l.scrollPos + totalHeight;
        if (tail <= WRAP_GAP) {
            // There is only one copy on screen
            l.scrollPos = tail;
        }
        topLeft := image.Point{0, topOffset + l.scrollPos + (line * l.font.Height())}
        l.graphics[line].Rect = image.Rectangle{topLeft, topLeft.Add(l.graphics[line].Bounds().Size())}
        l.graphics[line + len(l.lines)].Rect = l.graphics[line].Rect.Add(image.Point{0, totalHeight})
    } else {
        topLeft := image.Point{0, topOffset + (line * l.font.Height())}
        l.graphics[line].Rect = image.Rectangle{topLeft, topLeft.Add(l.graphics[line].Bounds().Size())}
        l.graphics[line + len(l.lines)].Rect = l.graphics[line].Rect
    }
}

func (l *Label) horizontalPosition(line int) {
    var tail, leftOffset int

    textWidth := l.font.Width(l.lines[line])
    totalWidth := textWidth + WRAP_GAP

    switch l.HAlign {
    case Left:
        leftOffset = 0
    case Centre:
        leftOffset = (l.Bounds().Dx() - textWidth) / 2
    case Right:
        leftOffset = (l.Bounds().Dx() - textWidth)
    }

    if ((textWidth > l.Dx()) && l.scroll) {
        // Alias it, so that tail is always to the right
        // scrollpos > 0 can only happen when scroll == true
        if (l.scrollPos > 0) {
            l.scrollPos -= totalWidth;
        }
        tail = l.scrollPos + totalWidth;
        if (tail <= WRAP_GAP) {
            l.scrollPos = tail
        }
        topLeft := image.Point{leftOffset + l.scrollPos, l.graphics[line].Rect.Min.Y}
        l.graphics[line].Rect = image.Rectangle{topLeft, topLeft.Add(l.graphics[line].Bounds().Size())} 
        l.graphics[line + len(l.lines)].Rect = l.graphics[line].Rect.Add(image.Point{totalWidth, 0}) 
    } else {
        topLeft := image.Point{leftOffset, l.graphics[line].Rect.Min.Y}
        l.graphics[line].Rect = image.Rectangle{topLeft, topLeft.Add(l.graphics[line].Bounds().Size())} 
        topLeft = image.Point{leftOffset, l.graphics[line + len(l.lines)].Rect.Min.Y}
        l.graphics[line + len(l.lines)].Rect = image.Rectangle{topLeft, topLeft.Add(l.graphics[line + len(l.lines)].Bounds().Size())} 
    }

}
