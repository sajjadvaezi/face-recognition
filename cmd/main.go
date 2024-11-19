package main

import (
	"encoding/json"
	"fmt"
	"github.com/sajjadvaezi/face-recognition/db"
	"github.com/sajjadvaezi/face-recognition/internal"
	"os/exec"

	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func callPythonFunction(command, faceData string) (string, error) {
	// Command to run Python script
	cmd := exec.Command("python3", "./face/main.py", command, faceData)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Parse JSON output
	var result map[string]string
	json.Unmarshal(output, &result)
	return result["result"], nil
}

func main() {
	db.InitSQLite()

	mux := internal.SetupRouter()

	//takePhotoRes, err := callPythonFunction("photo", "")
	//if err != nil {
	//	fmt.Println("Error registering face:", err)
	//	return
	//}

	registerMsg, err := callPythonFunction("register", "photo.jpg")
	if err != nil {
		fmt.Println("Error registering face:", err)
		return
	}

	//fmt.Println("takePhotoRes:", takePhotoRes)
	fmt.Println("registerMsg: ", registerMsg)

	err = http.ListenAndServe("localhost:8090", mux)
	if err != nil {
		return
	}

}
