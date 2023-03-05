package controllers

import (
	"acourse_tag_cart_bookmark_service/pkg/contracts"
	"acourse_tag_cart_bookmark_service/pkg/http/requests"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type CartHandler struct {
	CartUsecase contracts.CartUsecase
}

func (h CartHandler) FetchByUserID(c *gin.Context) {

	if c.Param("user_id") == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id is not provided in url parameter",
		})
		return
	}

	cart, err := h.CartUsecase.FetchByUserId(c.Request.Context(), c.Param("user_id"), []string{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "document not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, cart)
	return
}

func (h CartHandler) FetchByID(c *gin.Context) {

	if c.Param("id") == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id is not provided in url parameter",
		})
		return
	}

	cart, err := h.CartUsecase.FetchById(c.Request.Context(), c.Param("id"), []string{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "document not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, cart)
	return
}

func (h CartHandler) AddCourse(c *gin.Context) {

	if c.Param("user_id") == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id is not provided in url parameter",
		})
		return
	}

	var addCourseReq requests.AddCourseCartRequest
	err := c.ShouldBindJSON(&addCourseReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	status, err := h.CartUsecase.AddCourse(c.Request.Context(), &addCourseReq, c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if !status {
		c.JSON(http.StatusOK, gin.H{
			"status":  status,
			"message": "failed to add listed course",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
	return
}

func (h CartHandler) RevokeCourse(c *gin.Context) {

	if c.Param("user_id") == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id is not provided in url parameter",
		})
		return
	}

	var revokeCourseReq requests.RevokeCourseCartRequest
	err := c.ShouldBindJSON(&revokeCourseReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	status, err := h.CartUsecase.RevokeCourse(c.Request.Context(), &revokeCourseReq, c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if !status {
		c.JSON(http.StatusOK, gin.H{
			"status":  status,
			"message": "failed to revoke listed course",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
	return

}
