package api

import (
	"log"

	"dendro-age/internal/app/handler"
	"dendro-age/internal/app/repository"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	log.Println("Server start up")

	repo, err := repository.NewRepository()
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
	r.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(302, "/dendro_scale_patterns/feed")
	})

	if err := r.Run("localhost:8088"); err != nil {
		log.Fatal(err)
	}

	log.Println("Server down")
}
