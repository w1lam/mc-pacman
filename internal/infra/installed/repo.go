// Package installed holds installed repo implementation
package installed

type repo struct {
	path string
}

func New(path string) *repo {
	return &repo{
		path: path,
	}
}
