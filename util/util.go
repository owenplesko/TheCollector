package util

import (
	"strings"
	"time"
)

func ToRawName(displayName string) string {
	return strings.ToLower(strings.ReplaceAll(displayName, " ", ""))
}

func OnInterval(method func(), interval time.Duration) {
	go func() {
		for range time.Tick(interval) {
			go method()
		}
	}()
}
