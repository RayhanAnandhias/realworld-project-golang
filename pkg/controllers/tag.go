package controllers

import (
	"github.com/RayhanAnandhias/realworld-project-golang/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type TagController struct {
	DB *gorm.DB
}

func NewTagController(DB *gorm.DB) TagController {
	return TagController{DB}
}

func (tc *TagController) GetTags(ctx *gin.Context) {
	var rawTags []models.Tag
	var normalizedTags []string

	tc.DB.Raw("SELECT * FROM tags").Scan(&rawTags)

	for _, t := range rawTags {
		normalizedTags = append(normalizedTags, t.Name)
	}

	tagResponse := &models.TagResponse{
		Tags: normalizedTags,
	}

	ctx.JSON(http.StatusOK, tagResponse)
}
