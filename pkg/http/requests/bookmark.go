package requests

type CreateBookmarkRequest struct {
	UserID  string   `json:"user_id" binding:"required"`
	Courses []Course `json:"courses" binding:"required,dive"`
}

type AddCourseBookmarkRequest struct {
	UserID  string   `json:"user_id" binding:"required"`
	Courses []Course `json:"courses" binding:"required,dive"`
}

type DeleteAttachedCourseRequest struct {
	UserID  string   `json:"user_id" binding:"required"`
	Courses []Course `json:"courses" binding:"required,dive"`
}

type Course struct {
	ID string `json:"id" binding:"required"`
}
