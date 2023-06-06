package controllers

import (
	"github.com/RayhanAnandhias/realworld-project-golang/pkg/models"
	"github.com/RayhanAnandhias/realworld-project-golang/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ArticleController struct {
	DB *gorm.DB
}

func NewArticleController(DB *gorm.DB) ArticleController {
	return ArticleController{DB}
}

func (ac *ArticleController) CreateArticle(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	var payload *models.ArticleCreateRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	processedSlug := utils.GenerateSlug(payload.Article.Title)

	queryInsert := `INSERT INTO articles (id_author, slug, title, description, body)
						VALUES (?, ?, ?, ?, ?) RETURNING *`

	var article models.Article
	processInsert := ac.DB.Raw(queryInsert, currentUser.ID, processedSlug, payload.Article.Title, payload.Article.Description, payload.Article.Body).Scan(&article)
	if processInsert.Error != nil || article.ID == 0 {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": processInsert.Error.Error()})
		return
	}

	for _, t := range payload.Article.TagList {
		loweredTag := strings.ToLower(strings.TrimSpace(t))

		var tag models.Tag
		queryResultTag := ac.DB.Raw(`SELECT * FROM tags WHERE name = ?`, loweredTag).Scan(&tag)

		// if tag not found then insert into table tags
		if queryResultTag.Error != nil || tag.ID == 0 {
			processInsertTag := ac.DB.Raw(`INSERT INTO tags ("name") VALUES (?) RETURNING *`, loweredTag).Scan(&tag)
			if processInsertTag.Error != nil || article.ID == 0 {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": processInsertTag.Error.Error()})
				return
			}
		}

		processInsertArticleTag := ac.DB.Exec(`INSERT INTO article_tag VALUES (?, ?)`, article.ID, tag.ID)
		if processInsertArticleTag.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": processInsertArticleTag.Error.Error()})
			return
		}
	}

	//find article
	queryFind := `
        SELECT  
          --  a."id", 
		  a."slug", 
		  a."title", 
		  a."description", 
		  a."body", 
		  a."created_at", 
		  a."updated_at", 
		--  a."id_author", 
		  u."username", 
		  u."bio", 
		  u."image", 
		  CASE WHEN f."id_user_a" IS NULL THEN FALSE ELSE TRUE END AS "following", 
		--  att."id_tag" AS "tag_id", 
		  t."name" AS "tag_name", 
		--  l."id_user" AS "liked_by", 
		--  ul."username" AS "liked_by_username", 
		  CASE WHEN l."id_user" IS NULL THEN FALSE ELSE TRUE END AS "favorited",
		  z."favorites_count"
        FROM "articles" AS a 
        INNER JOIN "users" AS u ON u."id" = a."id_author" 
        LEFT JOIN "user_follow" AS f ON f."id_user_a" = ? AND f."id_user_b" = u."id"
        LEFT JOIN "article_tag" AS att ON att."id_article" = a."id" 
        LEFT JOIN "tags" AS t ON t."id" = att."id_tag" 
        LEFT JOIN "user_likes" AS l ON l."id_article" = a."id" 
        LEFT JOIN "users" AS ul ON ul."id" = l."id_user"
        LEFT JOIN (SELECT y."id_article", COUNT(y."id_article") AS "favorites_count" FROM "user_likes" AS y WHERE y."id_article" = ? GROUP BY y."id_article") AS z ON z."id_article" = a."id"
        WHERE a."id" = ? 
        ORDER BY a."created_at" DESC`

	var resultModel []models.ArticleQueryResult
	resultQueryArticle := ac.DB.Raw(queryFind, currentUser.ID, article.ID, article.ID).Scan(&resultModel)
	if resultQueryArticle.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": resultQueryArticle.Error.Error()})
		return
	}

	tagList := make(map[string]bool)

	for _, r := range resultModel {
		tagList[r.TagName] = true
	}

	// Convert map to slice of keys.
	var tagKeys []string
	for key, _ := range tagList {
		tagKeys = append(tagKeys, key)
	}

	uniqueArticleRes := resultModel[0]

	articleResponse := &models.ArticleResponse{
		Article: &models.ArticleCommon{
			Slug:           uniqueArticleRes.Slug,
			Title:          uniqueArticleRes.Title,
			Description:    uniqueArticleRes.Description,
			Body:           uniqueArticleRes.Body,
			TagList:        tagKeys,
			CreatedAt:      uniqueArticleRes.CreatedAt,
			UpdatedAt:      uniqueArticleRes.UpdatedAt,
			Favorited:      uniqueArticleRes.Favorited,
			FavoritesCount: uniqueArticleRes.FavoritesCount,
			Author: &models.UserProfile{
				Username:  uniqueArticleRes.Username,
				Bio:       uniqueArticleRes.Bio,
				Image:     uniqueArticleRes.Image,
				Following: uniqueArticleRes.Following,
			},
		},
	}

	ctx.JSON(http.StatusCreated, articleResponse)

}

