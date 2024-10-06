package utils_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zibbp/eos/internal/utils"
)

func TestGetFilesInDirectory(t *testing.T) {

	tmp, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	// Create a temporary file in the temporary directory
	file, err := os.CreateTemp(tmp, "file.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	files, err := utils.GetFilesInDirectory(tmp)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(files), 1)

	// Remove the temporary file
	err = os.Remove(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Remove the temporary directory
	err = os.Remove(tmp)
	if err != nil {
		t.Fatal(err)
	}

}

func TestStringInSlice(t *testing.T) {
	slice := []string{"one", "two", "three"}
	str := "one"

	found := utils.StringInSlice(slice, str)

	assert.True(t, found, "Expected string to be found in the slice")
}
