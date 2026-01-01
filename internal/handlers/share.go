package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/completaai/backend/internal/middleware"
	"github.com/completaai/backend/internal/models"
	"github.com/completaai/backend/internal/repository"
)

type ShareHandler struct {
	repo *repository.CollectionRepository
}

func NewShareHandler(repo *repository.CollectionRepository) *ShareHandler {
	return &ShareHandler{repo: repo}
}

func (h *ShareHandler) GetMissing(w http.ResponseWriter, r *http.Request) {
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

	missing := h.calculateMissing(collection.Stickers)
	text := h.formatMissingText(missing)

	respondJSON(w, http.StatusOK, models.ShareList{
		Missing: missing,
		Text:    text,
	})
}

func (h *ShareHandler) GetDuplicates(w http.ResponseWriter, r *http.Request) {
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

	duplicates := h.calculateDuplicates(collection.Stickers)
	text := h.formatDuplicatesText(duplicates)

	respondJSON(w, http.StatusOK, models.ShareList{
		Duplicates: duplicates,
		Text:       text,
	})
}

func (h *ShareHandler) GetBoth(w http.ResponseWriter, r *http.Request) {
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

	missing := h.calculateMissing(collection.Stickers)
	duplicates := h.calculateDuplicates(collection.Stickers)

	var textBuilder strings.Builder
	textBuilder.WriteString(h.formatMissingText(missing))
	textBuilder.WriteString("\n")
	textBuilder.WriteString(h.formatDuplicatesText(duplicates))

	respondJSON(w, http.StatusOK, models.ShareList{
		Missing:    missing,
		Duplicates: duplicates,
		Text:       textBuilder.String(),
	})
}

func (h *ShareHandler) calculateMissing(stickers map[string]int) *models.ShareSection {
	sections := make(map[string][]models.ShareItem)
	totalCount := 0

	for _, section := range models.AlbumSections {
		var missingItems []models.ShareItem
		for i := 1; i <= section.Total; i++ {
			code := fmt.Sprintf("%s%02d", section.Code, i)
			if qty, exists := stickers[code]; !exists || qty == 0 {
				missingItems = append(missingItems, models.ShareItem{Code: code})
				totalCount++
			}
		}
		if len(missingItems) > 0 {
			sections[section.Code] = missingItems
		}
	}

	return &models.ShareSection{
		Count:    totalCount,
		Sections: sections,
	}
}

func (h *ShareHandler) calculateDuplicates(stickers map[string]int) *models.ShareSection {
	sections := make(map[string][]models.ShareItem)
	totalCount := 0

	stickersBySection := make(map[string][]models.ShareItem)
	for code, qty := range stickers {
		if qty > 1 {
			sectionCode := h.extractSectionCode(code)
			stickersBySection[sectionCode] = append(stickersBySection[sectionCode], models.ShareItem{
				Code:     code,
				Quantity: qty - 1,
			})
			totalCount += qty - 1
		}
	}

	for sectionCode, items := range stickersBySection {
		sort.Slice(items, func(i, j int) bool {
			return items[i].Code < items[j].Code
		})
		sections[sectionCode] = items
	}

	return &models.ShareSection{
		Count:    totalCount,
		Sections: sections,
	}
}

func (h *ShareHandler) extractSectionCode(code string) string {
	for i := len(code) - 1; i >= 0; i-- {
		if code[i] < '0' || code[i] > '9' {
			return code[:i+1]
		}
	}
	return code
}

func (h *ShareHandler) formatMissingText(missing *models.ShareSection) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("FALTANTES (%d):\n", missing.Count))

	sectionCodes := make([]string, 0, len(missing.Sections))
	for code := range missing.Sections {
		sectionCodes = append(sectionCodes, code)
	}
	sort.Strings(sectionCodes)

	for _, sectionCode := range sectionCodes {
		items := missing.Sections[sectionCode]
		codes := make([]string, len(items))
		for i, item := range items {
			codes[i] = h.formatStickerNumber(item.Code)
		}
		builder.WriteString(fmt.Sprintf("%s %s\n", sectionCode, strings.Join(codes, ", ")))
	}

	return builder.String()
}

func (h *ShareHandler) formatDuplicatesText(duplicates *models.ShareSection) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("REPETIDAS (%d):\n", duplicates.Count))

	sectionCodes := make([]string, 0, len(duplicates.Sections))
	for code := range duplicates.Sections {
		sectionCodes = append(sectionCodes, code)
	}
	sort.Strings(sectionCodes)

	for _, sectionCode := range sectionCodes {
		items := duplicates.Sections[sectionCode]
		codes := make([]string, len(items))
		for i, item := range items {
			num := h.formatStickerNumber(item.Code)
			if item.Quantity > 1 {
				codes[i] = fmt.Sprintf("%s (x%d)", num, item.Quantity)
			} else {
				codes[i] = num
			}
		}
		builder.WriteString(fmt.Sprintf("%s %s\n", sectionCode, strings.Join(codes, ", ")))
	}

	return builder.String()
}

func (h *ShareHandler) formatStickerNumber(code string) string {
	for i := len(code) - 1; i >= 0; i-- {
		if code[i] < '0' || code[i] > '9' {
			return code[i+1:]
		}
	}
	return code
}
