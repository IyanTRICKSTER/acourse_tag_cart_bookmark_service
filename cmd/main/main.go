package main

import (
	"acourse_tag_cart_bookmark_service/cmd/grpc_client"
	"acourse_tag_cart_bookmark_service/pkg/config"
	"acourse_tag_cart_bookmark_service/pkg/database"
	"acourse_tag_cart_bookmark_service/pkg/database/migrations"
	"acourse_tag_cart_bookmark_service/pkg/http/controllers"
	"acourse_tag_cart_bookmark_service/pkg/repositories"
	"acourse_tag_cart_bookmark_service/pkg/usecase"
	"github.com/gin-gonic/gin"
)

func main() {

	engine := gin.Default()

	//Create Config Instance
	cfg := config.Construct(".env")

	//Connecting Databases
	db := database.Construct(cfg)
	db.Prepare()

	//Migrations
	mg := migrations.Construct(db)
	mg.MigrateSettings()

	//Setup Bookmarks
	//Repo
	bookmarkRepo := repositories.ConstructBookmarkDBRepository(
		db.GetConnection(),
		db.GetCollection(cfg.GetDBConfig()["COLLECTION_BOOKMARKS"]),
	)

	cartRepo := repositories.ConstructCartDBRepository(db.GetConnection(), db.GetCollection(cfg.GetDBConfig()["COLLECTION_CARTS"]))

	//Connect to Course Service via GRPC
	grpcCourseService := grpc_client.Construct(cfg)
	_, err := grpcCourseService.Dial()
	if err != nil {
		panic(err)
	}

	bookmarkUsecase := usecase.ConstructBookmarkUsecase(bookmarkRepo, grpcCourseService)
	cartUsecase := usecase.ConstructCartUsecase(cartRepo, grpcCourseService)

	//Setup Delivery/Controller
	controllers.SetupHandler(engine, &bookmarkUsecase, &cartUsecase)

	if port := cfg.GetAppConfig()["PORT"]; port == "" {
		err := engine.Run(":8080")
		if err != nil {
			panic(err)
		}
	} else {
		err := engine.Run(":" + port)
		if err != nil {
			panic(err)
		}
	}

}
