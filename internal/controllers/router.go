package controllers

import (
	"jwt-auth/internal/auth"

	"github.com/gin-gonic/gin"
)

type authRoutes struct {
	auth.Authentification
}

func NewRouter(handler *gin.Engine, a auth.Authentification) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	r := &authRoutes{a}

	h := handler.Group("/auth")
	{
		h.POST("/signup", r.signUp)
		h.POST("/signin", r.signIn)
		h.GET("/history", r.history)
		h.GET("/clear", r.deleteHistory)
	}
}