func (ac *ArticleController) GetAllArticles(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	tag := ctx.DefaultQuery("tag", "")
	author := ctx.DefaultQuery("author", "")
	favorited := ctx.DefaultQuery("favorited", "")

	if len(tag) == 0 && len(author) == 0 && len(favorited) == 0 {
		tag = "%" + tag + "%"
		author = "%" + author + "%"
		favorited = "%" + favorited + "%"
	} else {
		if len(tag) != 0 {
			tag = "%" + tag + "%"
		}

		if len(author) != 0 {
			author = "%" + author + "%"
		}

		if len(favorited) != 0 {
			favorited = "%" + favorited + "%"
		}
	}

	offset := (limit * page) - limit

	query := `
        SELECT
          a."id",
		  a."slug", 
		  a."title", 
		  a."description", 
		  a."body", 
		  a."created_at", 
		  a."updated_at", 
		  u."username", 
		  u."bio", 
		  u."image", 
		  CASE WHEN f."id_user_a" IS NULL THEN FALSE ELSE TRUE END AS "following", 
		  t."name" AS "tag_name", 
		  CASE WHEN l."id_user" IS NULL THEN FALSE ELSE FALSE END AS "favorited",
		  z."favorites_count"
        FROM "articles" AS a 
        INNER JOIN "users" AS u ON u."id" = a."id_author" 
        LEFT JOIN "user_follow" AS f ON f."id_user_b" = u."id"
		LEFT JOIN "article_tag" AS att ON att."id_article" = a."id" 
        LEFT JOIN "tags" AS t ON t."id" = att."id_tag" 
        LEFT JOIN "user_likes" AS l ON l."id_article" = a."id"
        LEFT JOIN "users" AS ul ON ul."id" = l."id_user"
        LEFT JOIN (SELECT y."id_article", COUNT(y."id_article") AS "favorites_count" FROM "user_likes" AS y GROUP BY y."id_article") AS z ON z."id_article" = a."id"
        WHERE u."username" ILIKE ? OR t."name" ILIKE ? OR ul."username" ILIKE ?
        ORDER BY a."created_at" DESC
        LIMIT ? 
        OFFSET ?`

	var resultModel []models.ArticleQueryResult
	processQuery := ac.DB.Raw(query, author, tag, favorited, limit, offset).Scan(&resultModel)
	if processQuery.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": processQuery.Error.Error()})
		return
	}

	var articleResponseArray []models.ArticleCommon

	uniqueArticles := make(map[int32]models.ArticleQueryResult)

	for _, r := range resultModel {
		uniqueArticles[r.ID] = r
	}

	var articles []models.ArticleQueryResult
	for _, value := range uniqueArticles {
		articles = append(articles, value)
	}

	for _, article := range articles {
		tagList := make(map[string]bool)

		for _, r := range resultModel {
			if article.ID == r.ID {
				tagList[r.TagName] = true
			}
		}

		// Convert map to slice of keys.
		var tagKeys []string
		for key, _ := range tagList {
			tagKeys = append(tagKeys, key)
		}

		articleObj := &models.ArticleCommon{
			Slug:           article.Slug,
			Title:          article.Title,
			Description:    article.Description,
			Body:           article.Body,
			TagList:        tagKeys,
			CreatedAt:      article.CreatedAt,
			UpdatedAt:      article.UpdatedAt,
			Favorited:      article.Favorited,
			FavoritesCount: article.FavoritesCount,
			Author: &models.UserProfile{
				Username:  article.Username,
				Bio:       article.Bio,
				Image:     article.Image,
				Following: article.Following,
			},
		}

		articleResponseArray = append(articleResponseArray, *articleObj)
	}

	ctx.JSON(http.StatusOK, gin.H{"articles": articleResponseArray, "articlesCount": len(articleResponseArray)})
}

