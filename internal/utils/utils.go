package utils

import (
	"io/fs"
	"os"
)

// GetFilesInDirectory reads the files in the specified directory and returns a slice of fs.DirEntry.
// It excludes directories and only includes files in the returned list.
func GetFilesInDirectory(path string) ([]fs.DirEntry, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var paths []fs.DirEntry
	for _, file := range files {
		if !file.IsDir() {
			paths = append(paths, file)
		}
	}

	return paths, nil
}

func StringInSlice(slice []string, str string) bool {
	stringMap := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		stringMap[s] = struct{}{}
	}
	_, found := stringMap[str]
	return found
}
