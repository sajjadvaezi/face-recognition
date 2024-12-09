package models

type Class struct {
	Classname string
	TeacherID int
	ClassID   int
}

type AddClassRequest struct {
	ClassName  string `json:"class_name"`
	UserNumber string `json:"user_number"`
}
type AddClassResponse struct {
	Status     string `json:"status"`
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
}
