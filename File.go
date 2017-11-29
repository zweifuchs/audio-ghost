package audioghost

import "time"

type MediaFile struct {
	Name string
	Path string
	AudioBook *Audiobook
	Track int
	Playtime time.Duration
}