package controllers

import (
	"acourse_tag_cart_bookmark_service/pkg/contracts"
	"acourse_tag_cart_bookmark_service/pkg/http/requests"
	"acourse_tag_cart_bookmark_service/pkg/http/responses"
	"acourse_tag_cart_bookmark_service/pkg/models"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type BookmarkHandler struct {
	BookmarkUsecase contracts.BookmarkUsecase
}

func (h *BookmarkHandler) Fetch(c *gin.Context) {

	var excludedField []string
	if c.Query("exclude") != "" {
		excludedField = strings.Split(c.Query("exclude"), ",")
	}

	page, ok := c.GetQuery("page")
	if page == "" || !ok {
		page = "1"
	} else if page == "0" {
		page = "1"
	}

	qPage, err2 := strconv.ParseInt(page, 10, 64)
	if err2 != nil {
		return
	}

	paginate := models.Pagination{
		Page:    qPage,
		PerPage: 25,
	}

	limit, skip := paginate.GetPagination()
	bookmarks, err := h.BookmarkUsecase.Fetch(c.Request.Context(), excludedField, limit, skip)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, responses.HttpPaginationResponse{
		PerPage: paginate.PerPage,
		Page:    paginate.Page,
		HttpResponse: responses.HttpResponse{
			Data:       bookmarks,
			StatusCode: http.StatusOK,
		},
	})
}

func (h BookmarkHandler) FetchById(c *gin.Context) {
	bookmark, err := h.BookmarkUsecase.FetchById(c.Request.Context(), c.Param("id"), []string{})
	if err != nil {
		//return 404 not found
		if err.Error() == mongo.ErrNoDocuments.Error() {
			c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookmark)
}

func (h BookmarkHandler) FetchByUserID(c *gin.Context) {
	bookmark, err := h.BookmarkUsecase.FetchByUserId(c.Request.Context(), c.Param("user_id"), []string{})
	if err != nil {
		//return 404 not found
		if err.Error() == mongo.ErrNoDocuments.Error() {
			c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookmark)
}

func (h BookmarkHandler) Create(c *gin.Context) {

	var createRequest requests.CreateBookmarkRequest

	err := c.ShouldBindJSON(&createRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	bookmark, err := h.BookmarkUsecase.Create(c.Request.Context(), &createRequest)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookmark)
	return
}

func (h BookmarkHandler) AddCourse(c *gin.Context) {

	var addCourse requests.AddCourseBookmarkRequest

	err := c.ShouldBindJSON(&addCourse)
	if err != nil {
		log.Println("BOOKMARK HANDLER: AddCourse", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = h.BookmarkUsecase.AddCourse(c.Request.Context(), &addCourse, c.Param("user_id"))
	if err != nil {
		//return 404 not found
		if err.Error() == mongo.ErrNoDocuments.Error() {
			log.Println("BOOKMARK HANDLER: AddCourse", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
			return
		}
		//return error on failed updates
		if err.Error() == errors.New("document not modified").Error() {
			log.Println("BOOKMARK HANDLER: AddCourse", err)
			c.JSON(http.StatusConflict, gin.H{"error": "duplication occurs"})
			return
		}
		log.Println("BOOKMARK HANDLER: AddCourse", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
	return

}

func (h BookmarkHandler) RevokeCourse(c *gin.Context) {

	var revokeCourse requests.DeleteAttachedCourseRequest

	err := c.ShouldBindJSON(&revokeCourse)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = h.BookmarkUsecase.RevokeCourse(c.Request.Context(), &revokeCourse, c.Param("user_id"))
	if err != nil {
		//return 404 not found
		if err.Error() == mongo.ErrNoDocuments.Error() {
			c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
			return
		}
		if err.Error() == errors.New("document not modified").Error() {
			c.JSON(http.StatusNotFound, gin.H{"error": "course id not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
	return
}
