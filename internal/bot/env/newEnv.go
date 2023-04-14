package env

import (
	"errors"
	"log"
	"strconv"

	"github.com/joho/godotenv"
)

// file with settings for enviroment
const envLoc = ".env"

var evn map[string]string

func Load() error {
	var err error
	if evn, err = godotenv.Read(envLoc); err != nil {
		return err
	}
	return nil
}

// return token
func LoadTGToken() (string, error) {
	token, ok := evn["TG_KEY"]
	if !ok {
		err := errors.New("telegram token not found in .evn")
		return "", err
	}
	return token, nil
}

// return env=id
func LoadAdminsID() map[string]int64 {
	adminID := make(map[string]int64)
	for admin, id := range evn {
		switch admin {
		case "ADMIN_ID":
			adminId, err := strconv.ParseInt(id, 0, 64)
			if err != nil {
				log.Println("error parse int: ", err)
			}
			adminID["ADMIN_ID"] = adminId

		case "MINTY_ID":

			mintyID, err := strconv.ParseInt(id, 0, 64)
			if err != nil {
				log.Println("error parse int: ", err)
			}
			adminID["MINTY_ID"] = mintyID

		case "OK_ID":

			okID, err := strconv.ParseInt(id, 0, 64)
			if err != nil {
				log.Println("error parse int: ", err)
			}
			adminID["OK_ID"] = okID
		case "MURS_ID":

			mursID, err := strconv.ParseInt(id, 0, 64)
			if err != nil {
				log.Println("error parse int: ", err)
			}
			adminID["MURS_ID"] = mursID

		}
	}

	return adminID
}

// return env=aiKey
func LoadAdminsAiKey() map[string]string {
	adminAiKey := make(map[string]string)
	for admin, aiKey := range evn {
		switch admin {
		case "ADMIN_KEY":
			adminAiKey["ADMIN_KEY"] = aiKey
		case "MINTY_KEY":
			adminAiKey["MINTY_KEY"] = aiKey
		case "OK_KEY":
			adminAiKey["OK_KEY"] = aiKey
		case "MURS_KEY":
			adminAiKey["MURS_KEY"] = aiKey
		}
	}

	return adminAiKey
}
