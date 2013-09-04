package progbar

import (
    "image"
    "testing"
    "time"
    "github.com/usedbytes/s4548"
)

//import _ "image/png"

func TestHorz(t *testing.T) {
    screen := s4548.NewS4548("/dev/s4548-0")

    bar := NewProgressBar(nil)
    bar.SetWidth(101)
    bar.SetHeight(7)
    bar.SetPos(image.Point{0, 10})
    
    for i := 0; i <= 110; i++ {
        bar.Progress = i
        bar.MakeGraphics()
        bar.Draw(screen)
        screen.Scanout()
        time.Sleep(100 * time.Millisecond)
    }
    
    bar.Direction = Vertical
    bar.SetWidth(30)
    bar.SetHeight(40)
    bar.SetPos(image.Point{10, 0})
    for i := 0; i <= 110; i++ {
        bar.Progress = i
        bar.MakeGraphics()
        bar.Draw(screen)
        screen.Scanout()
        time.Sleep(100 * time.Millisecond)
    }    

    screen.Close()
    time.Sleep(1000 * time.Millisecond)
}
