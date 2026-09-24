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
	PasswordHash *string   `gorm:"type:varchar(255)" json:"-"` // Nullable для пользователей OAuth
	AvatarURL    string    `json:"avatar_url"`

	// Поля для OAuth
	AuthProvider string `gorm:"type:varchar(50);default:'local'" json:"auth_provider"` // local, google, github, yandex, vk
	ProviderID   string `gorm:"type:varchar(255);index" json:"provider_id,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Reviews   []Review  `gorm:"foreignKey:UserID" json:"reviews,omitempty"`
}

type Movie struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Slug            string         `gorm:"uniqueIndex;not null" json:"slug"`
	Title           string         `gorm:"index;not null" json:"title"`  // Название с Кинопоиска
	OriginalTitle   string         `json:"original_title"`               // Оригинальное с IMDb
	Description     string         `gorm:"type:text" json:"description"` // Описание с Кинопоиска
	ReleaseYear     int            `json:"release_year"`
	PosterURL       string         `json:"poster_url"`
	IMDbID          string         `gorm:"column:imdb_id;uniqueIndex;not null" json:"imdb_id"` // tt1375666
	ExternalRatings datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"external_ratings"`    // {"imdb": 8.8, "kp": 8.7, "letterboxd": 4.3}
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
