package audioghost

type Audiobook struct {
	Name  string
	Path  string
	Files []string
}

type Audiobooks map[string]*Audiobook
