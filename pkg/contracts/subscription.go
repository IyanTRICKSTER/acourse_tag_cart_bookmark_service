package contracts

import (
	"acourse_tag_cart_bookmark_service/pkg/models"
	"context"
)

type SubscriptionDBRepository interface {
	// FetchById fetch data by id;
	// 'exclude' param specify which model fields you want to skip/unselect;
	FetchById(ctx context.Context, id string, exclude []string) (cart models.Subscription, err error)
	FetchByUserId(ctx context.Context, userID string, exclude []string) (cart models.Subscription, err error)
	Subscribe(ctx context.Context, coursesID []string) (err error)
	Unsubscribe(ctx context.Context, coursesID []string) (err error)
}
