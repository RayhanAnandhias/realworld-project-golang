package controllers

import (
	"fmt"
	"net/http"

	"github.com/RayhanAnandhias/realworld-project-golang/configs"
	"github.com/RayhanAnandhias/realworld-project-golang/pkg/models"
	"github.com/RayhanAnandhias/realworld-project-golang/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(DB *gorm.DB) UserController {
	return UserController{DB}
}

func (uc *UserController) RegisterUser(ctx *gin.Context) {
	var payload *models.UserRegisterRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(payload.User.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	newUser := &models.User{
		Username: payload.User.Username,
		Email:    payload.User.Email,
		Password: hashedPassword,
	}

	result := uc.DB.Raw("INSERT INTO users (username, email, password) VALUES (?, ?, ?) RETURNING *", newUser.Username, newUser.Email, newUser.Password).Scan(&newUser)

	fmt.Println(newUser)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}

	config, _ := configs.LoadConfig(".")

	// Generate Tokens
	accessToken, err := utils.CreateToken(config.AccessTokenExpiresIn, newUser.ID, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// refresh_token, err := utils.CreateToken(config.RefreshTokenExpiresIn, newUser.ID, config.RefreshTokenPrivateKey)
	// if err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
	// 	return
	// }

	// ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	// ctx.SetCookie("refresh_token", refresh_token, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	// ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	userResponse := &models.UserResponse{
		User: &models.UserCommon{
			Email:    newUser.Email,
			Username: newUser.Username,
			Bio:      newUser.Bio,
			Image:    newUser.Image,
			Token:    accessToken,
		},
	}

	ctx.JSON(http.StatusOK, userResponse)
}
