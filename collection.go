package audioghost

import "time"

type Collection struct {
	Name        string
	Path        string
	Audiobooks  Audiobooks
	Collections Collections
	Playtime    time.Duration
}

type Collections map[string]*Collection

