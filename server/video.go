package main

import "time"

type VideoState struct {
	Playing            bool
	Progress           float64
	LastProgressUpdate time.Time
}
