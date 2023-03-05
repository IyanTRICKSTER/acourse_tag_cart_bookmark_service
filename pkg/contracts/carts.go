package contracts

import (
	"acourse_tag_cart_bookmark_service/pkg/http/requests"
	"acourse_tag_cart_bookmark_service/pkg/models"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartDBRepository interface {
	// FetchById fetch data by id;
	// 'exclude' param specify which model fields you want to skip/unselect;
	FetchById(ctx context.Context, id string, exclude []string) (cart models.Cart, err error)
	FetchByUserId(ctx context.Context, userID string, exclude []string) (cart models.Cart, err error)
	Create(ctx context.Context, cart *models.Cart) (cartId primitive.ObjectID, err error)
	AddCourse(ctx context.Context, userID string, coursesID []string) (status bool, err error)
	RevokeCourse(ctx context.Context, userID string, coursesID []string) (status bool, err error)
	Delete(ctx context.Context, cartID string) (status bool, err error)
}

type CartUsecase interface {
	FetchById(ctx context.Context, id string, exclude []string) (cart models.Cart, err error)
	FetchByUserId(ctx context.Context, userID string, exclude []string) (cart models.Cart, err error)
	AddCourse(ctx context.Context, request *requests.AddCourseCartRequest, userID string) (status bool, err error)
	RevokeCourse(ctx context.Context, request *requests.RevokeCourseCartRequest, userID string) (status bool, err error)
}
