package apihandler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"dendro-age/internal/app/auth"
	"dendro-age/internal/app/ds"
	"dendro-age/internal/app/miniofile"
	"dendro-age/internal/app/repository"
	"dendro-age/internal/app/serializer"

	"github.com/gin-gonic/gin"
)

type API struct {
	Repo   *repository.Repository
	Minio  *miniofile.Client
	Tokens *auth.Tokens
	Gate   *auth.Gate
}

func New(repo *repository.Repository, minioClient *miniofile.Client, tokens *auth.Tokens, gate *auth.Gate) *API {
	return &API{Repo: repo, Minio: minioClient, Tokens: tokens, Gate: gate}
}

func (a *API) currentCreatorID(ctx *gin.Context) (uint, error) {
	id, ok := auth.ResearcherID(ctx)
	if !ok {
		return 0, fmt.Errorf("исследователь не авторизован")
	}
	return id, nil
}

func parseYearFromDate(raw string) *int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parsed, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return nil
	}
	year := parsed.Year()
	return &year
}

func (a *API) likeMap(patterns []ds.DendroScalePattern) map[uint]int {
	likes := make(map[uint]int, len(patterns))
	for _, pattern := range patterns {
		likes[pattern.ID] = a.Repo.CountLikes(pattern.ID)
	}
	return likes
}

func (a *API) likedFlag(ctx *gin.Context, patternID uint) int {
	creatorID, err := a.currentCreatorID(ctx)
	if err != nil {
		return 0
	}
	if a.Repo.HasLike(creatorID, patternID) {
		return 1
	}
	return 0
}

func (a *API) likedMap(ctx *gin.Context, patterns []ds.DendroScalePattern) map[uint]int {
	flags := make(map[uint]int, len(patterns))
	for _, pattern := range patterns {
		flags[pattern.ID] = a.likedFlag(ctx, pattern.ID)
	}
	return flags
}

func (a *API) mineFlag(ctx *gin.Context, creatorID uint) int {
	id, err := a.currentCreatorID(ctx)
	if err != nil || id != creatorID {
		return 0
	}
	return 1
}

func (a *API) mineMap(ctx *gin.Context, patterns []ds.DendroScalePattern) map[uint]int {
	flags := make(map[uint]int, len(patterns))
	for _, pattern := range patterns {
		flags[pattern.ID] = a.mineFlag(ctx, pattern.CreatorID)
	}
	return flags
}

func (a *API) patternDTO(ctx *gin.Context, pattern ds.DendroScalePattern) serializer.PatternDTO {
	return serializer.Pattern(pattern, a.Repo.CountLikes(pattern.ID), a.likedFlag(ctx, pattern.ID), a.mineFlag(ctx, pattern.CreatorID))
}

func (a *API) ListPatterns(ctx *gin.Context) {
	patterns, err := a.Repo.ListPublished(parseYearFromDate(ctx.Query("date_from")), parseYearFromDate(ctx.Query("date_to")))
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, serializer.MessageDTO{Message: "ошибка получения списка"})
		return
	}
	ctx.JSON(http.StatusOK, serializer.Patterns(patterns, a.likeMap(patterns), a.likedMap(ctx, patterns), a.mineMap(ctx, patterns)))
}

func (a *API) GetFeed(ctx *gin.Context) {
	idStr := ctx.Param("id")
	next := ctx.Query("next") == "true"

	var (
		pattern ds.DendroScalePattern
		err     error
	)
	if idStr == "" {
		pattern, err = a.Repo.GetFirstPublished()
	} else {
		id64, convErr := strconv.ParseUint(idStr, 10, 64)
		if convErr != nil {
			ctx.JSON(http.StatusBadRequest, serializer.MessageDTO{Message: "некорректный id"})
			return
		}
		id := uint(id64)
		if next {
			pattern, err = a.Repo.GetNextPublished(id)
		} else {
			pattern, err = a.Repo.GetPublishedByID(id)
		}
	}
	if err != nil {
		ctx.JSON(http.StatusNotFound, serializer.MessageDTO{Message: "паттерн не найден или удалён"})
		return
	}
	ctx.JSON(http.StatusOK, a.patternDTO(ctx, pattern))
}

func (a *API) GetDraft(ctx *gin.Context) {
	creatorID, err := a.currentCreatorID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, serializer.MessageDTO{Message: "создатель не найден"})
		return
	}
	pattern, err := a.Repo.GetDraftByCreator(creatorID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, serializer.MessageDTO{Message: "черновик не найден"})
		return
	}
	ctx.JSON(http.StatusOK, a.patternDTO(ctx, pattern))
}

func (a *API) CreatePattern(ctx *gin.Context) {
	creatorID, err := a.currentCreatorID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, serializer.MessageDTO{Message: "создатель не найден"})
		return
	}
	title := strings.TrimSpace(ctx.PostForm("title"))
	if title == "" {
		ctx.JSON(http.StatusBadRequest, serializer.MessageDTO{Message: "нужно название"})
		return
	}

	imageURL := ""
	videoURL := ""
	if file, err := ctx.FormFile("image"); err == nil && file != nil {
		url, upErr := a.Minio.Upload(file, "pattern_image")
		if upErr != nil {
			log.Println(upErr)
			ctx.JSON(http.StatusInternalServerError, serializer.MessageDTO{Message: "не удалось сохранить изображение"})
			return
		}
		imageURL = url
	}
	if file, err := ctx.FormFile("video"); err == nil && file != nil {
		url, upErr := a.Minio.Upload(file, "pattern_video")
		if upErr != nil {
			log.Println(upErr)
			ctx.JSON(http.StatusInternalServerError, serializer.MessageDTO{Message: "не удалось сохранить видео"})
			return
		}
		videoURL = url
	}

	pattern, err := a.Repo.CreateDraftWithMedia(creatorID, title, imageURL, videoURL)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, serializer.MessageDTO{Message: "не удалось создать черновик"})
		return
	}
	ctx.JSON(http.StatusCreated, a.patternDTO(ctx, pattern))
}

