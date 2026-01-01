package repository

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/completaai/backend/internal/models"
)

type CollectionRepository struct {
	db *sql.DB
}

func NewCollectionRepository(db *sql.DB) *CollectionRepository {
	return &CollectionRepository{db: db}
}

func (r *CollectionRepository) GetByUserID(userID string) (*models.Collection, error) {
	var collection models.Collection
	var stickersJSON []byte

	err := r.db.QueryRow(`
		SELECT id, user_id, stickers, updated_at
		FROM collections
		WHERE user_id = $1
	`, userID).Scan(&collection.ID, &collection.UserID, &stickersJSON, &collection.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(stickersJSON, &collection.Stickers); err != nil {
		return nil, err
	}

	return &collection, nil
}

func (r *CollectionRepository) Create(userID string) (*models.Collection, error) {
	collection := &models.Collection{
		UserID:    userID,
		Stickers:  make(map[string]int),
		UpdatedAt: time.Now(),
	}

	stickersJSON, err := json.Marshal(collection.Stickers)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(`
		INSERT INTO collections (user_id, stickers, updated_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`, collection.UserID, stickersJSON, collection.UpdatedAt).Scan(&collection.ID)

	if err != nil {
		return nil, err
	}

	return collection, nil
}

func (r *CollectionRepository) Update(collection *models.Collection) error {
	stickersJSON, err := json.Marshal(collection.Stickers)
	if err != nil {
		return err
	}

	collection.UpdatedAt = time.Now()

	_, err = r.db.Exec(`
		UPDATE collections
		SET stickers = $1, updated_at = $2
		WHERE id = $3 AND user_id = $4
	`, stickersJSON, collection.UpdatedAt, collection.ID, collection.UserID)

	return err
}

func (r *CollectionRepository) UpdateStickers(userID string, stickers map[string]int) (*models.Collection, error) {
	collection, err := r.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	if collection == nil {
		collection, err = r.Create(userID)
		if err != nil {
			return nil, err
		}
	}

	for code, quantity := range stickers {
		if quantity <= 0 {
			delete(collection.Stickers, code)
		} else {
			collection.Stickers[code] = quantity
		}
	}

	if err := r.Update(collection); err != nil {
		return nil, err
	}

	return collection, nil
}

func (r *CollectionRepository) GetOrCreate(userID string) (*models.Collection, error) {
	collection, err := r.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	if collection == nil {
		return r.Create(userID)
	}

	return collection, nil
}
