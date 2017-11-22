package audioghost

type Collection struct {
	Name        string
	Path        string
	Audiobooks  Audiobooks
	Collections Collections
}

type Collections map[string]*Collection
