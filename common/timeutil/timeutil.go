package timeutil

import (
	"math"
	"math/rand"
	"time"
)

func Sleep(duration time.Duration) {
	seconds := math.Round(rand.Float64() * duration.Seconds())
	time.Sleep(duration + time.Second*time.Duration(seconds))
}

func HugeSleep() {
	Sleep(time.Minute * 5)
}
