package usecase

import (
	"acourse_tag_cart_bookmark_service/pkg/contracts"
	"acourse_tag_cart_bookmark_service/pkg/http/requests"
	"acourse_tag_cart_bookmark_service/pkg/models"
	"context"
	"errors"
	"log"
	"time"
)

type BookmarkUsecase struct {
	DBRepository            contracts.BookmarksDBRepository
	GRPCCourseServiceClient contracts.GRPCCourseService
}

func (b BookmarkUsecase) Fetch(ctx context.Context, exclude []string, limit int64, skip int64) (bookmarks []models.Bookmark, err error) {
	bookmarks, err = b.DBRepository.Fetch(ctx, exclude, limit, skip)
	if err != nil {
		return nil, err
	}
	return bookmarks, nil
}

func (b BookmarkUsecase) FetchById(ctx context.Context, bookmarkID string, exclude []string) (bookmark models.Bookmark, err error) {

	//Fetch Bookmark Containing embedded course id
	bookmark, err = b.DBRepository.FetchById(ctx, bookmarkID, exclude)
	if err != nil {
		log.Println("BOOKMARK USECASE: FetchById ERROR", err)
		return models.Bookmark{}, err
	}

	//Fetch Course Data From CourseService through GRPC
	cIDs := make([]string, 0)
	for _, c := range bookmark.Courses {
		cIDs = append(cIDs, c.ID.Hex())
	}

	courseResults := b.GRPCCourseServiceClient.List(ctx, cIDs)
	//log.Println("BOOKMARK USECASE: FETCH BY ID: gRPC CourseService Result >>", courseResults)

	//Attach course data from CourseService to a Bookmark
	bookmark.Courses = courseResults

	//log.Println(bookmark)
	return bookmark, nil
}

func (b BookmarkUsecase) FetchByUserId(ctx context.Context, userID string, exclude []string) (bookmark models.Bookmark, err error) {
	bookmark, err = b.DBRepository.FetchByUserId(ctx, userID, exclude)
	if err != nil {
		log.Println("BOOKMARK USECASE: FetchByUserId ERROR >>", err)
		return models.Bookmark{}, err
	}

	//Fetch Course Data From CourseService through GRPC
	cIDs := make([]string, 0)
	for _, c := range bookmark.Courses {
		cIDs = append(cIDs, c.ID.Hex())
	}

	//Attach course data from CourseService to a Bookmark
	courseResults := b.GRPCCourseServiceClient.List(ctx, cIDs)
	bookmark.Courses = courseResults

	return bookmark, nil
}

func (b BookmarkUsecase) Create(ctx context.Context, request *requests.CreateBookmarkRequest) (bookmark models.Bookmark, err error) {

	courses := make([]models.Course, 0)
	for _, course := range request.Courses {
		courses = append(courses, models.Course{ID: b.DBRepository.GenerateObjectIDFromString(course.ID)})
	}

	timeNow := time.Now()
	newBookmark := models.Bookmark{
		ID:        b.DBRepository.GenerateModelID(),
		UserID:    request.UserID,
		Courses:   courses,
		UpdatedAt: &timeNow,
		CreatedAt: &timeNow,
	}

	bookmarkID, err := b.DBRepository.Create(ctx, &newBookmark)
	if err != nil {
		return models.Bookmark{}, err
	}

	newBookmark.ID = bookmarkID
	return newBookmark, nil
}

func (b BookmarkUsecase) AddCourse(ctx context.Context, request *requests.AddCourseBookmarkRequest, userID string) (status bool, err error) {

	//if a bookmark not found, then create a new one
	_, err = b.FetchByUserId(ctx, userID, []string{})
	if err != nil {
		if err.Error() == errors.New("mongo: no documents in result").Error() {
			_, err = b.Create(ctx, (*requests.CreateBookmarkRequest)(request))
			if err != nil {
				log.Println("BOOKMARK USECASE: AddCourse >>", err)
				return false, err
			}
			log.Println("BOOKMARK USECASE: AddCourse >>", "Course has been added", request.Courses)
			return true, nil
		}
		log.Println("BOOKMARK USECASE: AddCourse >>", err)
		return false, err
	}

	cID := make([]string, 0)
	for _, course := range request.Courses {
		cID = append(cID, course.ID)
	}

	status, err = b.DBRepository.AddCourse(ctx, userID, cID)
	if err != nil {
		log.Println("BOOKMARK USECASE: AddCourse >>", err)
		return false, err
	}

	return status, nil
}

func (b BookmarkUsecase) RevokeCourse(ctx context.Context, request *requests.DeleteAttachedCourseRequest, userID string) (status bool, err error) {

	coursesID := make([]string, 0)
	for _, c := range request.Courses {
		coursesID = append(coursesID, c.ID)
	}

	status, err = b.DBRepository.RevokeCourse(ctx, userID, coursesID)
	if err != nil {
		log.Println("BOOKMARK USECASE REVOKE COURSE:", err.Error())
		return false, err
	}

	return status, nil
}

func (b BookmarkUsecase) Delete(ctx context.Context, bookmarkID string) (status bool, err error) {
	status, err = b.DBRepository.Delete(ctx, bookmarkID)
	if err != nil {
		return false, err
	}
	return status, nil
}

func ConstructBookmarkUsecase(DBRepository contracts.BookmarksDBRepository, GRPCCourseServiceClient contracts.GRPCClient) contracts.BookmarkUsecase {
	return &BookmarkUsecase{DBRepository: DBRepository, GRPCCourseServiceClient: GRPCCourseServiceClient}
}
