package serializer

import (
	"dendro-age/internal/app/ds"
	"dendro-age/internal/app/repository"
)

type PatternDTO struct {
	ID               uint   `json:"id"`
	Title            string `json:"title"`
	ShortDescription string `json:"short_description"`
	FullDescription  string `json:"full_description"`
	Status           string `json:"status"`
	ImageURL         string `json:"image_url"`
	VideoURL         string `json:"video_url"`
	EventYear        *int   `json:"event_year"`
	RingCount        *int   `json:"ring_count"`
	LikeCount        int    `json:"like_count"`
	Liked            int    `json:"liked"`
	Mine             int    `json:"mine"`
	CreatorID        uint   `json:"creator_id"`
	CreatedAt        string `json:"created_at"`
	FormedAt         string `json:"formed_at,omitempty"`
}

type ResearcherDTO struct {
	ID       uint   `json:"id"`
	Login    string `json:"login"`
	FullName string `json:"full_name"`
}

type MessageDTO struct {
	Message string `json:"message"`
}

func Pattern(pattern ds.DendroScalePattern, likeCount int, liked int, mine int) PatternDTO {
	dto := PatternDTO{
		ID:               pattern.ID,
		Title:            pattern.Title,
		ShortDescription: pattern.ShortDescription,
		FullDescription:  pattern.FullDescription,
		Status:           pattern.Status,
		ImageURL:         repository.ResolveMedia(pattern.ImageURL, repository.DefaultImagePath),
		VideoURL:         repository.ResolveMedia(pattern.VideoURL, repository.DefaultVideoPath),
		EventYear:        pattern.EventYear,
		RingCount:        pattern.RingCount,
		LikeCount:        likeCount,
		Liked:            liked,
		Mine:             mine,
		CreatorID:        pattern.CreatorID,
		CreatedAt:        pattern.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if pattern.FormedAt != nil {
		dto.FormedAt = pattern.FormedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	return dto
}

func Patterns(patterns []ds.DendroScalePattern, likes map[uint]int, liked map[uint]int, mine map[uint]int) []PatternDTO {
	out := make([]PatternDTO, 0, len(patterns))
	for _, pattern := range patterns {
		out = append(out, Pattern(pattern, likes[pattern.ID], liked[pattern.ID], mine[pattern.ID]))
	}
	return out
}

func Researcher(r ds.DendroResearcher) ResearcherDTO {
	return ResearcherDTO{
		ID:       r.ID,
		Login:    r.Login,
		FullName: r.FullName,
	}
}
