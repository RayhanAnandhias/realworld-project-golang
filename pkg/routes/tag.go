package routes

import (
	"github.com/RayhanAnandhias/realworld-project-golang/pkg/controllers"
	"github.com/gin-gonic/gin"
)

type TagRouteController struct {
	tagController controllers.TagController
}

func NewTagRouteController(tagController controllers.TagController) TagRouteController {
	return TagRouteController{tagController}
}

func (trc *TagRouteController) TagRoute(rg *gin.RouterGroup) {

	router := rg.Group("tags")
	router.GET("/", trc.tagController.GetTags)
}