type publishBody struct {
	ShortDescription string `json:"short_description" form:"short_description"`
	FullDescription  string `json:"full_description" form:"full_description"`
	EventYear        int    `json:"event_year" form:"event_year"`
	RingCount        int    `json:"ring_count" form:"ring_count"`
}

func (a *API) PublishPattern(ctx *gin.Context) {
	creatorID, err := a.currentCreatorID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, serializer.MessageDTO{Message: "создатель не найден"})
		return
	}

	var body publishBody
	if strings.Contains(ctx.GetHeader("Content-Type"), "application/json") {
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, serializer.MessageDTO{Message: "некорректный JSON"})
			return
		}
	} else {
		_ = ctx.ShouldBind(&body)
	}

	pattern, err := a.Repo.PublishDraft(creatorID, strings.TrimSpace(body.ShortDescription), strings.TrimSpace(body.FullDescription), body.EventYear, body.RingCount)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, serializer.MessageDTO{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, a.patternDTO(ctx, pattern))
}

func (a *API) DeletePattern(ctx *gin.Context) {
	id64, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, serializer.MessageDTO{Message: "некорректный id"})
		return
	}
	creatorID, err := a.currentCreatorID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, serializer.MessageDTO{Message: "исследователь не авторизован"})
		return
	}
	if err := a.Repo.SoftDeletePatternORM(uint(id64), creatorID); err != nil {
		ctx.JSON(http.StatusNotFound, serializer.MessageDTO{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, serializer.MessageDTO{Message: "паттерн удалён"})
}

func (a *API) LikePattern(ctx *gin.Context) {
	id64, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, serializer.MessageDTO{Message: "некорректный id"})
		return
	}
	creatorID, err := a.currentCreatorID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, serializer.MessageDTO{Message: "создатель не найден"})
		return
	}

	likeValue := -1
	if raw := strings.TrimSpace(ctx.PostForm("like")); raw != "" {
		likeValue, _ = strconv.Atoi(raw)
	} else {
		var body struct {
			Like int `json:"like"`
		}
		if err := ctx.ShouldBindJSON(&body); err == nil {
			likeValue = body.Like
		}
	}
	if likeValue != 0 && likeValue != 1 {
		ctx.JSON(http.StatusBadRequest, serializer.MessageDTO{Message: "like должен быть 0 или 1"})
		return
	}

	if err := a.Repo.SetLike(creatorID, uint(id64), likeValue == 1); err != nil {
		ctx.JSON(http.StatusBadRequest, serializer.MessageDTO{Message: err.Error()})
		return
	}
	pattern, err := a.Repo.GetPublishedByID(uint(id64))
	if err != nil {
		ctx.JSON(http.StatusOK, serializer.MessageDTO{Message: "лайк обновлён"})
		return
	}
	ctx.JSON(http.StatusOK, a.patternDTO(ctx, pattern))
}

func (a *API) RegisterResearcher(ctx *gin.Context) {
	var body struct {
		Login    string `json:"login" form:"login"`
		Password string `json:"password" form:"password"`
		FullName string `json:"full_name" form:"full_name"`
	}
	if strings.Contains(ctx.GetHeader("Content-Type"), "application/json") {
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, serializer.MessageDTO{Message: "некорректный JSON"})
			return
		}
	} else {
		_ = ctx.ShouldBind(&body)
	}
	researcher, err := a.Repo.RegisterResearcher(body.Login, body.Password, body.FullName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, serializer.MessageDTO{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, serializer.Researcher(researcher))
}

func (a *API) LoginResearcher(ctx *gin.Context) {
	var body struct {
		Login    string `json:"login" form:"login"`
		Password string `json:"password" form:"password"`
	}
	if strings.Contains(ctx.GetHeader("Content-Type"), "application/json") {
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, serializer.MessageDTO{Message: "некорректный JSON"})
			return
		}
	} else {
		_ = ctx.ShouldBind(&body)
	}
	researcher, err := a.Repo.FindResearcherByLogin(body.Login)
	if err != nil || repository.CheckResearcherPassword(researcher, body.Password) != nil {
		ctx.JSON(http.StatusForbidden, serializer.MessageDTO{Message: "неверный логин или пароль"})
		return
	}
	token, expires, err := a.Tokens.Issue(researcher)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, serializer.MessageDTO{Message: "не удалось выдать токен"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   int(time.Until(expires).Seconds()),
		"login":        researcher.Login,
	})
}

func (a *API) LogoutResearcher(ctx *gin.Context) {
	if err := a.Gate.Revoke(ctx); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, serializer.MessageDTO{Message: "не удалось выйти"})
		return
	}
	ctx.JSON(http.StatusOK, serializer.MessageDTO{Message: "токен добавлен в чёрный список Redis"})
}
