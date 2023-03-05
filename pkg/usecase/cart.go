package usecase

import (
	"acourse_tag_cart_bookmark_service/pkg/contracts"
	"acourse_tag_cart_bookmark_service/pkg/http/requests"
	"acourse_tag_cart_bookmark_service/pkg/models"
	"context"
	"errors"
	"log"
)

var ErrNoDocuments = errors.New("mongo: no documents in result")

type CartUsecase struct {
	DBRepository            contracts.CartDBRepository
	GRPCCourseServiceClient contracts.GRPCCourseService
}

func (c CartUsecase) FetchById(ctx context.Context, id string, exclude []string) (models.Cart, error) {

	//Fetch a Cart
	cart, err := c.DBRepository.FetchById(ctx, id, exclude)
	if err != nil {
		return models.Cart{}, err
	}

	//Fetch Course Data From CourseService through GRPC
	cIDs := make([]string, 0)
	for _, c := range cart.Courses {
		cIDs = append(cIDs, c.ID.Hex())
	}

	courseResults := c.GRPCCourseServiceClient.List(ctx, cIDs)
	//log.Println("BOOKMARK USECASE: FETCH BY ID: gRPC CourseService Result >>", courseResults)

	//Attach course data from CourseService to a Bookmark
	cart.Courses = courseResults

	return cart, nil

}

func (c CartUsecase) FetchByUserId(ctx context.Context, userID string, exclude []string) (models.Cart, error) {

	//Fetch a Cart
	cart, err := c.DBRepository.FetchByUserId(ctx, userID, exclude)
	if err != nil {
		return models.Cart{}, err
	}

	//Fetch Course Data From CourseService through GRPC
	cIDs := make([]string, 0)
	for _, c := range cart.Courses {
		cIDs = append(cIDs, c.ID.Hex())
	}

	courseResults := c.GRPCCourseServiceClient.List(ctx, cIDs)
	//log.Println("BOOKMARK USECASE: FETCH BY ID: gRPC CourseService Result >>", courseResults)

	//Attach course data from CourseService to a Bookmark
	cart.Courses = courseResults

	return cart, nil
}

func (c CartUsecase) AddCourse(ctx context.Context, request *requests.AddCourseCartRequest, userID string) (status bool, err error) {

	if len(request.Courses) == 0 {
		return false, errors.New("you don't provide any course id, added nothing")
	}

	_, err = c.DBRepository.FetchByUserId(ctx, userID, []string{})
	if err != nil {

		//Create a cart if it doesn't exist yet
		cIDs := make([]models.Course, 0)
		for _, course := range request.Courses {
			cIDs = append(cIDs, models.Course{ID: models.GenerateObjectIDFromHex(course.ID)})
		}

		if err.Error() == ErrNoDocuments.Error() {

			_, err := c.DBRepository.Create(ctx, &models.Cart{
				ID:      models.GenerateObjectID(),
				UserID:  userID,
				Courses: cIDs,
			})

			if err != nil {
				log.Println("CART USECASE: AddCourse: Create Cart >>", err)
				return false, err
			}

			return true, nil
		}

		log.Println("CART USECASE: AddCourse: Create Cart >>", err)
		return false, err
	}

	//if it does exist then add course_id to existing cart
	cIDs := make([]string, 0)
	for _, course := range request.Courses {
		cIDs = append(cIDs, course.ID)
	}

	status, err = c.DBRepository.AddCourse(ctx, userID, cIDs)
	if err != nil {
		log.Println("CART USECASE: AddCourse: Add Courses >>", err)
		return false, err
	}

	return true, nil
	//c.DBRepository.AddCourse(ctx, userID)
}

func (c CartUsecase) RevokeCourse(ctx context.Context, request *requests.RevokeCourseCartRequest, userID string) (bool, error) {

	cIDs := make([]string, 0)
	if len(request.Courses) == 0 {
		return false, errors.New("you don't provide any course id, nothing removed")
	}

	for _, course := range request.Courses {
		cIDs = append(cIDs, course.ID)
	}

	status, err := c.DBRepository.RevokeCourse(ctx, userID, cIDs)
	if err != nil {
		return false, err
	}

	return status, nil

}

func ConstructCartUsecase(DBRepository contracts.CartDBRepository, grpcCourseService contracts.GRPCCourseService) contracts.CartUsecase {
	return &CartUsecase{DBRepository: DBRepository, GRPCCourseServiceClient: grpcCourseService}
}
