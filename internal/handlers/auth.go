package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/completaai/backend/internal/middleware"
	"github.com/completaai/backend/internal/repository"
	"github.com/completaai/backend/internal/services"
)

type AuthHandler struct {
	userRepo    *repository.UserRepository
	firebaseSvc *services.FirebaseService
}

func NewAuthHandler(userRepo *repository.UserRepository, firebaseSvc *services.FirebaseService) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, firebaseSvc: firebaseSvc}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

type FirebaseAuthRequest struct {
	IDToken string `json:"id_token"`
}

func jsonError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		jsonError(w, "email and password required", http.StatusBadRequest)
		return
	}

	user, token, err := h.userRepo.Create(req.Email, req.Password, req.Name)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{Token: token, User: user})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		jsonError(w, "email and password required", http.StatusBadRequest)
		return
	}

	user, token, err := h.userRepo.Authenticate(req.Email, req.Password)
	if err != nil {
		jsonError(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{Token: token, User: user})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r.Context())

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		jsonError(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r.Context())

	err := h.userRepo.Delete(userID)
	if err != nil {
		jsonError(w, "failed to delete account", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "account deleted successfully"})
}

func (h *AuthHandler) FirebaseAuth(w http.ResponseWriter, r *http.Request) {
	var req FirebaseAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[FirebaseAuth] Failed to decode request: %v", err)
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.IDToken == "" {
		log.Printf("[FirebaseAuth] Empty id_token")
		jsonError(w, "id_token required", http.StatusBadRequest)
		return
	}

	log.Printf("[FirebaseAuth] Verifying token (length: %d)", len(req.IDToken))
	firebaseToken, err := h.firebaseSvc.VerifyIDToken(r.Context(), req.IDToken)
	if err != nil {
		log.Printf("[FirebaseAuth] Token verification failed: %v", err)
		jsonError(w, "invalid firebase token", http.StatusUnauthorized)
		return
	}

	user, err := h.userRepo.GetByFirebaseUID(firebaseToken.UID)
	if err != nil {
		email, _ := firebaseToken.Claims["email"].(string)
		displayName, _ := firebaseToken.Claims["name"].(string)
		photoUrl, _ := firebaseToken.Claims["picture"].(string)

		var displayNamePtr, photoUrlPtr *string
		if displayName != "" {
			displayNamePtr = &displayName
		}
		if photoUrl != "" {
			photoUrlPtr = &photoUrl
		}

		user, err = h.userRepo.CreateFromFirebase(firebaseToken.UID, email, displayNamePtr, photoUrlPtr)
		if err != nil {
			jsonError(w, "failed to create user", http.StatusInternalServerError)
			return
		}
	} else {
		h.userRepo.UpdateLastLogin(user.ID)
	}

	token, err := h.userRepo.GenerateToken(user.ID)
	if err != nil {
		jsonError(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{Token: token, User: user})
}
