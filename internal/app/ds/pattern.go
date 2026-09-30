package ds

import "time"

const (
	StatusDraft     = "draft"
	StatusPublished = "published"
	StatusDeleted   = "deleted"
)

type DendroScalePattern struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	Title            string     `gorm:"size:160;not null" json:"title"`
	ShortDescription string     `gorm:"size:500" json:"short_description"`
	FullDescription  string     `gorm:"size:2000" json:"full_description"`
	Status           string     `gorm:"size:20;not null;index" json:"status"`
	ImageURL         string     `gorm:"size:255" json:"image_url"`
	VideoURL         string     `gorm:"size:255" json:"video_url"`
	EventYear        int        `json:"event_year"`
	RingCount        int        `json:"ring_count"`
	CreatedAt        time.Time  `json:"created_at"`
	FormedAt         *time.Time `json:"formed_at"`
	CreatorID        uint       `gorm:"not null;index" json:"creator_id"`
	Creator          DendroResearcher `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
}

func (DendroScalePattern) TableName() string {
	return "dendro_scale_patterns"
}
