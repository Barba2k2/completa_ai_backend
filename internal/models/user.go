package models

import "time"

type User struct {
	ID          string     `json:"id"`
	Email       string     `json:"email"`
	Name        *string    `json:"name,omitempty"`
	DisplayName *string    `json:"displayName,omitempty"`
	PhotoUrl    *string    `json:"photoUrl,omitempty"`
	PhoneNumber *string    `json:"phoneNumber,omitempty"`
	FirebaseUID *string    `json:"firebaseUid,omitempty"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}
