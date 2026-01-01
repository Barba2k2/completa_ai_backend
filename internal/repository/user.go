package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/completaai/backend/internal/models"
)

type UserRepository struct {
	db        *sql.DB
	jwtSecret string
}

func NewUserRepository(db *sql.DB, jwtSecret string) *UserRepository {
	return &UserRepository{db: db, jwtSecret: jwtSecret}
}

func (r *UserRepository) Create(email, password, name string) (*models.User, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user := &models.User{}
	err = r.db.QueryRow(`
		INSERT INTO users (email, password_hash, name)
		VALUES ($1, $2, $3)
		RETURNING id, email, name, created_at, updated_at
	`, email, string(hash), name).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, "", errors.New("email already exists")
	}

	token, err := r.generateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (r *UserRepository) Authenticate(email, password string) (*models.User, string, error) {
	user := &models.User{}
	var hash string

	err := r.db.QueryRow(`
		SELECT id, email, name, password_hash, created_at, updated_at
		FROM users WHERE email = $1
	`, email).Scan(&user.ID, &user.Email, &user.Name, &hash, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := r.generateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(`
		SELECT id, email, name, display_name, photo_url, phone_number, firebase_uid, last_login_at, created_at, updated_at
		FROM users WHERE id = $1
	`, id).Scan(&user.ID, &user.Email, &user.Name, &user.DisplayName, &user.PhotoUrl, &user.PhoneNumber, &user.FirebaseUID, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByFirebaseUID(firebaseUID string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(`
		SELECT id, email, name, display_name, photo_url, phone_number, firebase_uid, last_login_at, created_at, updated_at
		FROM users WHERE firebase_uid = $1
	`, firebaseUID).Scan(&user.ID, &user.Email, &user.Name, &user.DisplayName, &user.PhotoUrl, &user.PhoneNumber, &user.FirebaseUID, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) CreateFromFirebase(firebaseUID, email string, displayName, photoUrl *string) (*models.User, error) {
	user := &models.User{}
	now := time.Now()
	err := r.db.QueryRow(`
		INSERT INTO users (email, password_hash, firebase_uid, display_name, photo_url, last_login_at)
		VALUES ($1, '', $2, $3, $4, $5)
		RETURNING id, email, name, display_name, photo_url, phone_number, firebase_uid, last_login_at, created_at, updated_at
	`, email, firebaseUID, displayName, photoUrl, now).Scan(
		&user.ID, &user.Email, &user.Name, &user.DisplayName, &user.PhotoUrl,
		&user.PhoneNumber, &user.FirebaseUID, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) UpdateProfile(id string, displayName, photoUrl, phoneNumber *string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(`
		UPDATE users
		SET display_name = COALESCE($2, display_name),
		    photo_url = COALESCE($3, photo_url),
		    phone_number = COALESCE($4, phone_number),
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, email, name, display_name, photo_url, phone_number, firebase_uid, last_login_at, created_at, updated_at
	`, id, displayName, photoUrl, phoneNumber).Scan(
		&user.ID, &user.Email, &user.Name, &user.DisplayName, &user.PhotoUrl,
		&user.PhoneNumber, &user.FirebaseUID, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) UpdateLastLogin(id string) error {
	_, err := r.db.Exec(`
		UPDATE users SET last_login_at = NOW(), updated_at = NOW() WHERE id = $1
	`, id)
	return err
}

func (r *UserRepository) Delete(id string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM collections WHERE user_id = $1`, id)
	if err != nil {
		return err
	}

	result, err := tx.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return tx.Commit()
}

func (r *UserRepository) generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(r.jwtSecret))
}
