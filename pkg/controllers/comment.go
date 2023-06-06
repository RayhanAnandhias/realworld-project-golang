package controllers

import (
	"github.com/RayhanAnandhias/realworld-project-golang/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type CommentController struct {
	DB *gorm.DB
}

func NewCommentController(DB *gorm.DB) CommentController {
	return CommentController{DB}
}

func (cc *CommentController) CreateComment(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	slug := ctx.Param("slug")

	var payload *models.CommentCreateRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	querySingleArticle := `SELECT a."id" FROM "articles" AS a WHERE a."slug" = ?`
	var oldArticle models.Article
	processGetArticle := cc.DB.Raw(querySingleArticle, slug).Scan(&oldArticle)
	if processGetArticle.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": processGetArticle.Error.Error()})
		return
	} else if oldArticle.ID == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Data not found"})
		return
	}

	queryInsert := `INSERT INTO comments (id_author, id_article, body) VALUES (?, ?, ?) RETURNING *`

	var comment models.Comment
	process := cc.DB.Raw(queryInsert, currentUser.ID, oldArticle.ID, payload.Comment.Body).Scan(&comment)
	if process.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": process.Error.Error()})
		return
	}

	queryComment := `
		SELECT 
			c.id, 
			c.body, 
			c.created_at, 
			c.updated_at, 
			c.id_author, 
			c.id_article, 
			u.username, 
			u.bio, 
			u.image, 
			CASE WHEN f."id_user_a" IS NULL THEN FALSE ELSE TRUE END AS "following"
		FROM "comments" as c 
		INNER JOIN "users" AS u ON u."id" = c.id_author 
		LEFT JOIN "user_follow" AS f ON f."id_user_a" = ? AND f."id_user_b" = u."id" 
		WHERE c.id = ?`

	var commentResponse models.CommentQueryResult
	processComment := cc.DB.Raw(queryComment, currentUser.ID, comment.ID).Scan(&commentResponse)
	if processComment.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": processComment.Error.Error()})
		return
	}

	commentResponseObject := &models.CommentResponseData{
		Comment: &models.CommentResponse{
			ID:        commentResponse.ID,
			CreatedAt: commentResponse.CreatedAt,
			UpdatedAt: commentResponse.UpdatedAt,
			Body:      commentResponse.Body,
			Author: &models.UserProfile{
				Username:  commentResponse.Username,
				Bio:       commentResponse.Bio,
				Image:     commentResponse.Image,
				Following: commentResponse.Following,
			},
		},
	}

	ctx.JSON(http.StatusCreated, commentResponseObject)
}

func (cc *CommentController) GetCommentsForArticle(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	slug := ctx.Param("slug")

	querySingleArticle := `SELECT a."id" FROM "articles" AS a WHERE a."slug" = ?`
	var oldArticle models.Article
	processGetArticle := cc.DB.Raw(querySingleArticle, slug).Scan(&oldArticle)
	if processGetArticle.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": processGetArticle.Error.Error()})
		return
	} else if oldArticle.ID == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Data not found"})
		return
	}

	queryComments := `
		SELECT 
			c.id, 
			c.body, 
			c.created_at, 
			c.updated_at, 
			c.id_author, 
			c.id_article, 
			u.username, 
			u.bio, 
			u.image, 
			CASE WHEN f."id_user_a" IS NULL THEN FALSE ELSE TRUE END AS "following"
		FROM "comments" as c 
		INNER JOIN "users" AS u ON u."id" = c.id_author 
		LEFT JOIN "user_follow" AS f ON f."id_user_a" = ? AND f."id_user_b" = u."id" 
		WHERE c.id_article = ?
		ORDER BY c.created_at ASC`

	resultQuery := make([]models.CommentQueryResult, 0)
	process := cc.DB.Raw(queryComments, currentUser.ID, oldArticle.ID).Scan(&resultQuery)
	if process.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": process.Error.Error()})
		return
	}

	commentResponseArray := make([]models.CommentResponse, 0)

	for _, comment := range resultQuery {
		commentObj := &models.CommentResponse{
			ID:        comment.ID,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			Body:      comment.Body,
			Author: &models.UserProfile{
				Username:  comment.Username,
				Bio:       comment.Bio,
				Image:     comment.Image,
				Following: comment.Following,
			},
		}

		commentResponseArray = append(commentResponseArray, *commentObj)
	}

	ctx.JSON(http.StatusOK, gin.H{"comments": commentResponseArray})
}

func (cc *CommentController) DeleteCommentForArticle(ctx *gin.Context) {
	commentId := ctx.Param("commentId")

	commentIdInt, _ := strconv.Atoi(commentId)

	query := `DELETE FROM comments WHERE id = ?`

	process := cc.DB.Exec(query, commentIdInt)
	if process.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": process.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
