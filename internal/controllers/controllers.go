package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *authRoutes) signUp(c *gin.Context) {
	var body struct {
		Login    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	err := r.SignUp(c.Request.Context(), body.Login, body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprint(err),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "successfully created",
	})

}

func (r *authRoutes) signIn(c *gin.Context) {
	var body struct {
		Login    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	token, err := r.SignIn(c.Request.Context(), body.Login, body.Password)
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprint(err),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
	})
}

func (r *authRoutes) history(c *gin.Context) {
	token := c.Request.Header.Get("X-Token")

	history, err := r.History(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.IndentedJSON(http.StatusOK, history)
}

func (r *authRoutes) deleteHistory(c *gin.Context) {
	token := c.Request.Header.Get("X-Token")

	err := r.Clear(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": " history is deleted",
	})
}
