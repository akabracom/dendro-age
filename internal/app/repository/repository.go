package repository

import "fmt"

const (
	StatusDraft     = "draft"
	StatusPublished = "published"
	StatusDeleted   = "deleted"

	MinioBaseURL = "http://localhost:9000/drengo-age-img/"
	SharedVideo  = "drought_series_video.mp4"
)

type DendroScalePattern struct {
	ID                int
	Title             string
	ShortDescription  string
	FullDescription   string
	Status            string
	ImageKey          string
	VideoKey          string
	EventYear         int
	EventEndYear      int
	RingCount         int
	LikeUserIDs       []int
}

func (p DendroScalePattern) ImageURL() string {
	return MinioBaseURL + p.ImageKey
}

func (p DendroScalePattern) VideoURL() string {
	return MinioBaseURL + p.VideoKey
}

func likeUsers(count int) []int {
	ids := make([]int, count)
	for i := range ids {
		ids[i] = i + 1
	}
	return ids
}

type Repository struct {
	dendro_scale_patterns []DendroScalePattern
}

func NewRepository() (*Repository, error) {
	return &Repository{
		dendro_scale_patterns: []DendroScalePattern{
			{
				ID:               1,
				Title:            "Серия засух 1975–1977",
				ShortDescription: "Три узких кольца засухи подряд.",
				FullDescription:  "Паттерн включает три последовательных узких годичных кольца. Снижение ширины прироста указывает на неблагоприятные условия в период с 1975 по 1977 год.",
				Status:           StatusPublished,
				ImageKey:         "drought_series.jpg",
				VideoKey:         "drought_series.mp4",
				EventYear:        1975,
				EventEndYear:     1977,
				RingCount:        3,
				LikeUserIDs:      likeUsers(128),
			},
			{
				ID:               2,
				Title:            "Жара в Европе 2003–2004",
				ShortDescription: "Аномальная жара в кольцах 2003 года.",
				FullDescription:  "Кластер узких колец 2003–2004 годов фиксирует европейскую жару. Год начала и число колец помогают привязать находку к дендрошкале.",
				Status:           StatusPublished,
				ImageKey:         "fire_scar.jpg",
				VideoKey:         SharedVideo,
				EventYear:        2003,
				EventEndYear:     2004,
				RingCount:        2,
				LikeUserIDs:      likeUsers(96),
			},
			{
				ID:               3,
				Title:            "Мегазасуха 1540–1542",
				ShortDescription: "Реперный кластер мегазасухи XVI века.",
				FullDescription:  "Мегазасуха 1540–1542 годов даёт устойчивый маркер на мастер-хронологии. Паттерн используют при расчёте возраста находки.",
				Status:           StatusPublished,
				ImageKey:         "frost_damage.jpg",
				VideoKey:         "frost_damage.mp4",
				EventYear:        1540,
				EventEndYear:     1542,
				RingCount:        3,
				LikeUserIDs:      likeUsers(214),
			},
			{
				ID:               4,
				Title:            "Год без лета 1816–1817",
				ShortDescription: "Резкое падение прироста после 1816 года.",
				FullDescription:  "1816 год без лета оставляет сдвиг прироста на два сезона. Год начала и количество колец читаются и на шкале, и на образце.",
				Status:           StatusPublished,
				ImageKey:         "growth_shift.jpg",
				VideoKey:         "growth_shift.mp4",
				EventYear:        1816,
				EventEndYear:     1817,
				RingCount:        2,
				LikeUserIDs:      likeUsers(173),
			},
			{
				ID:               5,
				Title:            "Серия засух 1975–1977",
				ShortDescription: "Черновик: три узких кольца засухи.",
				FullDescription:  "Паттерн включает три последовательных узких годичных кольца. Снижение ширины прироста указывает на неблагоприятные условия в период с 1975 по 1977 год.",
				Status:           StatusDraft,
				ImageKey:         "drought_series.jpg",
				VideoKey:         "drought_series.mp4",
				EventYear:        1975,
				EventEndYear:     1977,
				RingCount:        3,
				LikeUserIDs:      []int{},
			},
			{
				ID:               6,
				Title:            "Повторный стресс 1784",
				ShortDescription: "Удалённый паттерн повторного стресса.",
				FullDescription:  "Карточка снята с публикации. В интерфейсе не показывается.",
				Status:           StatusDeleted,
				ImageKey:         "repeated_stress.jpg",
				VideoKey:         SharedVideo,
				EventYear:        1784,
				EventEndYear:     1784,
				RingCount:        12,
				LikeUserIDs:      []int{9},
			},
		},
	}, nil
}

func (r *Repository) published() []DendroScalePattern {
	out := make([]DendroScalePattern, 0)
	for _, pattern := range r.dendro_scale_patterns {
		if pattern.Status == StatusPublished {
			out = append(out, pattern)
		}
	}
	return out
}

func (r *Repository) ListPublished(yearFrom *int, yearTo *int) []DendroScalePattern {
	patterns := r.published()
	if yearFrom == nil && yearTo == nil {
		return patterns
	}
	filtered := make([]DendroScalePattern, 0)
	for _, pattern := range patterns {
		endYear := pattern.EventEndYear
		if endYear == 0 {
			endYear = pattern.EventYear
		}
		if yearFrom != nil && endYear < *yearFrom {
			continue
		}
		if yearTo != nil && pattern.EventYear > *yearTo {
			continue
		}
		filtered = append(filtered, pattern)
	}
	return filtered
}

func (r *Repository) GetPublishedByID(id int) (DendroScalePattern, error) {
	for _, pattern := range r.published() {
		if pattern.ID == id {
			return pattern, nil
		}
	}
	return DendroScalePattern{}, fmt.Errorf("паттерн дендрошкалы %d не найден", id)
}

func (r *Repository) GetFirstPublished() (DendroScalePattern, error) {
	patterns := r.published()
	if len(patterns) == 0 {
		return DendroScalePattern{}, fmt.Errorf("нет опубликованных паттернов")
	}
	return patterns[0], nil
}

func (r *Repository) GetNextPublished(id int) (DendroScalePattern, error) {
	patterns := r.published()
	if len(patterns) == 0 {
		return DendroScalePattern{}, fmt.Errorf("нет опубликованных паттернов")
	}
	for i, pattern := range patterns {
		if pattern.ID == id {
			if i+1 < len(patterns) {
				return patterns[i+1], nil
			}
			return patterns[0], nil
		}
	}
	return patterns[0], nil
}

func (r *Repository) GetDraft() (DendroScalePattern, error) {
	for _, pattern := range r.dendro_scale_patterns {
		if pattern.Status == StatusDraft {
			return pattern, nil
		}
	}
	return DendroScalePattern{}, fmt.Errorf("черновик паттерна не найден")
}
