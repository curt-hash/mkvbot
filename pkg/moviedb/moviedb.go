package moviedb

type Metadata struct {
	Name string
	Year int
	Tag  string
}

type DB interface {
	FuzzySearchTitle(q string) ([]*Metadata, error)
}
