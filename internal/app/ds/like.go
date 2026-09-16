package ds

type DendroPatternLike struct {
	ID           uint `gorm:"primaryKey"`
	ResearcherID uint `gorm:"not null;uniqueIndex:idx_like_pair"`
	PatternID    uint `gorm:"not null;uniqueIndex:idx_like_pair"`
	Researcher   DendroResearcher   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Pattern      DendroScalePattern `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

func (DendroPatternLike) TableName() string {
	return "dendro_pattern_likes"
}
