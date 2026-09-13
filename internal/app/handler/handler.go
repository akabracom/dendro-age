package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

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
	repository.DendroScalePattern
	LikeCount int
}

func toView(pattern repository.DendroScalePattern) patternView {
	return patternView{
		DendroScalePattern: pattern,
		LikeCount:          len(pattern.LikeUserIDs),
	}
}

func toViews(patterns []repository.DendroScalePattern) []patternView {
	views := make([]patternView, 0, len(patterns))
	for _, pattern := range patterns {
		views = append(views, toView(pattern))
	}
	return views
}

func parseYear(raw string) *int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	year, err := strconv.Atoi(raw)
	if err != nil {
		log.Println("некорректный год фильтра:", err)
		return nil
	}
	return &year
}

func (h *Handler) GetDendroScalePatterns(ctx *gin.Context) {
	yearFromRaw := ctx.Query("year_from")
	yearToRaw := ctx.Query("year_to")
	yearFrom := parseYear(yearFromRaw)
	yearTo := parseYear(yearToRaw)

	patterns := h.Repository.ListPublished(yearFrom, yearTo)
	ctx.HTML(http.StatusOK, "mosaic.html", gin.H{
		"ActiveTab":           "mosaic",
		"DendroScalePatterns": toViews(patterns),
		"YearFrom":            yearFromRaw,
		"YearTo":              yearToRaw,
	})
}

func (h *Handler) GetDendroScalePatternsFeed(ctx *gin.Context) {
	idStr := ctx.Param("id")
	next := ctx.Query("next") == "true"

	var (
		pattern repository.DendroScalePattern
		err     error
	)

	if idStr == "" {
		pattern, err = h.Repository.GetFirstPublished()
	} else {
		id, convErr := strconv.Atoi(idStr)
		if convErr != nil {
			log.Println("некорректный id паттерна:", convErr)
			ctx.Status(http.StatusBadRequest)
			return
		}
		if next {
			pattern, err = h.Repository.GetNextPublished(id)
		} else {
			pattern, err = h.Repository.GetPublishedByID(id)
		}
	}

	if err != nil {
		log.Println(err)
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.HTML(http.StatusOK, "feed.html", gin.H{
		"ActiveTab":         "feed",
		"DendroScalePattern": toView(pattern),
	})
}

func (h *Handler) GetDendroScalePatternsDraft(ctx *gin.Context) {
	pattern, err := h.Repository.GetDraft()
	if err != nil {
		log.Println(err)
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.HTML(http.StatusOK, "draft.html", gin.H{
		"ActiveTab":          "draft",
		"DendroScalePattern": toView(pattern),
	})
}
