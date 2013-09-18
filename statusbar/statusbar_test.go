package statusbar

import (
    "fmt"
    "github.com/usedbytes/s4548"
    "testing"
    "time"
    "math/rand"
)

func TestStatusBar(t* testing.T) {
    fmt.Println(s4548.GetS4548EnvPath())
    screen := s4548.NewS4548(s4548.GetS4548EnvPath())
    
    statusbar := NewStatusBar(nil)
    state := false
    statusbar.Tracks = 1000
    r := rand.New(rand.NewSource(99))

    for {
        time.Sleep(1000 * time.Millisecond)
        if state {
            statusbar.State = STATE_PAUSED
        } else {
            statusbar.State = STATE_PLAYING
        }
        state = !state
        statusbar.Repeat = !statusbar.Repeat
        statusbar.Shuffle = !statusbar.Shuffle

        statusbar.TrackNum = r.Intn(1000)
        
        statusbar.Draw(screen)
        screen.Scanout()
    }
}
