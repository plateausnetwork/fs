package fs

import (
	"errors"
)

// ErrDirDoesNotExist is a error indicating that a given directory does not exists.
var ErrDirDoesNotExist = errors.New("directory does not exist")

// ErrFileDoesNotExist is a error indicating that a given file does not exists.
var ErrFileDoesNotExist = errors.New("file does not exist")

// ErrNotFound is a error indicating that a given path does not exists.
var ErrNotFound = errors.New("Path not found")

// ErrPathIsEmpty is a error indicating that a given path is empty like ''.
var ErrPathIsEmpty = errors.New("Path is empty")

// ErrPathIsDirectory is a error indicating that a given path is a directory.
var ErrPathIsDirectory = errors.New("Path is a directory")

// ErrPathIsDirectoryDestFile is a error indicating that a given source path is a directory and the destination is a file.
var ErrPathIsDirectoryDestFile = errors.New("Path source is a directory and the destination file")

// ErrOneDirectoryOtherFile is a error indicating that a one is a directory and the other is a file.
var ErrOneDirectoryOtherFile = errors.New("One is a directory and the other is a file")

// ErrFilesNotEquals is a error indicating that the files source and destionation aren't equals.
var ErrFilesNotEquals = errors.New("Source and destination files aren't equals")
