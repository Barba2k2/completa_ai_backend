package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/completaai/backend/internal/middleware"
	"github.com/completaai/backend/internal/repository"
	"github.com/completaai/backend/internal/services"
)

type ProfileHandler struct {
	userRepo     *repository.UserRepository
	firebaseSvc  *services.FirebaseService
}

func NewProfileHandler(
	userRepo *repository.UserRepository,
	firebaseSvc *services.FirebaseService,
) *ProfileHandler {
	return &ProfileHandler{
		userRepo:    userRepo,
		firebaseSvc: firebaseSvc,
	}
}

type UpdateProfileRequest struct {
	DisplayName *string `json:"displayName,omitempty"`
	PhotoUrl    *string `json:"photoUrl,omitempty"`
	PhoneNumber *string `json:"phoneNumber,omitempty"`
}

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r.Context())

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r.Context())

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.UpdateProfile(userID, req.DisplayName, req.PhotoUrl, req.PhoneNumber)
	if err != nil {
		http.Error(w, `{"error":"failed to update profile"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *ProfileHandler) UpdatePhoto(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		respondError(w, http.StatusBadRequest, "invalid multipart form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		respondError(w, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	photoUrl, err := h.firebaseSvc.UploadProfileImage(r.Context(), userID, file, contentType)
	if err != nil {
		log.Printf("[Profile] Upload photo error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to upload image")
		return
	}

	user, err := h.userRepo.UpdateProfile(userID, nil, &photoUrl, nil)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update profile")
		return
	}

	respondJSON(w, http.StatusOK, user)
}
