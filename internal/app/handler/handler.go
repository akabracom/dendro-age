package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"dendro-age/internal/app/ds"
	"dendro-age/internal/app/repository"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{Repository: r}
}

type patternView struct {
	ds.DendroScalePattern
	LikeCount int
	ImageSrc  string
	VideoSrc  string
}

func (h *Handler) toView(pattern ds.DendroScalePattern) patternView {
	return patternView{
		DendroScalePattern: pattern,
		LikeCount:          h.Repository.CountLikes(pattern.ID),
		ImageSrc:           repository.ResolveMedia(pattern.ImageURL, repository.DefaultImagePath),
		VideoSrc:           repository.ResolveMedia(pattern.VideoURL, repository.DefaultVideoPath),
	}
}

func (h *Handler) toViews(patterns []ds.DendroScalePattern) []patternView {
	views := make([]patternView, 0, len(patterns))
	for _, pattern := range patterns {
		views = append(views, h.toView(pattern))
	}
	return views
}

func parseYearFromDate(raw string) *int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parsed, err := time.Parse("2006-01-02", raw)
	if err != nil {
		log.Println("некорректная дата фильтра:", err)
		return nil
	}
	year := parsed.Year()
	return &year
}

func (h *Handler) currentCreatorID() (uint, error) {
	researcher, err := h.Repository.CurrentResearcher()
	if err != nil {
		return 0, err
	}
	return researcher.ID, nil
}

func (h *Handler) renderNotFound(ctx *gin.Context) {
	ctx.HTML(http.StatusNotFound, "not_found.html", gin.H{
		"ActiveTab":   "feed",
		"RequestPath": ctx.Request.URL.Path,
	})
}

func (h *Handler) GetDendroScalePatterns(ctx *gin.Context) {
	dateFromRaw := ctx.Query("date_from")
	dateToRaw := ctx.Query("date_to")
	patterns, err := h.Repository.ListPublished(parseYearFromDate(dateFromRaw), parseYearFromDate(dateToRaw))
	if err != nil {
		log.Println(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.HTML(http.StatusOK, "mosaic.html", gin.H{
		"ActiveTab":           "mosaic",
		"DendroScalePatterns": h.toViews(patterns),
		"DateFrom":            dateFromRaw,
		"DateTo":              dateToRaw,
	})
}

func (h *Handler) GetDendroScalePatternsFeed(ctx *gin.Context) {
	idStr := ctx.Param("id")
	next := ctx.Query("next") == "true"

	var (
		pattern ds.DendroScalePattern
		err     error
	)

	if idStr == "" {
		pattern, err = h.Repository.GetFirstPublished()
	} else {
		id64, convErr := strconv.ParseUint(idStr, 10, 64)
		if convErr != nil {
			h.renderNotFound(ctx)
			return
		}
		id := uint(id64)
		if next {
			pattern, err = h.Repository.GetNextPublished(id)
		} else {
			pattern, err = h.Repository.GetPublishedByID(id)
		}
	}
	if err != nil {
		log.Println(err)
		h.renderNotFound(ctx)
		return
	}

	ctx.HTML(http.StatusOK, "feed.html", gin.H{
		"ActiveTab":          "feed",
		"DendroScalePattern": h.toView(pattern),
	})
}

func (h *Handler) GetDendroScalePatternsDraft(ctx *gin.Context) {
	creatorID, err := h.currentCreatorID()
	if err != nil {
		log.Println(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	pattern, err := h.Repository.GetDraftByCreator(creatorID)
	hasDraft := err == nil
	view := patternView{
		ImageSrc: repository.DefaultImagePath,
		VideoSrc: repository.DefaultVideoPath,
	}
	if hasDraft {
		view = h.toView(pattern)
	}
	ctx.HTML(http.StatusOK, "draft.html", gin.H{
		"ActiveTab":          "draft",
		"HasDraft":           hasDraft,
		"DendroScalePattern": view,
	})
}

func (h *Handler) PostDendroScalePatternNext(ctx *gin.Context) {
	creatorID, err := h.currentCreatorID()
	if err != nil {
		log.Println(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	title := strings.TrimSpace(ctx.PostForm("title"))
	if title == "" {
		ctx.Redirect(http.StatusSeeOther, "/dendro_scale_patterns/draft")
		return
	}
	if _, err := h.Repository.CreateDraft(creatorID, title); err != nil {
		log.Println(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Redirect(http.StatusSeeOther, "/dendro_scale_patterns/draft")
}

func (h *Handler) PostDendroScalePatternPublish(ctx *gin.Context) {
	creatorID, err := h.currentCreatorID()
	if err != nil {
		log.Println(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	shortDescription := strings.TrimSpace(ctx.PostForm("short_description"))
	fullDescription := strings.TrimSpace(ctx.PostForm("full_description"))
	eventYear, _ := strconv.Atoi(ctx.PostForm("event_year"))
	ringCount, _ := strconv.Atoi(ctx.PostForm("ring_count"))
	if _, err := h.Repository.PublishDraft(creatorID, shortDescription, fullDescription, eventYear, ringCount); err != nil {
		log.Println(err)
		ctx.Redirect(http.StatusSeeOther, "/dendro_scale_patterns/draft")
		return
	}
	ctx.Redirect(http.StatusSeeOther, "/dendro_scale_patterns")
}

func (h *Handler) PostDendroScalePatternDelete(ctx *gin.Context) {
	id64, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	if err := h.Repository.SoftDeletePattern(uint(id64)); err != nil {
		log.Println(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Redirect(http.StatusSeeOther, "/dendro_scale_patterns")
}
