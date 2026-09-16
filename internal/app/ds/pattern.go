package ds

import "time"

const (
	StatusDraft     = "draft"
	StatusPublished = "published"
	StatusDeleted   = "deleted"
)

type DendroScalePattern struct {
	ID               uint      `gorm:"primaryKey"`
	Title            string    `gorm:"size:160;not null"`
	ShortDescription string    `gorm:"size:500"`
	FullDescription  string    `gorm:"size:2000"`
	Status           string    `gorm:"size:20;not null;index"`
	ImageURL         string    `gorm:"size:255"`
	VideoURL         string    `gorm:"size:255"`
	EventYear        int
	RingCount        int
	CreatedAt        time.Time
	FormedAt         *time.Time
	CreatorID        uint `gorm:"not null;index"`
	Creator          DendroResearcher `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

func (DendroScalePattern) TableName() string {
	return "dendro_scale_patterns"
}
