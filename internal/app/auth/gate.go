package auth

import (
	"context"
	"net/http"
	"time"

	"dendro-age/internal/app/redisblk"
	"dendro-age/internal/app/serializer"

	"github.com/gin-gonic/gin"
)

type Gate struct {
	tokens *Tokens
	block  *redisblk.Blacklist
}

func NewGate(tokens *Tokens, block *redisblk.Blacklist) *Gate {
	return &Gate{tokens: tokens, block: block}
}

func (g *Gate) Optional(ctx *gin.Context) {
	raw := Bearer(ctx.GetHeader("Authorization"))
	if raw == "" {
		ctx.Next()
		return
	}
	claims, err := g.tokens.Parse(raw)
	if err != nil {
		ctx.Next()
		return
	}
	blocked, err := g.block.Has(context.Background(), raw)
	if err != nil || blocked {
		ctx.Next()
		return
	}
	ctx.Set(ContextResearcherID, claims.ResearcherID)
	ctx.Set(ContextAccessToken, raw)
	ctx.Next()
}

func (g *Gate) Required(ctx *gin.Context) {
	raw := Bearer(ctx.GetHeader("Authorization"))
	if raw == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, serializer.MessageDTO{Message: "нужен заголовок Authorization: Bearer"})
		return
	}
	claims, err := g.tokens.Parse(raw)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, serializer.MessageDTO{Message: "токен недействителен"})
		return
	}
	blocked, err := g.block.Has(context.Background(), raw)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, serializer.MessageDTO{Message: "не удалось проверить blacklist"})
		return
	}
	if blocked {
		ctx.AbortWithStatusJSON(http.StatusForbidden, serializer.MessageDTO{Message: "токен в чёрном списке, войдите снова"})
		return
	}
	ctx.Set(ContextResearcherID, claims.ResearcherID)
	ctx.Set(ContextAccessToken, raw)
	ctx.Next()
}

func (g *Gate) Revoke(ctx *gin.Context) error {
	raw := AccessToken(ctx)
	if raw == "" {
		raw = Bearer(ctx.GetHeader("Authorization"))
	}
	claims, err := g.tokens.Parse(raw)
	if err != nil {
		return err
	}
	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl < time.Second {
		ttl = time.Second
	}
	return g.block.Add(context.Background(), raw, ttl)
}
