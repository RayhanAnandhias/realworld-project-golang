package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

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

	refreshToken, err := utils.CreateToken(config.RefreshTokenExpiresIn, newUser.ID, config.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", accessToken, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	userResponse := &models.UserResponse{
		User: &models.UserCommon{
			Email:    newUser.Email,
			Username: newUser.Username,
			Bio:      newUser.Bio,
			Image:    newUser.Image,
			Token:    accessToken,
		},
	}

	ctx.JSON(http.StatusCreated, userResponse)
}

func (uc *UserController) LoginUser(ctx *gin.Context) {
	var payload *models.UserLoginRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := uc.DB.Raw("SELECT * FROM users WHERE email = ?", strings.ToLower(payload.User.Email)).Scan(&user)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	if err := utils.VerifyPassword(user.Password, payload.User.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	config, _ := configs.LoadConfig(".")

	// Generate Tokens
	accessToken, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refreshToken, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.ID, config.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", accessToken, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	userResponse := &models.UserResponse{
		User: &models.UserCommon{
			Email:    user.Email,
			Username: user.Username,
			Bio:      user.Bio,
			Image:    user.Image,
			Token:    accessToken,
		},
	}

	ctx.JSON(http.StatusOK, userResponse)
}

func (uc *UserController) RefreshToken(ctx *gin.Context) {
	message := "could not refresh access token"

	cookie, err := ctx.Cookie("refresh_token")

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	config, _ := configs.LoadConfig(".")

	sub, err := utils.ValidateToken(cookie, config.RefreshTokenPublicKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := uc.DB.Raw("SELECT * FROM users WHERE id = ?", fmt.Sprint(sub)).Scan(&user)
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no longer exists"})
		return
	}

	// Generate Tokens
	accessToken, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refreshToken, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.ID, config.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", accessToken, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	userResponse := &models.UserResponse{
		User: &models.UserCommon{
			Email:    user.Email,
			Username: user.Username,
			Bio:      user.Bio,
			Image:    user.Image,
			Token:    accessToken,
		},
	}

	ctx.JSON(http.StatusOK, userResponse)
}

func (uc *UserController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (uc *UserController) GetCurrentUser(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	token := ctx.MustGet("token").(string)

	userResponse := &models.UserResponse{
		User: &models.UserCommon{
			Email:    currentUser.Email,
			Username: currentUser.Username,
			Bio:      currentUser.Bio,
			Image:    currentUser.Image,
			Token:    token,
		},
	}

	ctx.JSON(http.StatusOK, userResponse)
}

func (uc *UserController) UpdateCurrentUser(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	token := ctx.MustGet("token").(string)

	var payload *models.UserUpdateRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	hashedPassword := currentUser.Password

	if payload.User.Password != nil {
		var err error
		hashedPassword, err = utils.HashPassword(*(payload.User.Password))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
	}

	newEmail := currentUser.Email
	if payload.User.Email != nil {
		newEmail = *(payload.User.Email)
	}

	newUsername := currentUser.Username
	if payload.User.Username != nil {
		newUsername = *(payload.User.Username)
	}

	newBio := currentUser.Bio
	if payload.User.Bio != nil {
		newBio = payload.User.Bio
	}

	newImage := currentUser.Image
	if payload.User.Image != nil {
		newImage = payload.User.Image
	}

	now := time.Now()
	var updatedUser models.User
	result := uc.DB.Raw("UPDATE users SET email = ?, password = ?, username = ?, bio = ?, image = ?, updated_at = ? WHERE id = ? RETURNING *", strings.ToLower(newEmail), hashedPassword, newUsername, utils.NewNullString(newBio), utils.NewNullString(newImage), now, currentUser.ID).Scan(&updatedUser)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}

	userResponse := &models.UserResponse{
		User: &models.UserCommon{
			Email:    updatedUser.Email,
			Username: updatedUser.Username,
			Bio:      updatedUser.Bio,
			Image:    updatedUser.Image,
			Token:    token,
		},
	}

	ctx.JSON(http.StatusOK, userResponse)
}

