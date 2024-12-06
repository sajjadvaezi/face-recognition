package internal

import (
	"encoding/json"
	"fmt"
	"github.com/sajjadvaezi/face-recognition/internal/services"
	"net/http"
)

func CheckHealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "fine")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	type registerRequest struct {
		Name          string `json:"name"`
		StudentNumber string `json:"student_number"`
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
	createdUserID, err := services.CreateUser(request.Name, request.StudentNumber)
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
		StudentNumber string `json:"student_number"`
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
	if afReq.StudentNumber == "" {
		http.Error(w, "StudentNumber is required", http.StatusBadRequest)
		fmt.Println("Validation error: StudentNumber is empty")
		return
	}

	// Call the service to add the face
	err = services.AddFace(afReq.StudentNumber)
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
		StudentNumber string `json:"student_number"`
		Image         string `json:"image"`
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

	name, err := services.RecognizeFaceWithImage(uploadReq.StudentNumber, uploadReq.Image)
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
		StudentNumber string `json:"student_number"`
		Image         string `json:"image"`
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

	err = services.AddFaceWithImage(uploadReq.StudentNumber, uploadReq.Image)
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
