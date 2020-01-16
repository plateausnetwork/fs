package fs_test

import (
	"bytes"
	"os"
	"reflect"
	"syscall"
	"testing"

	"github.com/rhizomplatform/fs"
)

func TestMkdir(t *testing.T) {
	WithTempDir(func(root string) {
		basePath := fs.Path(root + "testdata/another/dir/")
		tests := []struct {
			source   fs.Path
			expected error
		}{
			{
				source:   basePath.Join("./../test.txt"),
				expected: nil,
			},
			{
				source:   basePath.Join(".testsdsadaskdljsa.tar.gz"),
				expected: nil,
			},
			{
				source:   basePath,
				expected: nil,
			},
			{
				source: "",
				expected: &os.PathError{
					Op:   "mkdir",
					Path: "",
					Err:  syscall.ENOENT,
				},
			},
		}

		for i, test := range tests {
			if err := test.source.MkdirAll(); !reflect.DeepEqual(err, test.expected) {
				t.Errorf("Case %d, error testing MkdirAll: %v", i, err)
			}
			if err := fs.MkdirAll(test.source.String()); !reflect.DeepEqual(err, test.expected) {
				t.Errorf("Case %d, error testing MkdirAll: %v", i, err)
			}
		}
	})
}

func TestInfo(t *testing.T) {
	WithTempDir(func(root string) {
		basePath := fs.Path(root + "testdata/another/dir/")

		tests := []struct {
			path     fs.Path
			expected bool
		}{
			{
				path:     basePath.Join("test.txt"),
				expected: true,
			},
			{
				path:     basePath,
				expected: true,
			},
			{
				path:     fs.Path(root + "/dsadijsaio"),
				expected: false,
			},
		}

		file, err := tests[0].path.Create()
		if err != nil {
			t.Errorf("error creating file: '%v'", err)
		}
		defer file.Close()

		for i, test := range tests {
			filename := ""
			info := fs.Info(test.path.String())
			if info != nil {
				filename = info.Name()
			}
			received := filename == test.path.Basename()
			if received != test.expected {
				t.Errorf("Case %d, error testing Infos: expected '%v', received '%v'", i, test.expected, received)
			}

			info2 := test.path.Info()
			if info2 != nil {
				filename = info.Name()
			}
			received = filename == test.path.Basename()
			if received != test.expected {
				t.Errorf("Case %d, error testing Infos: expected '%v', received '%v'", i, test.expected, received)
			}
		}
	})
}

func TestExists(t *testing.T) {
	tests := []struct {
		dir  fs.Path
		file fs.Path
	}{
		{dir: "", file: "foo.txt"},
		{dir: "bar", file: "foo.txt"},
		{dir: "bar/baz", file: "foo.txt"},
	}

	WithTempDir(func(root string) {
		for i, test := range tests {
			dir := fs.Path(root).JoinP(test.dir)
			file := dir.JoinP(test.file)

			if dir.String() != root && dir.DirExists() {
				t.Errorf("Case %d, directory should NOT exist", i)
				continue
			}

			if file.FileExists() {
				t.Errorf("Case %d, file should NOT exist", i)
				continue
			}

			fileAp, err := file.Create()
			if err != nil {
				t.Errorf("Case %d, failed to create file: %v", i, err)
				continue
			}
			defer fileAp.Close()

			if !dir.DirExists() {
				t.Errorf("Case %d, directory should exist", i)
			}
			if !fs.DirExists(dir.String()) {
				t.Errorf("Case %d, directory should exist", i)
			}

			if !fs.Exists(dir.String()) && !fs.Exists(file.String()) {
				t.Errorf("Case %d, should exist", i)
			}
			if !dir.Exists() && !file.Exists() {
				t.Errorf("Case %d, should exist", i)
			}

			if !file.FileExists() {
				t.Errorf("Case %d, file should exist", i)
			}
			if !fs.FileExists(file.String()) {
				t.Errorf("Case %d, file should exist", i)
			}
		}
	})
}