func (uc *UserController) GetProfile(ctx *gin.Context) {
	profileUsername := ctx.Param("profileUsername")
	currentUser := ctx.MustGet("currentUser").(models.User)

	query :=
		`SELECT 
			u.username, 
			u.bio, 
			u.image, 
			CASE WHEN f.id_user_a IS NULL THEN FALSE ELSE TRUE END AS following 
		FROM users AS u 
		LEFT JOIN user_follow AS f ON f.id_user_a = ? AND f.id_user_b = u.id 
		WHERE u.username = ?`

	var profile models.UserProfile
	queryProfileResult := uc.DB.Raw(query, currentUser.ID, profileUsername).Scan(&profile)
	if queryProfileResult.Error != nil || len(profile.Username) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "fail", "message": "profile not found"})
		return
	}

	userProfileResponse := &models.UserProfileResponse{
		Profile: profile,
	}

	ctx.JSON(http.StatusOK, userProfileResponse)
}

func (uc *UserController) FollowUser(ctx *gin.Context) {
	profileUsername := ctx.Param("profileUsername")
	currentUser := ctx.MustGet("currentUser").(models.User)

	queryInsert := `INSERT INTO user_follow VALUES (?, ?)`
	queryUser := `SELECT id FROM users WHERE username = ?`
	queryProfile :=
		`SELECT 
			u.username, 
			u.bio, 
			u.image, 
			CASE WHEN f.id_user_a IS NULL THEN FALSE ELSE TRUE END AS following 
		FROM users AS u 
		LEFT JOIN user_follow AS f ON f.id_user_a = ? AND f.id_user_b = u.id 
		WHERE u.username = ?`

	var user models.User
	queryUserResult := uc.DB.Raw(queryUser, profileUsername).Scan(&user)
	if queryUserResult.Error != nil || user.ID == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "fail", "message": "profile not found"})
		return
	}

	execFollow := uc.DB.Exec(queryInsert, currentUser.ID, user.ID)
	if execFollow.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": execFollow.Error.Error()})
		return
	}

	var profile models.UserProfile
	queryProfileResult := uc.DB.Raw(queryProfile, currentUser.ID, profileUsername).Scan(&profile)
	if queryProfileResult.Error != nil || len(profile.Username) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "fail", "message": "profile not found"})
		return
	}

	userProfileResponse := &models.UserProfileResponse{
		Profile: profile,
	}

	ctx.JSON(http.StatusOK, userProfileResponse)
}

func (uc *UserController) UnfollowUser(ctx *gin.Context) {
	profileUsername := ctx.Param("profileUsername")
	currentUser := ctx.MustGet("currentUser").(models.User)

	queryDelete := `DELETE FROM user_follow WHERE id_user_a = ? AND id_user_b = ?`
	queryUser := `SELECT id FROM users WHERE username = ?`
	queryProfile :=
		`SELECT 
			u.username, 
			u.bio, 
			u.image, 
			CASE WHEN f.id_user_a IS NULL THEN FALSE ELSE TRUE END AS following 
		FROM users AS u 
		LEFT JOIN user_follow AS f ON f.id_user_a = ? AND f.id_user_b = u.id 
		WHERE u.username = ?`

	var user models.User
	queryUserResult := uc.DB.Raw(queryUser, profileUsername).Scan(&user)
	if queryUserResult.Error != nil || user.ID == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "fail", "message": "profile not found"})
		return
	}

	execUnfollow := uc.DB.Exec(queryDelete, currentUser.ID, user.ID)
	if execUnfollow.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": execUnfollow.Error.Error()})
		return
	}

	var profile models.UserProfile
	queryProfileResult := uc.DB.Raw(queryProfile, currentUser.ID, profileUsername).Scan(&profile)
	if queryProfileResult.Error != nil || len(profile.Username) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "fail", "message": "profile not found"})
		return
	}

	userProfileResponse := &models.UserProfileResponse{
		Profile: profile,
	}

	ctx.JSON(http.StatusOK, userProfileResponse)
}
