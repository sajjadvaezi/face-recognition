package internal

import (
	"encoding/json"
	"fmt"
	"github.com/sajjadvaezi/face-recognition/internal/services"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	type registerRequest struct {
		Name string `json:"name"`
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
	createdUserID, err := services.CreateUser(request.Name)
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
	type recognizeRequest struct {
		Hash string `json:"hash"`
	}

	type recognizeResponse struct {
		Name string `json:"name"`
	}
	recReq := recognizeRequest{}
	err := json.NewDecoder(r.Body).Decode(&recReq)
	if err != nil {
		fmt.Println(err.Error())

		return
	}
	name, err := services.RecognizeFace(recReq.Hash)
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
		UserId int    `json:"id"`
		Hash   string `json:"hash"`
	}
	afReq := addFaceRequest{}
	err := json.NewDecoder(r.Body).Decode(&afReq)
	if err != nil {
		fmt.Println("Decoding error ", err.Error())

		return
	}

	err = services.AddFace(int64(afReq.UserId), afReq.Hash)
	if err != nil {
		fmt.Println("Encoding error ", err.Error())

		return
	}
	fmt.Fprintf(w, "success")

}
