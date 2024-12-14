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
		Status     string `json:"status"`       // Indicates "success" or "failed"
		ID         int64  `json:"id,omitempty"` // User ID in case of success
		Error      string `json:"error,omitempty"`
		StatusCode int    `json:"status_code"`
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

	var req registerRequest
	var res registerResponse

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res = registerResponse{
			Status:     "failed",
			Error:      "Invalid request payload",
			StatusCode: http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// Validate input
	if req.Name == "" || req.UserNumber == "" || req.Role == "" {
		res = registerResponse{
			Status:     "failed",
			Error:      "name, user_number, and role are required fields",
			StatusCode: http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// Create the user
	createdUserID, err := services.CreateUser(req.Name, req.UserNumber, req.Role)
	if err != nil {
		// Log the detailed error for debugging
		fmt.Printf("Error creating user: %v\n", err)

		// Return a generic error message to the client
		res = registerResponse{
			Status:     "failed",
			Error:      "An unexpected error occurred. Please try again later.",
			StatusCode: http.StatusInternalServerError,
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	// Respond with success
	res = registerResponse{
		Status:     "success",
		ID:         createdUserID,
		StatusCode: http.StatusOK,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func RecognizeHandler(w http.ResponseWriter, r *http.Request) {
	type recognizeResponse struct {
		Status     string `json:"status"` // Indicates "success" or "failed"
		Name       string `json:"name,omitempty"`
		Error      string `json:"error,omitempty"`
		StatusCode int    `json:"status_code"`
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

	var res recognizeResponse

	// Call the service to recognize the face
	name, err := services.RecognizeFace()
	if err != nil {
		// Log the detailed error for debugging
		fmt.Printf("Error recognizing face: %v\n", err)

		// Return a generic error message to the client
		res = recognizeResponse{
			Status:     "failed",
			Error:      "An unexpected error occurred. Please try again later.",
			StatusCode: http.StatusInternalServerError,
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	// If face recognition is successful, respond with the name
	res = recognizeResponse{
		Status:     "success",
		Name:       name,
		StatusCode: http.StatusOK,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func AddFaceHandler(w http.ResponseWriter, r *http.Request) {
	type addFaceRequest struct {
		UserNumber string `json:"user_number"`
	}

	type response struct {
		Status     string `json:"status"` // Indicates "success" or "failed"
		Result     string `json:"result,omitempty"`
		Error      string `json:"error,omitempty"`
		StatusCode int    `json:"status_code"`
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

	var afReq addFaceRequest
	var res response

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&afReq); err != nil {
		res = response{
			Status:     "failed",
			Error:      "Invalid request payload",
			StatusCode: http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// Validate input
	if afReq.UserNumber == "" {
		res = response{
			Status:     "failed",
			Error:      "user_number is required",
			StatusCode: http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// Call the service to add the face
	err := services.AddFace(afReq.UserNumber)
	if err != nil {
		res = response{
			Status:     "failed",
			Error:      err.Error(), // Provide the error for debugging or logging
			StatusCode: http.StatusInternalServerError,
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	// Respond with success
	res = response{
		Status:     "success",
		Result:     "Face added successfully",
		StatusCode: http.StatusOK,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
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
		Status     string `json:"status"` // Indicates "success" or "failed"
		Error      string `json:"error,omitempty"`
		StatusCode int    `json:"status_code"`
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
	var res RegisterResponse

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&uploadReq); err != nil {
		res = RegisterResponse{
			Status:     "failed",
			Error:      "Invalid request body",
			StatusCode: http.StatusBadRequest,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// Process the image upload
	err := services.AddFaceWithImage(uploadReq.UserNumber, uploadReq.Image)
	if err != nil {
		// Log the error for debugging
		fmt.Printf("Error registering face: %v\n", err)

		// Return a detailed error message
		res = RegisterResponse{
			Status:     "failed",
			Error:      err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	// Success response
	res = RegisterResponse{
		Status:     "success",
		StatusCode: http.StatusOK,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
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
