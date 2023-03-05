package controllers

import (
	"acourse_tag_cart_bookmark_service/pkg/contracts"
	"github.com/gin-gonic/gin"
)

func SetupHandler(router *gin.Engine, bookmarkUsecase *contracts.BookmarkUsecase, cartUsecase *contracts.CartUsecase) {
	bookmarkHandler := BookmarkHandler{BookmarkUsecase: *bookmarkUsecase}
	cartHandler := CartHandler{CartUsecase: *cartUsecase}

	bRoute := router.Group("/bookmark")
	bRoute.GET("/", bookmarkHandler.Fetch)
	bRoute.GET("/:id", bookmarkHandler.FetchById)
	bRoute.GET("/u/:user_id", bookmarkHandler.FetchByUserID)
	//bRoute.POST("/create", bookmarkHandler.Create)
	bRoute.DELETE("/course/delete/:user_id", bookmarkHandler.RevokeCourse)
	bRoute.PATCH("/course/add/:user_id", bookmarkHandler.AddCourse)

	cRoute := router.Group("/cart")
	cRoute.GET("/:id", cartHandler.FetchByID)
	cRoute.GET("/u/:user_id", cartHandler.FetchByUserID)
	cRoute.PATCH("/course/add/:user_id", cartHandler.AddCourse)
	cRoute.DELETE("/course/revoke/:user_id", cartHandler.RevokeCourse)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

}