func (ac *ArticleController) GetFeedArticles(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	offset := (limit * page) - limit

	query := `
        SELECT
          a."id",
		  a."slug", 
		  a."title", 
		  a."description", 
		  a."body", 
		  a."created_at", 
		  a."updated_at", 
		  u."username", 
		  u."bio", 
		  u."image", 
		  CASE WHEN f."id_user_a" IS NULL THEN FALSE ELSE TRUE END AS "following", 
		  t."name" AS "tag_name", 
		  CASE WHEN l."id_user" IS NULL THEN FALSE ELSE TRUE END AS "favorited",
		  z."favorites_count"
        FROM "articles" AS a 
        INNER JOIN "users" AS u ON u."id" = a."id_author" 
        INNER JOIN "user_follow" AS f ON f."id_user_a" = ? AND f."id_user_b" = u."id"
		LEFT JOIN "article_tag" AS att ON att."id_article" = a."id" 
        LEFT JOIN "tags" AS t ON t."id" = att."id_tag" 
        LEFT JOIN "user_likes" AS l ON l."id_article" = a."id"
        LEFT JOIN "users" AS ul ON ul."id" = l."id_user"
        LEFT JOIN (SELECT y."id_article", COUNT(y."id_article") AS "favorites_count" FROM "user_likes" AS y GROUP BY y."id_article") AS z ON z."id_article" = a."id"
        ORDER BY a."created_at" DESC
        LIMIT ? 
        OFFSET ?`

	var resultModel []models.ArticleQueryResult
	processQuery := ac.DB.Raw(query, currentUser.ID, limit, offset).Scan(&resultModel)
	if processQuery.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": processQuery.Error.Error()})
		return
	}

	var articleResponseArray []models.ArticleCommon

	articleResponseArray = make([]models.ArticleCommon, 0)

	uniqueArticles := make(map[int32]models.ArticleQueryResult)

	for _, r := range resultModel {
		uniqueArticles[r.ID] = r
	}

	var articles []models.ArticleQueryResult
	for _, value := range uniqueArticles {
		articles = append(articles, value)
	}

	for _, article := range articles {
		tagList := make(map[string]bool)

		for _, r := range resultModel {
			if article.ID == r.ID {
				tagList[r.TagName] = true
			}
		}

		// Convert map to slice of keys.
		var tagKeys []string
		for key, _ := range tagList {
			tagKeys = append(tagKeys, key)
		}

		articleObj := &models.ArticleCommon{
			Slug:           article.Slug,
			Title:          article.Title,
			Description:    article.Description,
			Body:           article.Body,
			TagList:        tagKeys,
			CreatedAt:      article.CreatedAt,
			UpdatedAt:      article.UpdatedAt,
			Favorited:      article.Favorited,
			FavoritesCount: article.FavoritesCount,
			Author: &models.UserProfile{
				Username:  article.Username,
				Bio:       article.Bio,
				Image:     article.Image,
				Following: article.Following,
			},
		}

		articleResponseArray = append(articleResponseArray, *articleObj)
	}

	ctx.JSON(http.StatusOK, gin.H{"articles": articleResponseArray, "articlesCount": len(articleResponseArray)})
}

