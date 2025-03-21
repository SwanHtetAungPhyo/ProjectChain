package services

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

var log = logrus.New()

func AnimationLoop(message string, duration time.Duration) {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	start := time.Now()

	for {
		for _, frame := range frames {
			fmt.Printf("\r%s %s", frame, message)
			time.Sleep(100 * time.Millisecond)

			if time.Since(start) > duration {
				fmt.Printf("\r✅ %s\n", message)
				return
			}
		}
	}
}
