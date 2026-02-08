package services

import (
	"sync"
	"time"
)

// EntryMetaData is the metadata of a mod
type EntryMetaData struct {
	ID          string   `json:"id"` // id / slug
	Title       string   `json:"title"`
	Categories  []string `json:"categories"`
	Description string   `json:"description"`
	Wiki        string   `json:"wiki,omitempty"`
	Source      string   `json:"source,omitempty"`
	UpdatedAt   time.Time
}

// MetaData is the metadata map of mod metadata
type MetaData struct {
	SchemaVersion int                      `json:"schemaVersion"`
	Entries       map[string]EntryMetaData `json:"entries"`
	sync.Mutex
}

func (md *MetaData) FilterMissing(ids []string) []string {
	var missing []string
	for _, id := range ids {
		if _, ok := md.Entries[id]; !ok {
			missing = append(missing, id)
		}
	}
	return missing
}

func (md *MetaData) FilterStale(threshold time.Duration) []string {
	cutoff := time.Now().Add(-threshold)
	var stale []string

	for id, meta := range md.Entries {
		if meta.UpdatedAt.Before(cutoff) {
			stale = append(stale, id)
		}
	}
	return stale
}

func (md *MetaData) Merge(nmd *MetaData) {
	if md.Entries == nil {
		md.Entries = make(map[string]EntryMetaData)
	}

	for id, incoming := range nmd.Entries {
		existing, exists := md.Entries[id]
		if !exists {
			md.Entries[id] = incoming
			continue
		}

		merged := existing

		if incoming.Title != "" {
			merged.Title = incoming.Title
		}
		if incoming.Description != "" {
			merged.Description = incoming.Description
		}
		if len(incoming.Categories) > 0 {
			merged.Categories = incoming.Categories
		}
		if incoming.Wiki != "" {
			merged.Wiki = incoming.Wiki
		}
		if incoming.Source != "" {
			merged.Source = incoming.Source
		}

		merged.UpdatedAt = incoming.UpdatedAt

		md.Entries[id] = merged
	}
}