func (ac *ArticleController) GetArticleBySlug(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	slug := ctx.Param("slug")

	query := `
        SELECT
          a."id",
		  a."slug", 
		  a."title", 
		  a."description", 
		  a."body", 
		  a."created_at", 
		  a."updated_at", 
		  u."username", 
		  u."bio", 
		  u."image", 
		  CASE WHEN f."id_user_a" IS NULL THEN FALSE ELSE TRUE END AS "following", 
		  t."name" AS "tag_name", 
		  CASE WHEN l."id_user" IS NULL THEN FALSE ELSE TRUE END AS "favorited",
		  z."favorites_count"
        FROM "articles" AS a 
        INNER JOIN "users" AS u ON u."id" = a."id_author" 
        LEFT JOIN "user_follow" AS f ON f."id_user_a" = ? AND f."id_user_b" = u."id"
		LEFT JOIN "article_tag" AS att ON att."id_article" = a."id" 
        LEFT JOIN "tags" AS t ON t."id" = att."id_tag" 
        LEFT JOIN "user_likes" AS l ON l."id_article" = a."id"
        LEFT JOIN "users" AS ul ON ul."id" = l."id_user"
        LEFT JOIN (SELECT y."id_article", COUNT(y."id_article") AS "favorites_count" FROM "user_likes" AS y GROUP BY y."id_article") AS z ON z."id_article" = a."id"
        WHERE a."slug" = ?`

	var resultModel []models.ArticleQueryResult
	processQuery := ac.DB.Raw(query, currentUser.ID, slug).Scan(&resultModel)
	if processQuery.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": processQuery.Error.Error()})
		return
	} else if len(resultModel) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Data not found"})
		return
	}

	article := resultModel[0]

	tagList := make(map[string]bool)

	for _, r := range resultModel {
		if article.ID == r.ID {
			tagList[r.TagName] = true
		}
	}

	// Convert map to slice of keys.
	var tagKeys []string
	for key, _ := range tagList {
		tagKeys = append(tagKeys, key)
	}

	articleResponse := &models.ArticleCommon{
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		TagList:        tagKeys,
		CreatedAt:      article.CreatedAt,
		UpdatedAt:      article.UpdatedAt,
		Favorited:      article.Favorited,
		FavoritesCount: article.FavoritesCount,
		Author: &models.UserProfile{
			Username:  article.Username,
			Bio:       article.Bio,
			Image:     article.Image,
			Following: article.Following,
		},
	}

	ctx.JSON(http.StatusOK, gin.H{"article": articleResponse})
}

func (ac *ArticleController) UpdateArticle(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	slug := ctx.Param("slug")

	var payload *models.ArticleUpdateRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	querySingleArticle := `SELECT a."id", a."title", a."slug", a."description", a."body" FROM "articles" AS a WHERE a."slug" = ?`
	var oldArticle models.Article
	processGetArticle := ac.DB.Raw(querySingleArticle, slug).Scan(&oldArticle)
	if processGetArticle.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": processGetArticle.Error.Error()})
		return
	} else if oldArticle.ID == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Data not found"})
		return
	}

	titleUpdate := oldArticle.Title
	slugUpdate := oldArticle.Slug
	descriptionUpdate := oldArticle.Description
	bodyUpdate := oldArticle.Body

	if len(payload.Article.Title) != 0 {
		titleUpdate = payload.Article.Title
		slugUpdate = utils.GenerateSlug(payload.Article.Title)
	}

	if len(payload.Article.Description) != 0 {
		descriptionUpdate = payload.Article.Description
	}

	if len(payload.Article.Body) != 0 {
		bodyUpdate = payload.Article.Body
	}

	queryUpdateArticle := `UPDATE articles SET title = ?, slug = ?, description = ?, body = ?, updated_at = ? WHERE id = ? RETURNING id`
	now := time.Now()
	var articleUpdated models.Article
	processUpdate := ac.DB.Raw(queryUpdateArticle, titleUpdate, slugUpdate, descriptionUpdate, bodyUpdate, now, oldArticle.ID).Scan(&articleUpdated)
	if processUpdate.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": processUpdate.Error.Error()})
		return
	} else if articleUpdated.ID == 0 {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Update Failed"})
		return
	}

	// construct Article Response
	query := `
        SELECT
          a."id",
		  a."slug", 
		  a."title", 
		  a."description", 
		  a."body", 
		  a."created_at", 
		  a."updated_at", 
		  u."username", 
		  u."bio", 
		  u."image", 
		  CASE WHEN f."id_user_a" IS NULL THEN FALSE ELSE TRUE END AS "following", 
		  t."name" AS "tag_name", 
		  CASE WHEN l."id_user" IS NULL THEN FALSE ELSE TRUE END AS "favorited",
		  z."favorites_count"
        FROM "articles" AS a 
        INNER JOIN "users" AS u ON u."id" = a."id_author" 
        LEFT JOIN "user_follow" AS f ON f."id_user_a" = ? AND f."id_user_b" = u."id"
		LEFT JOIN "article_tag" AS att ON att."id_article" = a."id" 
        LEFT JOIN "tags" AS t ON t."id" = att."id_tag" 
        LEFT JOIN "user_likes" AS l ON l."id_article" = a."id"
        LEFT JOIN "users" AS ul ON ul."id" = l."id_user"
        LEFT JOIN (SELECT y."id_article", COUNT(y."id_article") AS "favorites_count" FROM "user_likes" AS y GROUP BY y."id_article") AS z ON z."id_article" = a."id"
        WHERE a."id" = ?`

	var resultModel []models.ArticleQueryResult
	processQuery := ac.DB.Raw(query, currentUser.ID, articleUpdated.ID).Scan(&resultModel)
	if processQuery.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": processQuery.Error.Error()})
		return
	} else if len(resultModel) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Data not found"})
		return
	}

	article := resultModel[0]

	tagList := make(map[string]bool)

	for _, r := range resultModel {
		if article.ID == r.ID {
			tagList[r.TagName] = true
		}
	}

	// Convert map to slice of keys.
	var tagKeys []string
	for key, _ := range tagList {
		tagKeys = append(tagKeys, key)
	}

	articleResponse := &models.ArticleCommon{
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		TagList:        tagKeys,
		CreatedAt:      article.CreatedAt,
		UpdatedAt:      article.UpdatedAt,
		Favorited:      article.Favorited,
		FavoritesCount: article.FavoritesCount,
		Author: &models.UserProfile{
			Username:  article.Username,
			Bio:       article.Bio,
			Image:     article.Image,
			Following: article.Following,
		},
	}

	ctx.JSON(http.StatusOK, gin.H{"article": articleResponse})
}

