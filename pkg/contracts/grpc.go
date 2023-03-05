package contracts

import (
	"acourse_tag_cart_bookmark_service/pkg/models"
	ps "acourse_tag_cart_bookmark_service/pkg/models/proto_schema"
	"context"
)

type GRPCClient interface {
	Dial() (ps.CoursesServiceClient, error)
	List(ctx context.Context, coursesID []string) []models.Course
}

type GRPCCourseService interface {
	GRPCClient
}
