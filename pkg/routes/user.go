package routes

import (
	"github.com/RayhanAnandhias/realworld-project-golang/pkg/controllers"
	"github.com/RayhanAnandhias/realworld-project-golang/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewUserRouteController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (urc *UserRouteController) UserRoute(rg *gin.RouterGroup) {

	router := rg.Group("users")
	router.POST("/", urc.userController.RegisterUser)
	router.POST("/login", urc.userController.LoginUser)
	router.POST("/logout", urc.userController.LogoutUser)
}

func (urc *UserRouteController) SingleUserRoute(rg *gin.RouterGroup) {
	router := rg.Group("user")
	router.GET("/", middlewares.DeserializeUser(), urc.userController.GetCurrentUser)
	router.PUT("/", middlewares.DeserializeUser(), urc.userController.UpdateCurrentUser)
}

func (urc *UserRouteController) ProfileRoute(rg *gin.RouterGroup) {
	router := rg.Group("profiles")
	router.GET("/:profileUsername", middlewares.DeserializeUser(), urc.userController.GetProfile)
	router.POST("/:profileUsername/follow", middlewares.DeserializeUser(), urc.userController.FollowUser)
	router.DELETE("/:profileUsername/follow", middlewares.DeserializeUser(), urc.userController.UnfollowUser)
}
