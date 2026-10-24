package auth

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"dendro-age/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ContextResearcherID = "researcher_id"
	ContextAccessToken  = "access_token"
	bearerPrefix        = "Bearer "
)

type Claims struct {
	jwt.RegisteredClaims
	ResearcherID uint   `json:"researcher_id"`
	Login        string `json:"login"`
}

type Tokens struct {
	secret []byte
	ttl    time.Duration
}

func NewTokensFromEnv() *Tokens {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dendro-age-lab4-secret"
	}
	hours := 24
	if raw := os.Getenv("JWT_HOURS"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			hours = parsed
		}
	}
	return &Tokens{secret: []byte(secret), ttl: time.Duration(hours) * time.Hour}
}

func (t *Tokens) TTL() time.Duration {
	return t.ttl
}

func (t *Tokens) Issue(researcher ds.DendroResearcher) (string, time.Time, error) {
	expires := time.Now().Add(t.ttl)
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "dendro-age",
			Subject:   strconv.FormatUint(uint64(researcher.ID), 10),
		},
		ResearcherID: researcher.ID,
		Login:        researcher.Login,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(t.secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expires, nil
}

func (t *Tokens) Parse(raw string) (Claims, error) {
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, bearerPrefix) {
		raw = strings.TrimSpace(raw[len(bearerPrefix):])
	}
	if raw == "" {
		return Claims{}, errors.New("пустой токен")
	}
	parsed, err := jwt.ParseWithClaims(raw, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("неверный алгоритм")
		}
		return t.secret, nil
	})
	if err != nil {
		return Claims{}, err
	}
	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return Claims{}, errors.New("токен недействителен")
	}
	return *claims, nil
}

func Bearer(header string) string {
	header = strings.TrimSpace(header)
	if strings.HasPrefix(header, bearerPrefix) {
		return strings.TrimSpace(header[len(bearerPrefix):])
	}
	return ""
}

func ResearcherID(ctx *gin.Context) (uint, bool) {
	value, ok := ctx.Get(ContextResearcherID)
	if !ok {
		return 0, false
	}
	id, ok := value.(uint)
	return id, ok && id > 0
}

func AccessToken(ctx *gin.Context) string {
	value, _ := ctx.Get(ContextAccessToken)
	token, _ := value.(string)
	return token
}
