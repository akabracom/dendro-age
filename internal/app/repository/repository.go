package repository

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"dendro-age/internal/app/currentresearcher"
	"dendro-age/internal/app/ds"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	DefaultImagePath = "/static/defaults/pattern_default.svg"
	DefaultVideoPath = "/static/defaults/pattern_default.mp4"
	minioBase        = "http://localhost:9000/drengo-age-img/"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&ds.DendroResearcher{}, &ds.DendroScalePattern{}, &ds.DendroPatternLike{}); err != nil {
		return nil, err
	}
	if err := db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_one_draft_per_researcher ON dendro_scale_patterns (creator_id) WHERE status = 'draft'`).Error; err != nil {
		return nil, err
	}
	repo := &Repository{db: db}
	if err := repo.seedIfEmpty(); err != nil {
		return nil, err
	}
	if err := repo.backfillFullDescriptions(); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *Repository) sqlDB() (*sql.DB, error) {
	return r.db.DB()
}

func (r *Repository) CurrentResearcher() (ds.DendroResearcher, error) {
	var researcher ds.DendroResearcher
	err := r.db.Where("login = ?", currentresearcher.Login()).First(&researcher).Error
	return researcher, err
}

func (r *Repository) ListPublished(yearFrom *int, yearTo *int) ([]ds.DendroScalePattern, error) {
	query := r.db.Where("status = ?", ds.StatusPublished)
	if yearFrom != nil {
		query = query.Where("event_year >= ?", *yearFrom)
	}
	if yearTo != nil {
		query = query.Where("event_year <= ?", *yearTo)
	}
	var patterns []ds.DendroScalePattern
	err := query.Order("id asc").Find(&patterns).Error
	return patterns, err
}

func (r *Repository) GetPublishedByID(id uint) (ds.DendroScalePattern, error) {
	var pattern ds.DendroScalePattern
	err := r.db.Where("id = ? AND status = ?", id, ds.StatusPublished).First(&pattern).Error
	if err != nil {
		return ds.DendroScalePattern{}, fmt.Errorf("паттерн дендрошкалы %d не найден", id)
	}
	return pattern, nil
}

func (r *Repository) GetFirstPublished() (ds.DendroScalePattern, error) {
	var pattern ds.DendroScalePattern
	err := r.db.Where("status = ?", ds.StatusPublished).Order("id asc").First(&pattern).Error
	return pattern, err
}

func (r *Repository) GetNextPublished(id uint) (ds.DendroScalePattern, error) {
	var pattern ds.DendroScalePattern
	err := r.db.Where("status = ? AND id > ?", ds.StatusPublished, id).Order("id asc").First(&pattern).Error
	if err == nil {
		return pattern, nil
	}
	return r.GetFirstPublished()
}

func (r *Repository) GetDraftByCreator(creatorID uint) (ds.DendroScalePattern, error) {
	var pattern ds.DendroScalePattern
	err := r.db.Where("creator_id = ? AND status = ?", creatorID, ds.StatusDraft).First(&pattern).Error
	return pattern, err
}

func (r *Repository) CreateDraft(creatorID uint, title string) (ds.DendroScalePattern, error) {
	var existing ds.DendroScalePattern
	err := r.db.Where("creator_id = ? AND status = ?", creatorID, ds.StatusDraft).First(&existing).Error
	if err == nil {
		return existing, nil
	}
	pattern := ds.DendroScalePattern{
		Title:     title,
		Status:    ds.StatusDraft,
		CreatorID: creatorID,
		CreatedAt: time.Now(),
	}
	if err := r.db.Create(&pattern).Error; err != nil {
		return ds.DendroScalePattern{}, err
	}
	return pattern, nil
}

func (r *Repository) PublishDraft(creatorID uint, shortDescription string, fullDescription string, eventYear int, ringCount int) (ds.DendroScalePattern, error) {
	pattern, err := r.GetDraftByCreator(creatorID)
	if err != nil {
		return ds.DendroScalePattern{}, fmt.Errorf("черновик паттерна не найден")
	}
	now := time.Now()
	pattern.ShortDescription = shortDescription
	pattern.FullDescription = fullDescription
	pattern.EventYear = eventYear
	pattern.RingCount = ringCount
	pattern.Status = ds.StatusPublished
	pattern.FormedAt = &now
	if err := r.db.Save(&pattern).Error; err != nil {
		return ds.DendroScalePattern{}, err
	}
	return pattern, nil
}

func (r *Repository) SoftDeletePattern(id uint) error {
	sqlDB, err := r.sqlDB()
	if err != nil {
		return err
	}
	tx, err := sqlDB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`UPDATE dendro_scale_patterns SET status = $1 WHERE id = $2`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(ds.StatusDeleted, id); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *Repository) CountLikes(patternID uint) int {
	var count int64
	r.db.Model(&ds.DendroPatternLike{}).Where("pattern_id = ?", patternID).Count(&count)
	return int(count)
}

func ResolveMedia(url string, fallback string) string {
	url = strings.TrimSpace(url)
	if url == "" {
		return fallback
	}
	if strings.HasPrefix(url, "/") {
		return url
	}
	client := &http.Client{Timeout: 800 * time.Millisecond}
	resp, err := client.Head(url)
	if err != nil || resp.StatusCode >= 400 {
		return fallback
	}
	return url
}

func (r *Repository) seedIfEmpty() error {
	var count int64
	if err := r.db.Model(&ds.DendroResearcher{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	shebe := ds.DendroResearcher{Login: currentresearcher.Login(), FullName: "Исследователь дендрошкалы"}
	anna := ds.DendroResearcher{Login: "anna", FullName: "Анна Кольцова"}
	if err := r.db.Create(&shebe).Error; err != nil {
		return err
	}
	if err := r.db.Create(&anna).Error; err != nil {
		return err
	}

	formed := time.Date(2026, 9, 2, 12, 0, 0, 0, time.Local)
	patterns := []ds.DendroScalePattern{
		{
			Title: "Серия засух 1975–1977", ShortDescription: "Три узких кольца засухи подряд.",
			FullDescription: "Паттерн включает три последовательных узких годичных кольца. Снижение ширины прироста указывает на неблагоприятные условия в период с 1975 по 1977 год.",
			Status: ds.StatusPublished, ImageURL: minioBase + "drought_series.jpg", VideoURL: minioBase + "drought_series.mp4",
			EventYear: 1975, RingCount: 3, CreatedAt: formed, FormedAt: &formed, CreatorID: shebe.ID,
		},
		{
			Title: "Жара в Европе 2003–2004", ShortDescription: "Аномальная жара в кольцах 2003 года.",
			FullDescription: "Кластер узких колец 2003–2004 годов фиксирует европейскую жару. Год начала и число колец помогают привязать находку к дендрошкале.",
			Status: ds.StatusPublished, ImageURL: minioBase + "fire_scar.jpg", VideoURL: minioBase + "drought_series_video.mp4",
			EventYear: 2003, RingCount: 2, CreatedAt: formed, FormedAt: &formed, CreatorID: anna.ID,
		},
		{
			Title: "Мегазасуха 1540–1542", ShortDescription: "Реперный кластер мегазасухи XVI века.",
			FullDescription: "Мегазасуха 1540–1542 годов даёт устойчивый маркер на мастер-хронологии. Паттерн используют при расчёте возраста находки.",
			Status: ds.StatusPublished, ImageURL: minioBase + "frost_damage.jpg", VideoURL: minioBase + "frost_damage.mp4",
			EventYear: 1540, RingCount: 3, CreatedAt: formed, FormedAt: &formed, CreatorID: anna.ID,
		},
		{
			Title: "Год без лета 1816–1817", ShortDescription: "Резкое падение прироста после 1816 года.",
			FullDescription: "1816 год без лета оставляет сдвиг прироста на два сезона. Год начала и количество колец читаются и на шкале, и на образце.",
			Status: ds.StatusPublished, ImageURL: minioBase + "growth_shift.jpg", VideoURL: minioBase + "growth_shift.mp4",
			EventYear: 1816, RingCount: 2, CreatedAt: formed, FormedAt: &formed, CreatorID: anna.ID,
		},
		{
			Title: "Черновик теплового стресса", ShortDescription: "",
			Status: ds.StatusDraft, ImageURL: "", VideoURL: "",
			EventYear: 0, RingCount: 0, CreatedAt: time.Now(), CreatorID: anna.ID,
		},
		{
			Title: "Повторный стресс 1784", ShortDescription: "Удалённый паттерн повторного стресса.",
			FullDescription: "Карточка снята с публикации. В интерфейсе не показывается.",
			Status: ds.StatusDeleted, ImageURL: minioBase + "repeated_stress.jpg", VideoURL: minioBase + "drought_series_video.mp4",
			EventYear: 1784, RingCount: 12, CreatedAt: formed, FormedAt: &formed, CreatorID: anna.ID,
		},
	}
	if err := r.db.Create(&patterns).Error; err != nil {
		return err
	}

	likes := []ds.DendroPatternLike{
		{ResearcherID: shebe.ID, PatternID: patterns[0].ID},
		{ResearcherID: anna.ID, PatternID: patterns[0].ID},
		{ResearcherID: shebe.ID, PatternID: patterns[1].ID},
		{ResearcherID: anna.ID, PatternID: patterns[2].ID},
		{ResearcherID: shebe.ID, PatternID: patterns[2].ID},
		{ResearcherID: anna.ID, PatternID: patterns[3].ID},
	}
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&likes).Error
}

func (r *Repository) backfillFullDescriptions() error {
	texts := map[string]string{
		"Серия засух 1975–1977":      "Паттерн включает три последовательных узких годичных кольца. Снижение ширины прироста указывает на неблагоприятные условия в период с 1975 по 1977 год.",
		"Жара в Европе 2003–2004":    "Кластер узких колец 2003–2004 годов фиксирует европейскую жару. Год начала и число колец помогают привязать находку к дендрошкале.",
		"Мегазасуха 1540–1542":       "Мегазасуха 1540–1542 годов даёт устойчивый маркер на мастер-хронологии. Паттерн используют при расчёте возраста находки.",
		"Год без лета 1816–1817":     "1816 год без лета оставляет сдвиг прироста на два сезона. Год начала и количество колец читаются и на шкале, и на образце.",
		"Повторный стресс 1784":      "Карточка снята с публикации. В интерфейсе не показывается.",
	}
	for title, text := range texts {
		if err := r.db.Model(&ds.DendroScalePattern{}).
			Where("title = ? AND (full_description = '' OR full_description IS NULL)", title).
			Update("full_description", text).Error; err != nil {
			return err
		}
	}
	return nil
}
