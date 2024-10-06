package yt_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVideoInfo(t *testing.T) {
	// Create a temporary file to use as the infoPath
	file, err := os.CreateTemp("", "info.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	// Write some test data to the file
	testData := yt.VideoInfo{
		ID:      "test-id",
		Channel: "test-channel",
	}
	err = json.NewEncoder(file).Encode(testData)
	if err != nil {
		t.Fatal(err)
	}

	ytService := yt.NewYTVideoProcessor()

	// Call the function being tested
	info, err := ytService.GetVideoInfo(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Assert that the returned data matches the test data
	assert.Equal(t, testData, info)
}
