package ds

type DendroResearcher struct {
	ID       uint   `gorm:"primaryKey"`
	Login    string `gorm:"size:64;uniqueIndex;not null"`
	FullName string `gorm:"size:128;not null"`
}

func (DendroResearcher) TableName() string {
	return "dendro_researchers"
}
