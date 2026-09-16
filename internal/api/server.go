package api

import (
	"fmt"
	"log"
	"net/http"

	"dendro-age/internal/app/config"
	"dendro-age/internal/app/dsn"
	"dendro-age/internal/app/handler"
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

	h := handler.NewHandler(repo)
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

	addr := fmt.Sprintf("%s:%d", cfg.ServiceHost, cfg.ServicePort)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
	log.Println("Server down")
}
