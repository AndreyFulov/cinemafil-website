package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Username     string    `gorm:"uniqueIndex;not null" json:"username"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Reviews      []Review  `gorm:"foreignKey:UserID" json:"reviews,omitempty"`
}

type Movie struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Slug            string         `gorm:"uniqueIndex;not null" json:"slug"`
	Title           string         `gorm:"index;not null" json:"title"`
	OriginalTitle   string         `json:"original_title"`
	ReleaseYear     int            `json:"release_year"`
	PosterURL       string         `json:"poster_url"`
	ExternalRatings datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"external_ratings"`
	AvgTotalScore   float32        `gorm:"default:0" json:"avg_total_score"`
	ReviewsCount    int            `gorm:"default:0" json:"reviews_count"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	Reviews         []Review       `gorm:"foreignKey:MovieID" json:"reviews,omitempty"`
}

type Review struct {
	ID      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID  uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_user_movie;not null" json:"user_id"`
	MovieID uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_user_movie;not null" json:"movie_id"`

	// Оценки по категориям (1-10)
	ScoreVisual  int8 `gorm:"not null" json:"score_visual"`
	ScoreAudio   int8 `gorm:"not null" json:"score_audio"`
	ScoreStory   int8 `gorm:"not null" json:"score_story"`
	ScoreActing  int8 `gorm:"not null" json:"score_acting"`
	ScoreEmotion int8 `gorm:"not null" json:"score_emotion"`

	TotalScore  float32        `gorm:"not null" json:"total_score"`
	ReviewText  string         `gorm:"type:text" json:"review_text"`
	ExtraScores datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"extra_scores"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// Автоматический пересчет итоговой оценки перед сохранением
func (r *Review) BeforeSave(tx *gorm.DB) error {
	sum := float32(r.ScoreVisual + r.ScoreAudio + r.ScoreStory + r.ScoreActing + r.ScoreEmotion)
	r.TotalScore = sum / 5.0
	return nil
}
