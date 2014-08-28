package ui

import (
    "fmt"
    "image"
    "time"
    "image/draw"
    "image/color"
    //"fmt"
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

const debug bool = false

type Label struct {
    *Widget
    font *fonts.Font

    Text string
    Scroll, Wrap, Active, Invert bool

    text string
    lines []string
    scroll, wrap, active, invert bool
    scrollSpeed, scrollPos int
    drawMode int
    VAlign Alignment
    HAlign Alignment
    graphics []*image.Paletted
    canvas *image.Paletted
    
    lastTime time.Time
    
}

const WRAP_GAP = 12

func NewLabel(p *Widget, f *fonts.Font) *Label {
    label := new(Label)
    label.Widget = NewWidget(p)
    label.font = f
    label.VAlign = Top
    label.HAlign = Left
    label.scrollSpeed = -15
    label.lastTime = time.Now()
    
    
    label.makeGraphics()
    return label
}


func (l *Label) Font() *fonts.Font {
    return l.font
}

func (l *Label) Update() {

    width := l.Bounds().Dx()
    height := l.Bounds().Dy()

    textChanged, wrapChanged := (l.text != l.Text), (l.wrap != l.Wrap)
    l.text, l.scroll, l.wrap, l.invert = l.Text, l.Scroll, l.Wrap, l.Invert
   
    if (textChanged && l.Widget.AutoWidth) {
        l.SetWidth(l.font.Width(l.text));
    }

    if (textChanged || wrapChanged) {
        if (l.wrap) {
            l.lines = l.font.WrapText(l.text, l.Dx());
        } else {
            l.lines = []string{l.text}
        }
    }
    
    if ((textChanged || wrapChanged) && l.Widget.AutoHeight) {
        l.SetHeight(len(l.lines) * l.font.Height());
    }

    if (!l.Scroll) {
        l.scrollPos = 0;
    }

    if (textChanged || wrapChanged || (width != l.Bounds().Dx()) ||
        (height != l.Bounds().Dy())) {
        l.makeGraphics();
    }
    
    if (l.Active && !l.active) {
        l.lastTime = time.Now()
    }
    l.active = l.Active
    
    l.Dirty = false
}

func (l *Label) IsDirty() bool {    
    
    if l.Dirty {
        fmt.Println("Dirty cos dirty")
        return true
    }
            
    if (l.scroll != l.Scroll) || (l.wrap != l.Wrap) || (l.text != l.Text) || 
        (l.invert != l.Invert) {
        fmt.Println("Dirty cos changed")
        l.Dirty = true
    } else if l.Scroll {
        fmt.Println("Dirty cos scroll")
        l.Dirty = true;
    }
    
    return l.Dirty
}

func (l *Label) makeGraphics() {
    if (debug) {
        fmt.Println("label.makeGraphics");
    }
    l.graphics = make([]*image.Paletted, len(l.lines)*2)

    for i, line := range l.lines {
        l.graphics[i] = l.font.MakeWordColor(line, l.Widget.Background, l.Widget.Foreground)
        //*(l.graphics[i + len(l.lines)]) = *(l.graphics[i])
        l.graphics[i + len(l.lines)] = l.font.MakeWordColor(line, l.Widget.Background, l.Widget.Foreground)
    }

    if (l.canvas == nil) || (l.Bounds().Size() != l.canvas.Bounds().Size()) {    
        //p := color.Palette{l.Widget.Background, l.Widget.Foreground}
        p := color.Palette{l.Widget.Foreground, l.Widget.Background}
        l.canvas = image.NewPaletted(image.Rectangle{image.Point{0,0}, l.Bounds().Size()}, p)
    }
    
    //p := color.Palette{l.Widget.Background, l.Widget.Foreground}
    //l.canvas = image.NewPaletted(image.Rectangle{image.Point{0,0}, l.Bounds().Size()}, p)
}

/*
func (l *Label) InvertColors() {
    for _, im := range l.graphics {
        im.Palette[0], im.Palette[1] = im.Palette[1], im.Palette[0]
    }
    l.canvas.Palette[0], l.canvas.Palette[1] = l.canvas.Palette[1], l.canvas.Palette[0]
    l.inverted = !l.inverted
    l.Dirty = true
}

func (l *Label) Inverted() bool {
    return l.inverted
}
*/

func (l *Label) ResetScroll() {
    l.scrollPos = 0
}

func (l *Label) Draw(to draw.Image) image.Rectangle {
    
    if (l.IsDirty()) {
        l.Update();
        
        if (debug) {
            fmt.Printf("Label bounds: %v\n", l.Bounds())
        }
        for i, _ := range l.lines {
            l.verticalPosition(i)
            l.horizontalPosition(i)
        }

        bg := l.Widget.Background
        if l.Invert {
            bg = l.Widget.Foreground
        }

        //draw.Draw(l.canvas, image.Rectangle{image.ZP, l.Bounds().Size()}, &image.Uniform{bg}, image.ZP, draw.Src)
        if (l.IsVisible()) {
            draw.Draw(to, l.Bounds(), &image.Uniform{bg}, image.ZP, draw.Src)
            for _, g := range l.graphics {
                if l.Invert {
                    ip := color.Palette{l.Widget.Foreground, l.Widget.Background}
                    g.Palette = ip
                } else {
                    p := color.Palette{l.Widget.Background, l.Widget.Foreground}
                    g.Palette = p
                }
                //topLeft := l.Bounds().Min.Add(image.Point{0, l.font.Height() * i})
                //dest := image.Rectangle{topLeft, topLeft.Add(g.Bounds().Max)}
                if (debug) {
                    fmt.Printf("\tSprite bounds: %v\n", g.Bounds())
                }
                draw.Draw(to, l.Bounds(), g, image.ZP, draw.Src)
            }
        }
        
        if (l.active && l.scroll) {
            now := time.Now()
            nanoseconds := now.Sub(l.lastTime).Nanoseconds()
            l.scrollPos += int(int64(l.scrollSpeed) * nanoseconds / 1e9);
        }
        l.lastTime = time.Now()

        return l.Bounds()
    }
    return image.ZR
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

func (l *Label) MakeDirty() {
    l.Dirty = true
}
