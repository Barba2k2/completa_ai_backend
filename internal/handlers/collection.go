package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/completaai/backend/internal/middleware"
	"github.com/completaai/backend/internal/models"
	"github.com/completaai/backend/internal/repository"
)

type CollectionHandler struct {
	repo *repository.CollectionRepository
}

func NewCollectionHandler(repo *repository.CollectionRepository) *CollectionHandler {
	return &CollectionHandler{repo: repo}
}

func (h *CollectionHandler) Get(w http.ResponseWriter, r *http.Request) {
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

	respondJSON(w, http.StatusOK, collection)
}

func (h *CollectionHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req models.UpdateStickersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	collection, err := h.repo.UpdateStickers(userID, req.Stickers)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update collection")
		return
	}

	respondJSON(w, http.StatusOK, collection)
}

func (h *CollectionHandler) Sync(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req models.SyncRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	serverCollection, err := h.repo.GetOrCreate(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get collection")
		return
	}

	hasChanges := false
	mergedStickers := make(map[string]int)

	for code, qty := range serverCollection.Stickers {
		mergedStickers[code] = qty
	}

	if req.LastSyncedAt == nil || req.ClientTime.After(serverCollection.UpdatedAt) {
		for code, qty := range req.Stickers {
			if currentQty, exists := mergedStickers[code]; !exists || qty != currentQty {
				hasChanges = true
			}
			if qty > 0 {
				mergedStickers[code] = qty
			} else {
				delete(mergedStickers, code)
			}
		}
	} else {
		for code, qty := range req.Stickers {
			if _, exists := mergedStickers[code]; !exists {
				mergedStickers[code] = qty
				hasChanges = true
			}
		}
	}

	serverCollection.Stickers = mergedStickers
	if hasChanges {
		if err := h.repo.Update(serverCollection); err != nil {
			respondError(w, http.StatusInternalServerError, "failed to sync collection")
			return
		}
	}

	response := models.SyncResponse{
		Stickers:   mergedStickers,
		SyncedAt:   time.Now(),
		HasChanges: hasChanges,
	}

	respondJSON(w, http.StatusOK, response)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.APIResponse{
		Success: true,
		Data:    data,
	})
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.APIResponse{
		Success: false,
		Error:   message,
	})
}
