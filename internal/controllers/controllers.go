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

	// hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "Failed to hash body",
	// 	})
	// }

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

	if c.BindHeader(&token) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "can't read header",
		})
		return
	}

	history, err := r.History(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal error",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, history)
}

func (r *authRoutes) deleteHistory(c *gin.Context) {
	// take token

	// auth validate token

	// response status
}
