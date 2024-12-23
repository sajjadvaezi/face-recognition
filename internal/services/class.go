package services

import (
	"fmt"
	"github.com/sajjadvaezi/face-recognition/db"
	"github.com/sajjadvaezi/face-recognition/models"
	"log/slog"
)

func AddClass(addClass models.AddClassRequest) error {
	if addClass.ClassName == "" {
		return fmt.Errorf("empty classname")
	}
	if addClass.UserNumber == "" {
		return fmt.Errorf("empty usernumber")
	}
	err := db.AddClass(addClass.ClassName, addClass.UserNumber)
	if err != nil {
		return fmt.Errorf("error adding to db, error: %s", err.Error())
	}
	return nil
}

func Attendance(request models.AttendanceRequest) error {
	user, err := RecognizeUser(request.Image)
	if err != nil {
		slog.Error(fmt.Sprintf("error recognizing user %s", err.Error()))
		return fmt.Errorf("could not recognize user")
	}

	_, err = db.Attendance(user.UserNumber, request.ClassName)
	if err != nil {
		slog.Error(fmt.Sprintf("couldn't not attend, error: %s", err.Error()))
		return fmt.Errorf("user %s could not attend the class %s ", user.UserNumber, request.ClassName)
	}
	slog.Info(fmt.Sprintf("success attendance of %s with user_number %s at class %s", user.Name, user.UserNumber, request.ClassName))
	return nil
}

func AttendedUsers(className string) (models.Users, error) {
	if className == "" {
		return nil, fmt.Errorf("empty class name")
	}
	users, err := db.AttendedUsers(className)
	if err != nil {
		slog.Error(fmt.Sprintf("error retriving users in class %s, error: %s", className, err.Error()))
		return nil, fmt.Errorf("error getting users")
	}
	return users, nil
}
