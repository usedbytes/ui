package ui

import (
    "github.com/usedbytes/fonts"
    "image"
    "image/draw"
    "fmt"
)

const (
    STATE_PLAYING = iota
    STATE_PAUSED
    STATE_NONE
)

const (
    ICON_PLAY = string(1)
    ICON_PAUSE = string(2)
    ICON_REPEAT = string(3)
    ICON_SHUFFLE = string(4)
    ICON_BLANK = string(5)
)

type StatusBar struct {
    view *View
    icons, tiny *fonts.Font
    stateLabel, shuffleLabel, repeatLabel, tracksLabel *Label 
    
    State int
    Repeat, Shuffle bool
    TrackNum, Tracks int
}

func NewStatusBar(p *Widget) (*StatusBar) {
    sb := new(StatusBar)
    sb.view = NewView(p, "Status Bar")
    sb.view.AutoWidth = false
    sb.view.AutoHeight = false
    sb.view.SetWidth(101)
    sb.view.SetHeight(8)
    
    sb.icons = fonts.NewFontFromFile("/home/kernelcode/icon_font.fnt")
    sb.tiny = fonts.NewFontFromFile("/home/kernelcode/tiny_font.fnt")
    
    sb.State = STATE_NONE
    sb.Repeat, sb.Shuffle = false, false
    sb.TrackNum, sb.Tracks = 0, 0
    
    sb.stateLabel = NewLabel(sb.view.Widget, sb.icons)
    sb.stateLabel.SetPos(image.Point{2, 0})
    sb.repeatLabel = NewLabel(sb.view.Widget, sb.icons)
    sb.repeatLabel.SetPos(image.Point{10, 0})
    sb.shuffleLabel = NewLabel(sb.view.Widget, sb.icons)
    sb.shuffleLabel.SetPos(image.Point{19, 0})
    sb.tracksLabel = NewLabel(sb.view.Widget, sb.tiny)
    sb.tracksLabel.AutoWidth = false
    sb.tracksLabel.SetWidth(70)
    width := sb.view.Bounds().Dx()
    sb.tracksLabel.SetPos(image.Point{width - 70, 0})
    sb.tracksLabel.HAlign = Right
    
    sb.view.AddChild(sb.stateLabel)
    sb.view.AddChild(sb.repeatLabel)
    sb.view.AddChild(sb.shuffleLabel)
    sb.view.AddChild(sb.tracksLabel)
    
    return sb
}

func (sb *StatusBar) IsDirty() bool {
    //sb.Update()
    return sb.view.IsDirty()
}

func (sb *StatusBar) Update() {
    switch sb.State {
        case STATE_PLAYING:
            sb.stateLabel.Text = ICON_PLAY
        case STATE_PAUSED:
            sb.stateLabel.Text = ICON_PAUSE
        case STATE_NONE:
            fallthrough
        default:
            sb.stateLabel.Text = ICON_BLANK
    }
    
    if (sb.Repeat) {
        sb.repeatLabel.Text = ICON_REPEAT
    } else {
        sb.repeatLabel.Text = ICON_BLANK
    }
    
    if (sb.Shuffle) {
        sb.shuffleLabel.Text = ICON_SHUFFLE
    } else {
        sb.shuffleLabel.Text = ICON_BLANK
    }

    sb.tracksLabel.Text = fmt.Sprintf("%d/%d", sb.TrackNum, sb.Tracks)
    
    sb.view.Update()
}

func (sb *StatusBar) IsVisible() bool {
    return sb.view.IsVisible()
}

func (sb *StatusBar) Draw(to draw.Image) {
    
    if (sb.IsDirty()) {
        sb.Update()
    }

    if (sb.IsVisible()) {
        sb.view.Draw(to)
    }
}