func TestOpen(t *testing.T) {
	WithTempDir(func(dir string) {
		root := fs.Path(dir)

		f, err := root.Join("does").Join("exists.txt").Create()
		if err != nil {
			t.Error(err)
			return
		}
		f.Close()

		tests := []struct {
			path     fs.Path
			expected error
		}{
			{path: root.Join("does/exists.txt"), expected: nil},
			{path: root.Join("does"), expected: fs.ErrFileDoesNotExist},
			{path: root.Join("does/not/exists.txt"), expected: fs.ErrFileDoesNotExist},
			{path: root.Join("does/not"), expected: fs.ErrFileDoesNotExist},
			{path: root.Join("does/not/../exists.txt"), expected: nil},
			{path: fs.Path(""), expected: fs.ErrFileDoesNotExist},
		}

		for i, test := range tests {
			path := test.path

			if _, err := path.Open(); err != test.expected {
				t.Errorf("Case %d, error testing open: expected '%v', received '%v'", i, test.expected, err)
			}

			if _, err := fs.Open(path.String()); err != test.expected {
				t.Errorf("Case %d, error testing open: expected '%v', received '%v'", i, test.expected, err)
			}
		}
	})
}

func TestWrite(t *testing.T) {
	WithTempDir(func(root string) {
		src := fs.Path(root).Join("a.txt")
		dest := fs.Path(root).Join("b.txt")

		srcFile, err := fs.Create(src.String())
		if err != nil {
			t.Errorf("Create text on source error %v", err)
		}

		if _, err := srcFile.Write([]byte("this is a test")); err != nil {
			t.Errorf("Write text on source error %v", err)
		}

		if err := src.CopyTo(dest); err != nil {
			t.Errorf("Copy dir to dir error %v", err)
		}

		a, err := src.ReadAll()
		if err != nil {
			t.Errorf("Error reading from source: %v", err)
		}

		b, err := dest.ReadAll()
		if err != nil {
			t.Errorf("Error reading from destination: %v", err)
		}

		if !bytes.Equal(a, b) {
			t.Errorf("Content of the copied fles does not match")
		}
	})
}

func TestReadAllError(t *testing.T) {
	basePath := fs.Path("testdata/another/dir/")

	tests := []struct {
		source   fs.Path
		expected error
	}{
		{
			source:   basePath.Join("test.txt"),
			expected: fs.ErrFileDoesNotExist,
		},
		{
			source:   basePath.Join("testsdsadaskdljsa"),
			expected: fs.ErrFileDoesNotExist,
		},
	}

	for i, test := range tests {
		if _, err := test.source.ReadAll(); err != test.expected {
			t.Errorf("Case %d, error testing ReadAll: %v", i, err)
		}
		if _, err := fs.ReadAll(test.source.String()); err != test.expected {
			t.Errorf("Case %d, error testing ReadAll: %v", i, err)
		}
	}
}

func TestReadDir(t *testing.T) {
	WithTempDir(func(root string) {
		basePath := fs.Path(root + "testdata/another/dir/")

		tests := []struct {
			source         fs.Path
			expectedParent error
			expected       error
		}{
			{
				source:         basePath.Join("testsdsadaskdljsa"),
				expectedParent: nil,
				expected:       fs.ErrDirDoesNotExist,
			},
			{
				source:         basePath.Join("testsdsadaskdljsa.c"),
				expectedParent: nil,
				expected:       fs.ErrDirDoesNotExist,
			},
			{
				source:         basePath.Join("anotheragain/.c"),
				expectedParent: nil,
				expected:       fs.ErrDirDoesNotExist,
			},
		}

		for i, test := range tests {
			file, err := test.source.Append()
			if err != nil {
				t.Errorf("Case %d, error apending on source: '%v'", i, err)
			}
			defer file.Close()

			if _, err := file.Write([]byte("a")); err != nil {
				t.Errorf("Case %d, error writing to file: %v", i, err)
			}

			if _, err := test.source.Parent().ReadDir(); err != test.expectedParent {
				t.Errorf("Case %d, error testing read dir: expected '%v', received '%v'", i, test.expected, err)
			}
			if _, err := test.source.ReadDir(); err != test.expected {
				t.Errorf("Case %d, error testing read dir: expected '%v', received '%v'", i, test.expected, err)
			}
		}
	})
}
