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
    bar.SetPos(image.Point{0, 32})
    bar.MakeGraphics()
    bar.Draw(screen)
    screen.Scanout()

    screen.Close()
    time.Sleep(1000 * time.Millisecond)
}