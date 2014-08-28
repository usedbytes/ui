package ui

import (
    "fmt"
    "image"
    "image/draw"
    "image/color"
    "github.com/usedbytes/fonts"
    "github.com/usedbytes/input"
    "time"
)

type List struct {
    *Widget
    font, iconFont *fonts.Font

    Title, title string
    halt bool
    titleLabel *Label
        
    items, onscreen []*ListItem
    
    Selected, selected, canFit int

    labelCache, iconCache []*Label
    canvas *image.Paletted
    
    drawables []*Drawable
}

type actionT func(...interface{})

type ListItem  struct {
    index int
    Text string
    IconIndex int   
    Action actionT
    Tag interface{}
}

func (l *ListItem) String() string {
    return fmt.Sprintf("[%v: %v%s - %v]", 
                l.index, l.IconIndex, l.Text, l.Action)
}

func (l *List) newListItem(i int, t string, icon int, act actionT,
    tag interface{}) *ListItem {
    li := new(ListItem)
    li.index = i
    li.Text = t
    li.IconIndex = icon
    li.Action = act
    li.Tag = tag
    return li
}

func (l *List) AddItem(text string, icon int, act actionT, tag interface{}) {
    li := l.newListItem(len(l.items), text, icon, act, tag)
    l.items = append(l.items, li)
}

func NewList(p *Widget, f *fonts.Font, iF *fonts.Font) *List {
    list := new(List)
    list.Widget = NewWidget(p)
    list.font = f
    list.iconFont = iF
    list.onscreen = make([]*ListItem, 0)
    
    return list
}

func (l *List) SetCapacity(capacity int) {
    if (capacity < len(l.items)) {
        return
    } else {
        l.items = append(l.items, make([]*ListItem, capacity - len(l.items))...)
    }
}

func (l *List) NumItems() int {
    return len(l.items)
}

func (l *List) Update() {

    if (l.title != l.Title) {
        l.title = l.Title
        l.makeTitle()
    }

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
    canFit := activeHeight / l.font.Height();
    if (l.font.Height() * canFit) < activeHeight {
        canFit++
    }
    
    //fmt.Printf("numItems: %v, canFit: %v\n", numItems, l.canFit)
    
    //fmt.Println("Selected:", l.Selected)
    if (canFit >= numItems) {
        //fmt.Println("Can fit all")
        l.onscreen = l.items
    } else {
        if (l.Selected < canFit / 2) {
            //fmt.Println("At Top")
            l.onscreen = l.items[0:canFit]
            //fmt.Println(l.items)
        } else if (l.Selected >= (numItems - (canFit / 2))) {
            //fmt.Println("At Bottom")
            l.onscreen = l.items[numItems - canFit:numItems]
            //fmt.Println(l.items)
        } else {
            l.onscreen = make([]*ListItem, canFit)
            for i := 0; i < canFit; i++ {
                el := (l.Selected - (canFit / 2) + i) % numItems
                l.onscreen[i] = l.items[el]
                //fmt.Printf("%v: %v\t", i, el)
            }
            //fmt.Printf("\n")
            //fmt.Println(l.items)
        }
    }
    //fmt.Println(l.onscreen)
    /*
    for _, i := range l.items {
        i.SetWidth(l.Bounds().Dx())
        i.Update()
    }
    */
    if (canFit != l.canFit) { 
        l.canFit = canFit
        l.makeGraphics()
    }

    l.Dirty = false
}

func (l *List) IsDirty() bool {    
    
    
    if l.Dirty {
        return true
    }
            
    if (l.title != l.Title) {
        l.Dirty = true
    }
        
    if (l.selected != l.Selected) {
        l.selected = l.Selected
        return true
    }

    for _, c := range l.labelCache {
        if c.IsDirty() {
            return true
        }
    }

    for _, i := range l.iconCache {
        if i.IsDirty() {
            return true
        }
    }

    return l.Dirty
}

func (l *List) makeTitle() {

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

}

func (l *List) makeGraphics() {

    if (len(l.iconCache) != l.canFit) {
        l.iconCache = make([]*Label, len(l.onscreen)) 
        for i := 0; i < len(l.onscreen); i++ {
            lbl := NewLabel(l.Widget, l.iconFont)
            l.iconCache[i] = lbl
            lbl.AutoWidth = false
            lbl.Foreground = l.Foreground
            lbl.Background = l.Background
            lbl.HAlign = Centre
            lbl.VAlign = Middle
            ypos := i * l.font.Height()
            if (l.title != "") {
                ypos += l.font.Height()
            }
            lbl.SetPos(image.Point{0, ypos})
            lbl.SetWidth(l.font.Height())
        }
    }

    if (len(l.labelCache) != l.canFit) {
        l.labelCache = make([]*Label, len(l.onscreen)) 
        for i := 0; i < len(l.onscreen); i++ {
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

func (l *List) Draw(to draw.Image) image.Rectangle {

    if (l.IsDirty()) {
        l.Update();
    }
    
    if (l.title != "") {
        l.titleLabel.Draw(l.canvas)
    }
    
    //for i := 0; i < l.canFit && i < l.NumItems(); i++ {
    for i := 0; i < len(l.onscreen); i++ {
        l.labelCache[i].Text = l.onscreen[i].Text
        l.iconCache[i].Text = string(l.onscreen[i].IconIndex)
        if (l.Selected == l.onscreen[i].index) {
            l.labelCache[i].Invert = true
            l.iconCache[i].Invert = true
        } else {
            l.labelCache[i].Invert = false
            l.iconCache[i].Invert = false
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
        l.iconCache[i].Draw(l.canvas)
    }
    //l.wasSelected = l.Selected
    
    if (l.IsVisible()) {    
        draw.Draw(to, l.Bounds(), l.canvas, image.ZP, draw.Src)
        return l.Bounds()
    }

    return image.ZR

}

func (l* List) Item(i int) *ListItem {
    return l.items[i]
}

func (l *List) HandleInput(key rune) bool {
    switch key {
        case input.KEY_UP, input.KEY_SCROLLDOWN:
            if (l.Selected > 0) {
                l.Selected--
                l.halt = false
                if (l.Selected == 0) {
                    l.halt = true;
                    go l.haltTimer()
                }
            } else if (l.Selected == 0) {
                if (!l.halt) {
                    l.Selected = l.NumItems() - 1
                }
            }
            return true
        case input.KEY_DOWN, input.KEY_SCROLLUP:
            if (l.Selected < l.NumItems() - 1) {
                l.Selected++
                l.halt = false
                if (l.Selected == l.NumItems() - 1) {
                    l.halt = true;
                    go l.haltTimer()
                }
            } else if (l.Selected == l.NumItems() - 1) {
                if (!l.halt) {
                    l.Selected = 0
                }
            }
            return true
        case input.KEY_ENTER:
            li := l.Item(l.Selected)
            if (li.Action != nil) {
                li.Action(li)
            }
            return true
    }
    return false
}

func (l *List) haltTimer() {
    time.Sleep(1200 * time.Millisecond)
    l.halt = false;
}

func (l *List) MakeDirty() {
    l.Dirty = true
}
