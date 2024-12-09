package services

import (
	"fmt"
	"github.com/sajjadvaezi/face-recognition/db"
	"github.com/sajjadvaezi/face-recognition/models"
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
