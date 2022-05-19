package util

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/labstack/gommon/log"
)

func GetDirFolders(dirname string) ([]string, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Errorf("Failed to read dir %s", dirname)
		return nil, err
	}
	var filenames []string
	for _, file := range files {
		if file.IsDir() {
			filenames = append(filenames, file.Name())
		}
	}
	return filenames, nil
}

func GetFilesinDir(dirname string) ([]string, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Errorf("Failed to read dir %s", dirname)
		return nil, err
	}
	var filenames []string
	for _, file := range files {
		filenames = append(filenames, file.Name())
	}
	return filenames, nil
}

func StringDateToDateTime(date string) time.Time {
	hypenDate := date[:4] + "-" + date[4:6] + "-" + date[6:]
	myDate, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT15:04:05Z", hypenDate))
	if err != nil {
		log.Errorf("Failed to parse date: %s", err)
		return time.Now()
	}
	return myDate
}

// Check if string is in array of strings
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
