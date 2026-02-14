package meta

import (
	"time"

	"github.com/w1lam/Packages/modrinth"
)

// ResolveMetaData resolved metadata of given slugs
func ResolveMetaData(ids []string) (*MetaData, error) {
	mrProj, err := modrinth.BatchFetchModrinthProjects(ids)
	if err != nil {
		return nil, err
	}

	out := MetaData{
		SchemaVersion: 1,
		Entries:       make(map[string]EntryMetaData),
	}

	for _, p := range mrProj {
		out.Entries[p.Slug] = EntryMetaData{
			ID:          p.ID,
			Title:       p.Title,
			Categories:  p.Categories,
			Description: p.Description,
			Wiki:        p.Wiki,
			Source:      p.Source,
			UpdatedAt:   time.Now(),
		}
	}

	return &out, nil
}
