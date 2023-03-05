package tests

import (
	"acourse_tag_cart_bookmark_service/pkg/models"
	"acourse_tag_cart_bookmark_service/pkg/repositories"
	"context"
	"github.com/go-playground/assert/v2"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
)

func TestMongo(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("model id generator", func(mt *mtest.T) {
		db := mt.Client.Database("acourse")
		coll := mt.Coll

		bookmarkDBRepo := repositories.ConstructBookmarkDBRepository(db, coll)

		objectID := bookmarkDBRepo.GenerateModelID()
		expected := bookmarkDBRepo.GenerateObjectIDFromString(objectID.Hex())

		assert.Equal(t, objectID.String(), expected.String())
	})

	mt.Run("insert data", func(mt *mtest.T) {
		db := mt.Client.Database("acourse")
		coll := mt.Coll

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		bookmarkDBRepo := repositories.ConstructBookmarkDBRepository(db, coll)

		bookmarkID := bookmarkDBRepo.GenerateModelID()
		resID, err := bookmarkDBRepo.Create(context.TODO(), &models.Bookmark{
			ID:        bookmarkID,
			UserID:    "132",
			Courses:   []models.Course{{ID: bookmarkDBRepo.GenerateModelID()}},
			UpdatedAt: nil,
			CreatedAt: nil,
			DeletedAt: nil,
		})
		if err != nil {
			return
		}

		assert.Equal(t, bookmarkID, resID)
		assert.Equal(t, nil, err)
	})

	mt.Run("fetch data", func(mt *mtest.T) {
		db := mt.Client.Database("acourse")
		coll := mt.Coll

		bookmarkDBRepo := repositories.ConstructBookmarkDBRepository(db, coll)
		//_, _ = bookmarkDBRepo.Create(context.TODO(), &models.Bookmark{
		//	ID:        bookmarkDBRepo.GenerateModelID(),
		//	UserID:    "132",
		//	Courses:   []models.Course{{ID: bookmarkDBRepo.GenerateModelID()}},
		//	UpdatedAt: nil,
		//	CreatedAt: nil,
		//	DeletedAt: nil,
		//})

		bookmarks, err := bookmarkDBRepo.Fetch(context.TODO(), []string{}, 10, 0)
		if err != nil {
			return
		}

		assert.Equal(t, bookmarks, []models.Bookmark{})
	})

	mt.Run("find by id", func(mt *mtest.T) {
		db := mt.Client.Database("acourse")
		coll := mt.Coll

		bookmarkDBRepo := repositories.ConstructBookmarkDBRepository(db, coll)

		id := bookmarkDBRepo.GenerateModelID().Hex()
		bookmark, _ := bookmarkDBRepo.FetchById(context.TODO(), id, []string{})

		//assert.Equal(t, err, mongo.ErrNoDocuments)
		assert.Equal(t, bookmark, models.Bookmark{})
	})

	mt.Run("delete course", func(mt *mtest.T) {
		db := mt.Client.Database("acourse")
		coll := mt.Coll

		bookmarkDBRepo := repositories.ConstructBookmarkDBRepository(db, coll)

		status, err := bookmarkDBRepo.RevokeCourse(context.TODO(), bookmarkDBRepo.GenerateModelID().Hex(), []string{"123", "456"})
		if err != nil {
			return
		}

		t.Log(err)
		t.Log(status)
	})
}
