// Package fs implements utility routines for path creation and other common
// filesystem operations
package fs

import (
	"os"
	"path/filepath"
)

// MkdirAll creates all directories that doesn't exists
func MkdirAll(path string) error {
	return Path(path).MkdirAll()
}

// Info returns a info of a string path
func Info(path string) os.FileInfo {
	return Path(path).Info()
}

// FileExists returns true if the given path exists and is a regular file.
func FileExists(path string) bool {
	return Path(path).FileExists()
}

// DirExists returns true if the given path exists and is a directory.
func DirExists(path string) bool {
	return Path(path).DirExists()
}

// Exists returns true if the given path exists
func Exists(path string) bool {
	return Path(path).Exists()
}

// Open opens an existing file for reading.
func Open(path string) (*os.File, error) {
	return Path(path).Open()
}

// Create opens the specified file for writing, creating a new file if necessary.
// If the file already exists, it is overridden.
func Create(path string) (*os.File, error) {
	return Path(path).Create()
}

// Append works like create, but instead of discarding the content of an existing file,
// it just appends the new data at the end of the file.
func Append(path string) (*os.File, error) {
	return Path(path).Append()
}

// RemoveAll files or directory in the given path
func RemoveAll(path string) {
	Path(path).RemoveAll()
}

// ReadAll returns all the content of a file
func ReadAll(path string) ([]byte, error) {
	return Path(path).ReadAll()
}

// Abs returns an absolute representation of path.
// If the path is not absolute it will be joined with the current
// working directory to turn it into an absolute path. The absolute
// path name for a given file is not guaranteed to be unique.
// Abs calls Clean on the result.
func Abs(path string) string {
	path, _ = filepath.Abs(path)
	return path
}

// Basename returns the name of the last element of the path
func Basename(path string) string {
	return Path(path).Basename()
}

// Clean returns the shortest path name equivalent to path
func Clean(path string) string {
	return Path(path).Clean().String()
}

// Dirname returns all but the last element of path, typically the path's directory
func Dirname(path string) string {
	return Path(path).Parent().String()
}