func (ac *ArticleController) FavoriteArticle(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	slug := ctx.Param("slug")

	querySingleArticle := `SELECT a."id" FROM "articles" AS a WHERE a."slug" = ?`
	var oldArticle models.Article
	processGetArticle := ac.DB.Raw(querySingleArticle, slug).Scan(&oldArticle)
	if processGetArticle.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": processGetArticle.Error.Error()})
		return
	} else if oldArticle.ID == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Data not found"})
		return
	}

	query := `INSERT INTO user_likes (id_user, id_article) VALUES (?, ?)`
	process := ac.DB.Exec(query, currentUser.ID, oldArticle.ID)
	if process.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": process.Error.Error()})
		return
	}

	// construct Article Response
	queryArticle := `
        SELECT
          a."id",
		  a."slug", 
		  a."title", 
		  a."description", 
		  a."body", 
		  a."created_at", 
		  a."updated_at", 
		  u."username", 
		  u."bio", 
		  u."image", 
		  CASE WHEN f."id_user_a" IS NULL THEN FALSE ELSE TRUE END AS "following", 
		  t."name" AS "tag_name", 
		  CASE WHEN l."id_user" IS NULL THEN FALSE ELSE TRUE END AS "favorited",
		  z."favorites_count"
        FROM "articles" AS a 
        INNER JOIN "users" AS u ON u."id" = a."id_author" 
        LEFT JOIN "user_follow" AS f ON f."id_user_a" = ? AND f."id_user_b" = u."id"
		LEFT JOIN "article_tag" AS att ON att."id_article" = a."id" 
        LEFT JOIN "tags" AS t ON t."id" = att."id_tag" 
        LEFT JOIN "user_likes" AS l ON l."id_article" = a."id"
        LEFT JOIN "users" AS ul ON ul."id" = l."id_user"
        LEFT JOIN (SELECT y."id_article", COUNT(y."id_article") AS "favorites_count" FROM "user_likes" AS y GROUP BY y."id_article") AS z ON z."id_article" = a."id"
        WHERE a."id" = ?`

	var resultModel []models.ArticleQueryResult
	processQuery := ac.DB.Raw(queryArticle, currentUser.ID, oldArticle.ID).Scan(&resultModel)
	if processQuery.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": processQuery.Error.Error()})
		return
	} else if len(resultModel) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Data not found"})
		return
	}

	article := resultModel[0]

	tagList := make(map[string]bool)

	for _, r := range resultModel {
		if article.ID == r.ID {
			tagList[r.TagName] = true
		}
	}

	// Convert map to slice of keys.
	var tagKeys []string
	for key, _ := range tagList {
		tagKeys = append(tagKeys, key)
	}

	articleResponse := &models.ArticleCommon{
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		TagList:        tagKeys,
		CreatedAt:      article.CreatedAt,
		UpdatedAt:      article.UpdatedAt,
		Favorited:      article.Favorited,
		FavoritesCount: article.FavoritesCount,
		Author: &models.UserProfile{
			Username:  article.Username,
			Bio:       article.Bio,
			Image:     article.Image,
			Following: article.Following,
		},
	}

	ctx.JSON(http.StatusOK, gin.H{"article": articleResponse})

}

