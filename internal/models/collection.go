package models

import (
	"time"

	"github.com/google/uuid"
)

type Collection struct {
	ID        uuid.UUID         `json:"id"`
	UserID    uuid.UUID         `json:"user_id"`
	Stickers  map[string]int    `json:"stickers"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type UpdateStickersRequest struct {
	Stickers map[string]int `json:"stickers"`
}

type SyncRequest struct {
	Stickers     map[string]int `json:"stickers"`
	ClientTime   time.Time      `json:"client_time"`
	LastSyncedAt *time.Time     `json:"last_synced_at,omitempty"`
}

type SyncResponse struct {
	Stickers   map[string]int `json:"stickers"`
	SyncedAt   time.Time      `json:"synced_at"`
	HasChanges bool           `json:"has_changes"`
}

type CollectionStats struct {
	TotalStickers    int                    `json:"total_stickers"`
	OwnedStickers    int                    `json:"owned_stickers"`
	MissingStickers  int                    `json:"missing_stickers"`
	DuplicateCount   int                    `json:"duplicate_count"`
	CompletionRate   float64                `json:"completion_rate"`
	SectionProgress  map[string]*SectionStats `json:"section_progress"`
}

type SectionStats struct {
	Name    string  `json:"name"`
	Total   int     `json:"total"`
	Owned   int     `json:"owned"`
	Missing int     `json:"missing"`
	Rate    float64 `json:"rate"`
}
