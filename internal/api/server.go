package api

import (
	"fmt"
	"log"
	"net/http"

	"dendro-age/internal/app/apihandler"
	"dendro-age/internal/app/auth"
	"dendro-age/internal/app/config"
	"dendro-age/internal/app/dsn"
	"dendro-age/internal/app/handler"
	"dendro-age/internal/app/miniofile"
	"dendro-age/internal/app/redisblk"
	"dendro-age/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func StartServer() {
	log.Println("Server start up")
	_ = godotenv.Load()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	repo, err := repository.New(dsn.FromEnv())
	if err != nil {
		log.Fatal(err)
	}

	minioClient, err := miniofile.NewFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	block, err := redisblk.NewFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	tokens := auth.NewTokensFromEnv()
	gate := auth.NewGate(tokens, block)

	h := handler.NewHandler(repo)
	api := apihandler.New(repo, minioClient, tokens, gate)
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/dendro_scale_patterns", h.GetDendroScalePatterns)
	r.GET("/dendro_scale_patterns/feed", h.GetDendroScalePatternsFeed)
	r.GET("/dendro_scale_patterns/feed/:id", h.GetDendroScalePatternsFeed)
	r.GET("/dendro_scale_patterns/draft", h.GetDendroScalePatternsDraft)
	r.GET("/dendro_scale_patterns/:id", h.GetDendroScalePatternsFeed)
	r.POST("/dendro_scale_patterns/draft/next", h.PostDendroScalePatternNext)
	r.POST("/dendro_scale_patterns/draft/publish", h.PostDendroScalePatternPublish)
	r.POST("/dendro_scale_patterns/:id/remove", h.PostDendroScalePatternDelete)
	r.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusFound, "/dendro_scale_patterns/feed")
	})

	r.GET("/swagger", func(ctx *gin.Context) {
		ctx.File("./docs/swagger.html")
	})
	r.StaticFile("/swagger/openapi.json", "./docs/openapi.json")

	apiGroup := r.Group("/api")
	publicRead := apiGroup.Group("", gate.Optional)
	{
		publicRead.GET("/dendro_scale_patterns", api.ListPatterns)
		publicRead.GET("/dendro_scale_patterns/feed", api.GetFeed)
		publicRead.GET("/dendro_scale_patterns/feed/:id", api.GetFeed)
	}

	apiGroup.POST("/dendro_researchers/register", api.RegisterResearcher)
	apiGroup.POST("/dendro_researchers/login", api.LoginResearcher)

	creator := apiGroup.Group("", gate.Required)
	{
		creator.GET("/dendro_scale_patterns/draft", api.GetDraft)
		creator.POST("/dendro_scale_patterns", api.CreatePattern)
		creator.PUT("/dendro_scale_patterns/publish", api.PublishPattern)
		creator.DELETE("/dendro_scale_patterns/:id", api.DeletePattern)
		creator.POST("/dendro_scale_patterns/:id/like", api.LikePattern)
		creator.POST("/dendro_researchers/logout", api.LogoutResearcher)
	}

	addr := fmt.Sprintf("%s:%d", cfg.ServiceHost, cfg.ServicePort)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
	log.Println("Server down")
}
