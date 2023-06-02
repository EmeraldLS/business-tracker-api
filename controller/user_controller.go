package controller

import (
	"fmt"
	"net/http"

	"github.com/EmeraldLS/phsps-api/config"
	"github.com/EmeraldLS/phsps-api/encryption"
	"github.com/EmeraldLS/phsps-api/model"
	"github.com/EmeraldLS/phsps-api/token"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-module/carbon"
)

func Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "bind_error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "struct_error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}

	hashwd, err := encryption.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "pass_hash_error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	user.Password = hashwd
	user.JoinDate = carbon.Now().ToDateTimeString()
	user.UpdatedDate = carbon.Now().ToDateTimeString()
	if err := config.ValidateEmail(user.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "email_error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	if err := config.CheckEmailExist(user.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "email_error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}

	token, refreshToken, expiresAt, err, claim := token.GenerateToken(user.FirstName, user.LastName, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "token_error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	user.Token = token
	user.RefreshToken = refreshToken
	user.ExpiresAt = expiresAt

	if err := config.Register(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "registration_error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}

	c.SetCookie("token", token, int(expiresAt), "", "", false, true)
	c.JSON(http.StatusCreated, gin.H{
		"response": "registration_successful",
		"user":     claim,
	})

}

func Login(c *gin.Context) {
	var user model.Login
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "bind_error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"respnse": "struct_error",
			"message": fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	userDetail, err := config.Login(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}

	if err := encryption.ValidatePassword(user.Password, userDetail.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "pass_error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	signedToken, refreshToken, expirationTime, err, claim := token.GenerateToken(userDetail.FirstName, userDetail.LastName, userDetail.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "token_error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	updatedAt := carbon.Now().ToDateTimeString()
	_, err = token.UpdateToken(signedToken, refreshToken, updatedAt, expirationTime, userDetail.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "update_error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	c.SetCookie("token", signedToken, int(expirationTime), "", "", false, true)
	c.JSON(http.StatusOK, claim)
}

func Logout(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"response": "token_error",
			"message":  "no authorization token",
		})
		c.Abort()
		return
	}
	count, err := config.Logout(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"response": "logout_error",
			"message":  fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}
	fmt.Println("Modified Count:", count)
	c.JSON(http.StatusOK, gin.H{
		"response": "success",
		"message":  "logout successful.",
	})
}
