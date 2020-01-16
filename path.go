package fs

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Path represents a valid filesystem path.
type Path string

// WalkType determines the type of walking to be done by the Walk routine
type WalkType uint

const (
	// WalkBoth will consider both directories and files in the walking routine
	WalkBoth WalkType = iota

	// WalkFiles should consider only files in the walking routine
	WalkFiles

	// WalkDirs should consider only directories in the walking routine
	WalkDirs
)

const (
	defaultFileMode os.FileMode = 0644 // rw-r--r--
	defaultDirMode  os.FileMode = 0755 // rwxr-xr-x

	openFileFlag   int = os.O_RDONLY
	createFileFlag int = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	appendFileFlag int = os.O_WRONLY | os.O_CREATE | os.O_APPEND
)

// Info returns a info of a path
func (p Path) Info() os.FileInfo {
	if info, err := os.Stat(p.String()); err == nil {
		return info
	}
	return nil
}

// Exists returns true if the given path exists
func (p Path) Exists() bool {
	return p.Info() != nil
}

// FileExists returns true if the given path exists and is a regular file.
func (p Path) FileExists() bool {
	if info := p.Info(); info != nil {
		return info.Mode().IsRegular()
	}

	return false
}

// DirExists returns true if the given path exists and is a directory.
func (p Path) DirExists() bool {
	if info := p.Info(); info != nil {
		return info.IsDir()
	}

	return false
}

// Open opens the file specified by path for reading.
func (p Path) Open() (*os.File, error) {
	if !p.FileExists() {
		return nil, ErrFileDoesNotExist
	}

	return open(p, openFileFlag, defaultFileMode)
}

// Create open the specified file for writing, creating a new file if necessary.
// If the file already exists, it is overridden.
func (p Path) Create() (*os.File, error) {
	return open(p, createFileFlag, defaultFileMode)
}

// Append works like create, but instead of discarding the content of an existing file,
// it just appends the new data at the end of the file.
func (p Path) Append() (*os.File, error) {
	file, err := open(p, appendFileFlag, defaultFileMode)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// RemoveAll files or directory in the given path
func (p Path) RemoveAll() {
	os.RemoveAll(p.String())
}

// MkdirAll creates all directories that doesn't exists
func (p Path) MkdirAll() error {
	if err := os.MkdirAll(p.String(), defaultDirMode); err != nil {
		return err
	}
	return nil
}

// ReadAll returns all the content of a file
func (p Path) ReadAll() ([]byte, error) {
	if !p.FileExists() {
		return nil, ErrFileDoesNotExist
	}
	return ioutil.ReadFile(p.String())
}

// ReadDir reads the directory named by dirname and returns
// a list of directory entries sorted by filename.
func (p Path) ReadDir() ([]Path, error) {
	if !p.DirExists() {
		return nil, ErrDirDoesNotExist
	}
	osFiles, err := ioutil.ReadDir(p.String())
	paths := make([]Path, len(osFiles))

	if err != nil {
		return paths, err
	}

	for i := range osFiles {
		paths[i] = Path(osFiles[i].Name())
	}

	return paths, nil
}

// CopyTo copies the the data at the location represented by the receiver to
// a given destination. If the receiver is a directory, a recursive copy of
// its contents is made.
func (p Path) CopyTo(dest Path) error {
	return copy(p, dest)
}

// Join join the current path with the specified string value
// and returns a new path
func (p Path) Join(other string) Path {
	return Path(filepath.Join(string(p), other))
}

// JoinP join the current path with the specified path
// and returns a new path
func (p Path) JoinP(other Path) Path {
	return p.Join(string(other))
}

// String converts a path to its string representation
func (p Path) String() string {
	return string(p)
}

// Empty returns true when the path is empty
func (p Path) Empty() bool {
	return strings.TrimSpace(p.String()) == ""
}

// Basename returns the name of the last element of the path
func (p Path) Basename() string {
	return filepath.Base(p.String())
}

// Ext returns the extension of the path, including the "." character.
func (p Path) Ext() string {
	return filepath.Ext(p.String())
}

// Parent returns the parent directory of the current path.
func (p Path) Parent() Path {
	return Path(filepath.Dir(p.String()))
}

// Clean returns the shortest path name equivalent to path
func (p Path) Clean() Path {
	return Path(filepath.Clean(p.String()))
}

// Walk walks on every item (configurable by the 'walkType') parameter and call
// the walker function.
func (p Path) Walk(walkType WalkType, walker func(path Path, isDirectory bool) error) error {
	if !p.DirExists() {
		return ErrDirDoesNotExist
	}

	return filepath.Walk(p.String(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == p.String() {
			return nil
		}

		skipper := filepath.SkipDir
		if Dirname(path) == p.String() {
			skipper = nil
		}

		if info.IsDir() && walkType == WalkFiles {
			return skipper
		}

		if !info.IsDir() && walkType == WalkDirs {
			return skipper
		}

		return walker(Path(path), info.IsDir())
	})
}

// Abs returns an absolute representation of path, when possible
func (p Path) Abs() Path {
	if path, err := filepath.Abs(p.String()); err == nil {
		return Path(path)
	}
	return p
}

func open(p Path, flag int, mode os.FileMode) (*os.File, error) {
	if p.Empty() {
		return nil, ErrPathIsEmpty
	}

	if p.DirExists() {
		return nil, ErrPathIsDirectory
	}

	file, err := os.OpenFile(p.String(), flag, mode)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		if err = p.Clean().Parent().MkdirAll(); err != nil {
			return nil, err
		}
		return os.OpenFile(p.String(), flag, mode)
	}

	return file, nil
}

// copy copy one path to another
func copy(src, dest Path) error {
	if !src.Exists() {
		return ErrNotFound
	}

	if src.FileExists() {
		if dest.DirExists() {
			return copyFiles(src, dest.Join(src.Basename()))
		}
		return copyFiles(src, dest)
	}

	if dest.FileExists() {
		return ErrPathIsDirectoryDestFile
	}

	return copyDirs(src, dest)
}

// copyDirs copy one dir to another
func copyDirs(src, dest Path) error {
	if !dest.DirExists() {
		if err := dest.MkdirAll(); err != nil {
			return err
		}
	}
	return src.Walk(WalkBoth, func(path Path, isDirectory bool) error {
		newDest := Path(strings.Replace(path.String(), src.String(), dest.String(), 1))

		if isDirectory {
			if err := newDest.MkdirAll(); err != nil {
				return err
			}
		} else {
			if err := copyFiles(path, newDest); err != nil {
				return err
			}
		}

		return nil
	})
}

// copyFiles copy one file to another
func copyFiles(src, dest Path) error {
	info := src.Info()
	if info == nil {
		return ErrFileDoesNotExist
	}

	srcFile, err := open(src, openFileFlag, 0400) //r--------
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := open(dest, createFileFlag, info.Mode())
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

// Count how many files are on some path 'p'
func (p Path) Count(walkType WalkType) (count uint64) {
	if !p.DirExists() {
		return
	}

	p.Walk(walkType, func(path Path, isDirectory bool) error { // nolint: errcheck
		count++
		return nil
	})
	return
}
