package requests

type AddCourseCartRequest struct {
	UserID  string   `json:"user_id" binding:"required"`
	Courses []Course `json:"courses" binding:"required,dive"`
}

type RevokeCourseCartRequest struct {
	UserID  string   `json:"user_id" binding:"required"`
	Courses []Course `json:"courses" binding:"required,dive"`
}
