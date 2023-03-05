package tests

import (
	"acourse_tag_cart_bookmark_service/pkg/config"
	"acourse_tag_cart_bookmark_service/pkg/database"
	"acourse_tag_cart_bookmark_service/pkg/database/migrations"
	"acourse_tag_cart_bookmark_service/pkg/http/requests"
	"acourse_tag_cart_bookmark_service/pkg/models"
	"acourse_tag_cart_bookmark_service/pkg/repositories"
	"acourse_tag_cart_bookmark_service/pkg/usecase"
	"context"
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestCartUsecase(t *testing.T) {

	//Create Config Instance
	cfg := config.Construct("../../../.env")

	//Connecting Databases
	db := database.Construct(cfg)
	db.Prepare()

	//Migrations
	mg := migrations.Construct(db)
	mg.MigrateSettings()

	//
	CartDBRepo := repositories.ConstructCartDBRepository(
		db.GetConnection(),
		db.GetCollection(cfg.GetDBConfig()["COLLECTION_CARTS"]))

	cartUsecase := usecase.CartUsecase{DBRepository: CartDBRepo}

	cartID := ""
	userID := "88"
	courseID1 := models.GenerateObjectID()
	courseID2 := models.GenerateObjectID()

	t.Run("AddCourse+", func(t *testing.T) {

		status, err := cartUsecase.AddCourse(context.TODO(), &requests.AddCourseCartRequest{
			UserID:  userID,
			Courses: []requests.Course{{ID: courseID1.Hex()}, {ID: courseID2.Hex()}},
		}, userID)

		assert.Equal(t, err, nil)
		assert.Equal(t, status, true)
	})

	t.Run("FetchByUserId", func(t *testing.T) {

		cart, err := cartUsecase.FetchByUserId(context.TODO(), userID, []string{})
		if err != nil {
			t.Fatal(err)
		}
		cartID = cart.ID.Hex()
		assert.Equal(t, cart.UserID, userID)
	})

	t.Run("FetchById", func(t *testing.T) {

		cart, err := cartUsecase.FetchById(context.TODO(), cartID, []string{})
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, cart.UserID, userID)
	})

	t.Run("AddCourse_WithDuplicateCourseIds", func(t *testing.T) {

		_, err := cartUsecase.AddCourse(context.TODO(), &requests.AddCourseCartRequest{
			UserID:  userID,
			Courses: []requests.Course{{ID: courseID1.Hex()}, {ID: courseID2.Hex()}},
		}, userID)
		if err != nil {
			t.Fatal(err)
		}

		cart, err := cartUsecase.FetchByUserId(context.TODO(), userID, []string{})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, len(cart.Courses), 2)
	})

	t.Run("RevokeCourse+", func(t *testing.T) {

		status, err := cartUsecase.RevokeCourse(context.TODO(), &requests.RevokeCourseCartRequest{
			UserID:  userID,
			Courses: []requests.Course{{ID: courseID1.Hex()}, {ID: courseID2.Hex()}},
		}, userID)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, err, nil)
		assert.Equal(t, status, true)
	})

}
