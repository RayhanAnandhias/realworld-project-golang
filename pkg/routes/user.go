package routes

import (
	"github.com/RayhanAnandhias/realworld-project-golang/pkg/controllers"
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
}
