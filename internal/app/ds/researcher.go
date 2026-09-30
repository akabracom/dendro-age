package ds

type DendroResearcher struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Login    string `gorm:"size:64;uniqueIndex;not null" json:"login"`
	Password string `gorm:"size:128" json:"-"`
	FullName string `gorm:"size:128;not null" json:"full_name"`
}

func (DendroResearcher) TableName() string {
	return "dendro_researchers"
}
