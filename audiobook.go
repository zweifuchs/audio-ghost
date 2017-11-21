package audio_ghost

type Audiobook struct {
	Name  string
	Path  string
	Files []string
}

type Audiobooks map[string]*Audiobook
