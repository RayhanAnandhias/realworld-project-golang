package routes

import (
	"github.com/RayhanAnandhias/realworld-project-golang/pkg/controllers"
	"github.com/RayhanAnandhias/realworld-project-golang/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

type ArticleRouteController struct {
	ArticleController controllers.ArticleController
	CommentController controllers.CommentController
}

func NewArticleRouteController(ArticleController controllers.ArticleController, CommentController controllers.CommentController) ArticleRouteController {
	return ArticleRouteController{ArticleController, CommentController}
}

func (arc *ArticleRouteController) ArticleRoute(rg *gin.RouterGroup) {
	router := rg.Group("articles")
	router.POST("/", middlewares.DeserializeUser(), arc.ArticleController.CreateArticle)
	router.GET("/", arc.ArticleController.GetAllArticles)
	router.GET("/feed", middlewares.DeserializeUser(), arc.ArticleController.GetFeedArticles)
	router.GET("/:slug", middlewares.DeserializeUser(), arc.ArticleController.GetArticleBySlug)
	router.PUT("/:slug", middlewares.DeserializeUser(), arc.ArticleController.UpdateArticle)
	router.POST("/:slug/favorite", middlewares.DeserializeUser(), arc.ArticleController.FavoriteArticle)
	router.DELETE("/:slug/favorite", middlewares.DeserializeUser(), arc.ArticleController.UnfavoriteArticle)
	router.POST("/:slug/comments", middlewares.DeserializeUser(), arc.CommentController.CreateComment)
	router.GET("/:slug/comments", middlewares.DeserializeUser(), arc.CommentController.GetCommentsForArticle)
	router.DELETE("/:slug/comments/:commentId", middlewares.DeserializeUser(), arc.CommentController.DeleteCommentForArticle)
	router.DELETE("/:slug", middlewares.DeserializeUser(), arc.ArticleController.DeleteArticle)
}
