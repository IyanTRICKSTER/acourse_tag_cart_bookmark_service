package tests

import (
	"acourse_tag_cart_bookmark_service/pkg/config"
	"acourse_tag_cart_bookmark_service/pkg/database"
	"acourse_tag_cart_bookmark_service/pkg/database/migrations"
	"acourse_tag_cart_bookmark_service/pkg/models"
	"acourse_tag_cart_bookmark_service/pkg/repositories"
	"context"
	"github.com/go-playground/assert/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestCartDBRepo(t *testing.T) {

	//Create Config Instance
	cfg := config.Construct("../../.env")

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

	timeNow := time.Now()

	//Positive
	cartID := models.GenerateObjectID()
	userId := strconv.Itoa(rand.Int())
	courseID1 := models.GenerateObjectID()
	courseID2 := models.GenerateObjectID()

	coursesData := []models.Course{{
		ID:   courseID1,
		Name: "Kelas Test",
	}}

	cartData := &models.Cart{
		ID:        cartID,
		UserID:    userId,
		Courses:   coursesData,
		UpdatedAt: &timeNow,
		CreatedAt: &timeNow,
		DeletedAt: nil,
	}

	t.Run("CreateCart+", func(t *testing.T) {

		createdCartId, err := CartDBRepo.Create(context.TODO(), cartData)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, cartID, createdCartId)
	})

	t.Run("FetchById+", func(t *testing.T) {

		cart, err := CartDBRepo.FetchById(context.TODO(), cartID.Hex(), []string{})
		if err != nil {
			t.Fatal(err.Error())
		}

		assert.Equal(t, cart.ID, cartID)
	})

	t.Run("FetchByUserID+", func(t *testing.T) {
		cart, err := CartDBRepo.FetchByUserId(context.TODO(), userId, []string{})
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, cart.ID, cartID)
	})

	t.Run("AddCourseToACart+", func(t *testing.T) {

		status, err := CartDBRepo.AddCourse(context.TODO(), userId, []string{courseID2.Hex()})
		if err != nil {
			t.Fatal(err)
		}

		cart, err := CartDBRepo.FetchByUserId(context.TODO(), userId, []string{})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, status, true)
		assert.Equal(t, cart.Courses[1].ID, courseID2)
	})

	t.Run("AddCourseToACart_WithInvalidUserID-", func(t *testing.T) {
		status, err := CartDBRepo.AddCourse(context.TODO(), "invaliduserid", []string{"ivalidhexid1", "invalidhexid2"})
		if err == nil {
			t.Fatal("Something went wrong! this should raises error no document in result")
		}
		assert.Equal(t, status, false)
		assert.Equal(t, err, mongo.ErrNoDocuments)

	})

	t.Run("AddCourseToACart_WithInvalidCourseID-", func(t *testing.T) {
		status, err := CartDBRepo.AddCourse(context.TODO(), userId, []string{"ivalidhexid1", "invalidhexid2"})
		if err != nil {
			t.Fatal(err)
		}

		cart, err := CartDBRepo.FetchByUserId(context.TODO(), userId, []string{})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, status, true)
		assert.Equal(t, len(cart.Courses), 2)
	})

	t.Run("RevokeCourseFromCart+", func(t *testing.T) {
		status, err := CartDBRepo.RevokeCourse(context.TODO(), userId, []string{courseID1.Hex()})
		if err != nil {
			t.Fatal(err)
		}

		cart, err := CartDBRepo.FetchByUserId(context.TODO(), userId, []string{})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, status, true)
		assert.Equal(t, len(cart.Courses), 1)
		assert.Equal(t, cart.Courses[0].ID, courseID2)
	})

	t.Run("RevokeCourse_WithInvalidUserId-", func(t *testing.T) {
		status, err := CartDBRepo.RevokeCourse(context.TODO(), "invaliduserid", []string{})
		if err == nil {
			t.Fatal("This should raises error no document in result")
		}
		assert.Equal(t, status, false)
		assert.Equal(t, err, mongo.ErrNoDocuments)

	})

	t.Run("RevokeCourse_WithNoExistsCourseID-", func(t *testing.T) {

		status, _ := CartDBRepo.RevokeCourse(context.TODO(), userId, []string{courseID1.Hex()})

		assert.Equal(t, status, true)

	})

	t.Run("DeleteByID+", func(t *testing.T) {

		status, err := CartDBRepo.Delete(context.TODO(), cartID.Hex())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, status, true)

	})

	//Negative
	t.Run("CreateCartWithDuplicationUserID-", func(t *testing.T) {

		_, err := CartDBRepo.Create(context.TODO(), &models.Cart{
			ID:     models.GenerateObjectID(),
			UserID: "88",
			Courses: []models.Course{{
				ID:   models.GenerateObjectID(),
				Name: "Kelas Test",
			}},
			UpdatedAt: &timeNow,
			CreatedAt: &timeNow,
			DeletedAt: nil,
		})

		_, err = CartDBRepo.Create(context.TODO(), &models.Cart{
			ID:     models.GenerateObjectID(),
			UserID: "88",
			Courses: []models.Course{{
				ID:   models.GenerateObjectID(),
				Name: "Kelas Test",
			}},
			UpdatedAt: &timeNow,
			CreatedAt: &timeNow,
			DeletedAt: nil,
		})

		if err == nil {
			t.Error("Something Went Wrong, Duplication of user id occurs")
		}

		assert.Equal(t, true, mongo.IsDuplicateKeyError(err))

	})

	t.Run("FetchById_NoExistDocument-", func(t *testing.T) {

		_, err := CartDBRepo.FetchById(context.TODO(), cartID.Hex(), []string{})

		assert.Equal(t, err, mongo.ErrNoDocuments)
	})

	t.Run("FetchById_WithInvalidHexID-", func(t *testing.T) {

		_, err := CartDBRepo.FetchById(context.TODO(), "terekjkdfjdfhd", []string{})

		assert.Equal(t, err, primitive.ErrInvalidHex)
	})

	t.Run("FetchByUserID-", func(t *testing.T) {
		_, err := CartDBRepo.FetchByUserId(context.TODO(), userId, []string{})
		if err == nil {
			t.Fatal("Something went wrong, user id not match and should raises error")
		}
		assert.Equal(t, err, mongo.ErrNoDocuments)
	})

	t.Run("FetchById_WithInvalidHex", func(t *testing.T) {
		_, err := CartDBRepo.FetchById(context.TODO(), "thisisarandomhexformat", []string{})

		assert.Equal(t, err, primitive.ErrInvalidHex)
	})

	t.Run("Delete_NoExistsDocument-", func(t *testing.T) {
		status, err := CartDBRepo.Delete(context.TODO(), models.GenerateObjectID().Hex())

		assert.Equal(t, status, false)
		assert.Equal(t, err, mongo.ErrNoDocuments)
	})

	t.Run("Delete_WithInvalidHex-", func(t *testing.T) {
		_, err := CartDBRepo.Delete(context.TODO(), "thisisarandomhexformat")

		assert.Equal(t, err, primitive.ErrInvalidHex)
	})

}