func (ac *ArticleController) UnfavoriteArticle(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	slug := ctx.Param("slug")

	querySingleArticle := `SELECT a."id" FROM "articles" AS a WHERE a."slug" = ?`
	var oldArticle models.Article
	processGetArticle := ac.DB.Raw(querySingleArticle, slug).Scan(&oldArticle)
	if processGetArticle.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": processGetArticle.Error.Error()})
		return
	} else if oldArticle.ID == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Data not found"})
		return
	}

	query := `DELETE FROM user_likes WHERE id_user = ? AND id_article = ?`
	process := ac.DB.Exec(query, currentUser.ID, oldArticle.ID)
	if process.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": process.Error.Error()})
		return
	}

	// construct Article Response
	queryArticle := `
        SELECT
          a."id",
		  a."slug", 
		  a."title", 
		  a."description", 
		  a."body", 
		  a."created_at", 
		  a."updated_at", 
		  u."username", 
		  u."bio", 
		  u."image", 
		  CASE WHEN f."id_user_a" IS NULL THEN FALSE ELSE TRUE END AS "following", 
		  t."name" AS "tag_name", 
		  CASE WHEN l."id_user" IS NULL THEN FALSE ELSE TRUE END AS "favorited",
		  z."favorites_count"
        FROM "articles" AS a 
        INNER JOIN "users" AS u ON u."id" = a."id_author" 
        LEFT JOIN "user_follow" AS f ON f."id_user_a" = ? AND f."id_user_b" = u."id"
		LEFT JOIN "article_tag" AS att ON att."id_article" = a."id" 
        LEFT JOIN "tags" AS t ON t."id" = att."id_tag" 
        LEFT JOIN "user_likes" AS l ON l."id_article" = a."id"
        LEFT JOIN "users" AS ul ON ul."id" = l."id_user"
        LEFT JOIN (SELECT y."id_article", COUNT(y."id_article") AS "favorites_count" FROM "user_likes" AS y GROUP BY y."id_article") AS z ON z."id_article" = a."id"
        WHERE a."id" = ?`

	var resultModel []models.ArticleQueryResult
	processQuery := ac.DB.Raw(queryArticle, currentUser.ID, oldArticle.ID).Scan(&resultModel)
	if processQuery.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": processQuery.Error.Error()})
		return
	} else if len(resultModel) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Data not found"})
		return
	}

	article := resultModel[0]

	tagList := make(map[string]bool)

	for _, r := range resultModel {
		if article.ID == r.ID {
			tagList[r.TagName] = true
		}
	}

	// Convert map to slice of keys.
	var tagKeys []string
	for key, _ := range tagList {
		tagKeys = append(tagKeys, key)
	}

	articleResponse := &models.ArticleCommon{
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		TagList:        tagKeys,
		CreatedAt:      article.CreatedAt,
		UpdatedAt:      article.UpdatedAt,
		Favorited:      article.Favorited,
		FavoritesCount: article.FavoritesCount,
		Author: &models.UserProfile{
			Username:  article.Username,
			Bio:       article.Bio,
			Image:     article.Image,
			Following: article.Following,
		},
	}

	ctx.JSON(http.StatusOK, gin.H{"article": articleResponse})

}

func (ac *ArticleController) DeleteArticle(ctx *gin.Context) {
	slug := ctx.Param("slug")

	query := `DELETE FROM articles WHERE slug = ?`

	process := ac.DB.Exec(query, slug)
	if process.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": process.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
