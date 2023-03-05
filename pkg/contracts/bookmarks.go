package contracts

import (
	"acourse_tag_cart_bookmark_service/pkg/http/requests"
	"acourse_tag_cart_bookmark_service/pkg/models"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookmarksDBRepository interface {
	// Fetch List all data from database;
	// 'exclude' param specify which model fields you want to skip/unselect;
	// 'limit' and 'skip param are used to perform some kind of pagination
	Fetch(ctx context.Context, exclude []string, limit int64, skip int64) (bookmarks []models.Bookmark, err error)
	// FetchById fetch data by id;
	// 'exclude' param specify which model fields you want to skip/unselect;
	FetchById(ctx context.Context, id string, exclude []string) (bookmark models.Bookmark, err error)
	FetchByUserId(ctx context.Context, userID string, exclude []string) (bookmark models.Bookmark, err error)
	Create(ctx context.Context, bookmark *models.Bookmark) (bookmarkID primitive.ObjectID, err error)
	Update(ctx context.Context, bookmark *models.Bookmark, bookmarkID string) (status bool, err error)
	AddCourse(ctx context.Context, userID string, coursesID []string) (status bool, err error)
	Delete(ctx context.Context, bookmarkID string) (status bool, err error)
	RevokeCourse(ctx context.Context, userID string, coursesID []string) (status bool, err error)
	GenerateModelID() primitive.ObjectID
	GenerateObjectIDFromString(id string) primitive.ObjectID
}

type BookmarkUsecase interface {
	// Fetch List all data from database;
	// 'exclude' param specify which model fields you want to skip/unselect;
	// 'limit' and 'skip param are used to perform some kind of pagination
	Fetch(ctx context.Context, exclude []string, limit int64, skip int64) (bookmarks []models.Bookmark, err error)

	// FetchById fetch data by id;
	// 'exclude' param specify which model fields you want to skip/unselect;
	FetchById(ctx context.Context, id string, exclude []string) (bookmark models.Bookmark, err error)

	FetchByUserId(ctx context.Context, userID string, exclude []string) (bookmark models.Bookmark, err error)
	Create(ctx context.Context, request *requests.CreateBookmarkRequest) (bookmark models.Bookmark, err error)
	AddCourse(ctx context.Context, request *requests.AddCourseBookmarkRequest, userID string) (status bool, err error)
	RevokeCourse(ctx context.Context, request *requests.DeleteAttachedCourseRequest, userID string) (status bool, err error)
	Delete(ctx context.Context, bookmarkID string) (status bool, err error)
}
