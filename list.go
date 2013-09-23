package ui

import (
    "fmt"
    "image"
    "image/draw"
    "image/color"
    "github.com/usedbytes/fonts"
)

type List struct {
    *Widget
    font *fonts.Font

    Title, title string
    titleLabel *Label
        
    items, onscreen []*listItem
    
    selected, canFit int

    labelCache []*Label
    icons []*image.Paletted
    canvas *image.Paletted
    
    drawables []*Drawable
}

type listItem  struct {
    index int
    text string
    iconIndex int
    //label *label.Label
    //canvas *image.Paletted
    
    action int
}

func (l *listItem) String() string {
    return fmt.Sprintf("[%v: %v%s - %v]", l.index, l.iconIndex, l.text, l.action)
}

func (l *List) newListItem(i int, t string, icon int, act int) *listItem {
    li := new(listItem)
    li.index = i
    li.text = t
    li.iconIndex = icon
    li.action = act
    
    return li
}

func (l *List) AddItem(text string) {
    li := l.newListItem(len(l.items), text, 0, 0)
    l.items = append(l.items, li)
}

func NewList(p *Widget, f *fonts.Font) *List {
    list := new(List)
    list.Widget = NewWidget(p)
    list.font = f
    
    return list
}

func (l *List) SetCapacity(capacity int) {
    if (capacity < len(l.items)) {
        return
    } else {
        l.items = append(l.items, make([]*listItem, capacity - len(l.items))...)
    }
}


func (l *List) Update() {

    l.title = l.Title
    fmt.Println(l.title)
    
    /*
    if (l.Widget.AutoWidth) {
        l.SetWidth(l.font.Width(l.title))
    }
    */

    numItems := len(l.items)
    h := numItems
    if l.title != "" {
        h++
    }

    if (l.Widget.AutoHeight) {
        l.SetHeight(h * l.font.Height())
    }
    
    activeHeight := l.Bounds().Dy()
    if (l.title != "") {
        activeHeight -= l.font.Height()
    }
    l.canFit = activeHeight / l.font.Height();
    if (l.font.Height() * l.canFit) < activeHeight {
        l.canFit++
    }
    
    fmt.Printf("numItems: %v, canFit: %v\n", numItems, l.canFit)
    
    fmt.Println("Selected:", l.selected)
    if (l.canFit >= numItems) {
        fmt.Println("Can fit all")
        l.onscreen = l.items
    } else {
        if (l.selected < l.canFit / 2) {
            fmt.Println("At Top")
            l.onscreen = l.items[0:l.canFit]
            fmt.Println(l.items)
        } else if (l.selected >= (numItems - (l.canFit / 2))) {
            fmt.Println("At Bottom")
            l.onscreen = l.items[numItems - l.canFit:numItems]
            fmt.Println(l.items)
        } else {
            l.onscreen = make([]*listItem, l.canFit)
            for i := 0; i < l.canFit; i++ {
                el := (l.selected - (l.canFit / 2) + i) % numItems
                l.onscreen[i] = l.items[el]
                fmt.Printf("%v: %v\t", i, el)
            }
            fmt.Printf("\n")
            fmt.Println(l.items)
        }
    }
    fmt.Println(l.onscreen)
    /*
    for _, i := range l.items {
        i.SetWidth(l.Bounds().Dx())
        i.Update()
    }
    */
    
    l.makeGraphics()
    
    l.Dirty = false
}

func (l *List) IsDirty() bool {    
    
    return true
    
    if l.Dirty {
        return true
    }
            
    if (l.title != l.Title) {
        l.Dirty = true
    }
        
    return l.Dirty
}

func (l *List) makeGraphics() {

    if l.title != "" {
        if (l.titleLabel == nil) {
            l.titleLabel = NewLabel(l.Widget, l.font);
            l.titleLabel.HAlign = Left
            l.titleLabel.Foreground = l.Foreground
            l.titleLabel.Background = l.Background
        }
        l.titleLabel.SetWidth(l.Bounds().Dx())
        l.titleLabel.SetHeight(l.font.Height())        
        l.titleLabel.Text = l.title
    }
    
    if (len(l.labelCache) != l.canFit) {
        l.labelCache = make([]*Label, l.canFit) 
        for i := 0; i < l.canFit; i++ {
            lbl := NewLabel(l.Widget, l.font)
            l.labelCache[i] = lbl
            lbl.AutoWidth = false
            lbl.Foreground = l.Foreground
            lbl.Background = l.Background
            ypos := i * l.font.Height()
            if (l.title != "") {
                ypos += l.font.Height()
            }
            lbl.SetPos(image.Point{l.font.Height(), ypos})
            lbl.SetWidth(l.Bounds().Dx() - l.font.Height())
        }
    }

    if (l.canvas == nil) || (l.Bounds().Size() != l.canvas.Bounds().Size()) {    
        p := color.Palette{l.Widget.Background, l.Widget.Foreground}
        l.canvas = image.NewPaletted(image.Rectangle{image.Point{0,0}, l.Bounds().Size()}, p)
    }
}

func (l *List) Draw(to draw.Image) {
    
    if (l.IsDirty()) {
        l.Update();
    }
    
    if (l.title != "") {
        l.titleLabel.Draw(l.canvas)
    }
    
    for i := 0; i < l.canFit; i++ {
        l.labelCache[i].Text = l.onscreen[i].text
        if (l.selected == l.onscreen[i].index) {
            l.labelCache[i].Invert = true
        } else {
            l.labelCache[i].Invert = false
        }
                // l.labelCache[i].Foreground = l.Background
                // l.labelCache[i].Background = l.Foreground
                //l.labelCache[i].Dirty = true
            /*
            } else if (l.labelCache[i].Foreground != l.Foreground) {
            l.labelCache[i].Foreground = l.Foreground
            l.labelCache[i].Background = l.Background
            l.labelCache[i].Dirty = true
            */
        l.labelCache[i].Draw(l.canvas)
    }
    //l.wasSelected = l.selected
    
    if (l.IsVisible()) {    
        draw.Draw(to, l.Bounds(), l.canvas, image.ZP, draw.Src)
    }

}