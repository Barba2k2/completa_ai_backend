package handlers

import (
	"fmt"
	"net/http"

	"github.com/completaai/backend/internal/middleware"
	"github.com/completaai/backend/internal/models"
	"github.com/completaai/backend/internal/repository"
)

type StatsHandler struct {
	repo *repository.CollectionRepository
}

func NewStatsHandler(repo *repository.CollectionRepository) *StatsHandler {
	return &StatsHandler{repo: repo}
}

func (h *StatsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	collection, err := h.repo.GetOrCreate(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get collection")
		return
	}

	stats := h.calculateStats(collection.Stickers)
	respondJSON(w, http.StatusOK, stats)
}

func (h *StatsHandler) calculateStats(stickers map[string]int) *models.CollectionStats {
	totalStickers := models.GetTotalStickers()
	ownedStickers := 0
	duplicateCount := 0
	sectionProgress := make(map[string]*models.SectionStats)

	for _, section := range models.AlbumSections {
		sectionProgress[section.Code] = &models.SectionStats{
			Name:    section.Name,
			Total:   section.Total,
			Owned:   0,
			Missing: section.Total,
			Rate:    0,
		}
	}

	for code, qty := range stickers {
		if qty > 0 {
			ownedStickers++
			sectionCode := h.extractSectionCode(code)
			if section, exists := sectionProgress[sectionCode]; exists {
				section.Owned++
				section.Missing--
			}
		}
		if qty > 1 {
			duplicateCount += qty - 1
		}
	}

	for _, section := range sectionProgress {
		if section.Total > 0 {
			section.Rate = float64(section.Owned) / float64(section.Total) * 100
		}
	}

	completionRate := float64(0)
	if totalStickers > 0 {
		completionRate = float64(ownedStickers) / float64(totalStickers) * 100
	}

	return &models.CollectionStats{
		TotalStickers:   totalStickers,
		OwnedStickers:   ownedStickers,
		MissingStickers: totalStickers - ownedStickers,
		DuplicateCount:  duplicateCount,
		CompletionRate:  completionRate,
		SectionProgress: sectionProgress,
	}
}

func (h *StatsHandler) extractSectionCode(code string) string {
	for i := len(code) - 1; i >= 0; i-- {
		if code[i] < '0' || code[i] > '9' {
			return code[:i+1]
		}
	}
	return code
}

func (h *StatsHandler) GetSections(w http.ResponseWriter, r *http.Request) {
	sections := make([]map[string]interface{}, len(models.AlbumSections))
	for i, section := range models.AlbumSections {
		stickers := make([]map[string]interface{}, section.Total)
		for j := 1; j <= section.Total; j++ {
			stickers[j-1] = map[string]interface{}{
				"code":     fmt.Sprintf("%s%02d", section.Code, j),
				"number":   j,
				"is_shiny": j == 1,
			}
		}
		sections[i] = map[string]interface{}{
			"code":     section.Code,
			"name":     section.Name,
			"total":    section.Total,
			"type":     section.Type,
			"stickers": stickers,
		}
	}

	respondJSON(w, http.StatusOK, sections)
}
