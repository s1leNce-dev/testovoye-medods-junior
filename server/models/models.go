package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Login    string    `json:"login" gorm:"unique"`
	Password string    `json:"password"`
}

type RefreshTokenSessions struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SessionID uuid.UUID `gorm:"type:uuid"`
	UserID    uuid.UUID `gorm:"type:uuid"`
	UserAgent string
	ClientIP  string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time

	User User `gorm:"forgeinKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
