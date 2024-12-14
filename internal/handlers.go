package internal

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/sajjadvaezi/face-recognition/internal/services"
	"github.com/sajjadvaezi/face-recognition/models"
	"log/slog"
	"net/http"
)

func CheckHealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "fine")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	type registerRequest struct {
		Name       string `json:"name"`
		UserNumber string `json:"user_number"`
		Role       string `json:"role"`
	}

	type registerResponse struct {
		ID int64
	}
	response := registerResponse{}
	request := registerRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	createdUserID, err := services.CreateUser(request.Name, request.UserNumber, request.Role)
	if err != nil {
		return
	}
	response.ID = createdUserID
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println(err.Error())

		return
	}

}

func RecognizeHandler(w http.ResponseWriter, r *http.Request) {

	type recognizeResponse struct {
		Name string `json:"name"`
	}

	name, err := services.RecognizeFace()
	if err != nil {
		fmt.Println(err.Error())

		return
	}
	recRes := recognizeResponse{
		Name: name,
	}
	err = json.NewEncoder(w).Encode(&recRes)
	if err != nil {
		fmt.Println(err.Error())

		return
	}

}

func AddFaceHandler(w http.ResponseWriter, r *http.Request) {
	type addFaceRequest struct {
		UserNumber string `json:"user_number"`
	}

	type response struct {
		Result string `json:"result"`
		Error  string `json:"error,omitempty"`
	}

	// Decode the request body
	var afReq addFaceRequest
	err := json.NewDecoder(r.Body).Decode(&afReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		fmt.Println("Decoding error:", err.Error())
		return
	}

	// Validate input
	if afReq.UserNumber == "" {
		http.Error(w, "StudentNumber is required", http.StatusBadRequest)
		fmt.Println("Validation error: StudentNumber is empty")
		return
	}

	// Call the service to add the face
	err = services.AddFace(afReq.UserNumber)
	if err != nil {
		http.Error(w, "Failed to add face to database", http.StatusInternalServerError)
		fmt.Println("Service error:", err.Error())
		return
	}

	// Respond with success
	res := response{
		Result: "successful",
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		fmt.Println("Encoding response error:", err.Error())
	}
}

func RecognizeWithImageHandler(w http.ResponseWriter, r *http.Request) {
	type ImageUploadRequest struct {
		UserNumber string `json:"user_number"`
		Image      string `json:"image"`
	}

	type RecognizeResponse struct {
		Name  string `json:"name,omitempty"`
		Error string `json:"error,omitempty"`
	}

	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var uploadReq ImageUploadRequest
	var response RecognizeResponse

	err := json.NewDecoder(r.Body).Decode(&uploadReq)
	if err != nil {
		response.Error = "invalid request body"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	name, err := services.RecognizeFaceWithImage(uploadReq.UserNumber, uploadReq.Image)
	if err != nil {
		response.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	response.Name = name

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func RegisterFaceWithImageHandler(w http.ResponseWriter, r *http.Request) {
	type ImageUploadRequest struct {
		UserNumber string `json:"user_number"`
		Image      string `json:"image"`
	}

	type RegisterResponse struct {
		Error string `json:"error"`
	}

	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var uploadReq ImageUploadRequest
	var response RegisterResponse

	err := json.NewDecoder(r.Body).Decode(&uploadReq)
	if err != nil {
		response.Error = "invalid request body"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err = services.AddFaceWithImage(uploadReq.UserNumber, uploadReq.Image)
	if err != nil {
		response.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Success response

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func AddClassHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AddClassRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res := models.AddClassResponse{
			Status:     "failed",
			StatusCode: http.StatusBadRequest,
			Error:      "Invalid request body",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		slog.Error("Failed to decode request body", "error", err)
		return
	}

	if err := services.AddClass(req); err != nil {
		res := models.AddClassResponse{
			Status:     "failed",
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		slog.Error("Failed to add class", "error", err, "class_name", req.ClassName, "user_number", req.UserNumber)
		return
	}

	res := models.AddClassResponse{
		Status:     "success",
		StatusCode: http.StatusOK,
		Error:      "null",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		slog.Error("Failed to encode response", "error", err)
	}
}

func AttendanceHandler(w http.ResponseWriter, r *http.Request) {
	atReq := models.AttendanceRequest{}

	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var response models.AttendanceClassResponse

	err := json.NewDecoder(r.Body).Decode(&atReq)
	if err != nil {
		response.Error = "invalid request body"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err = services.Attendance(atReq)
	if err != nil {
		response.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	response = models.AttendanceClassResponse{
		Status:     "success",
		StatusCode: http.StatusOK,
		Error:      "null",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func AttendedUsersHandler(c fiber.Ctx) error {
	// Get the class name from the request parameters
	className := c.Params("classname")

	// Check if the class name is empty
	if className == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "class name is required",
		})
	}

	// Call the service layer to get the list of attended users
	users, err := services.AttendedUsers(className)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve attended users",
		})
	}

	// Transform the data to only include user number and name
	type UserResponse struct {
		UserNumber string `json:"user_number"`
		Name       string `json:"name"`
	}

	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			UserNumber: user.UserNumber,
			Name:       user.Name,
		})
	}

	// Return the transformed list of users as a JSON response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"users": userResponses,
	})
}
