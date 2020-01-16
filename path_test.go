package fs_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/rhizomplatform/fs"
)

// TempDirHandler is a handler that receives the path of a temporary
// directory as parameter
type TempDirHandler func(string)

// WithTempDir runs the specified handler in a context with a
// temporary directory available
func WithTempDir(handler TempDirHandler) {
	if dir, err := ioutil.TempDir("", ""); err != nil {
		panic(err)
	} else {
		defer fs.RemoveAll(dir)
		handler(dir)
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		start     string
		fragments []string
		expected  string
	}{
		{
			start:     "/foo",
			fragments: []string{"bar", "baz"},
			expected:  "/foo/bar/baz",
		},
		{
			start:     "/foo",
			fragments: []string{"", "baz"},
			expected:  "/foo/baz",
		},
		{
			start:     "/foo",
			fragments: []string{"bar", ""},
			expected:  "/foo/bar",
		},
		{
			start:     "/foo",
			fragments: []string{"", ""},
			expected:  "/foo",
		},
		{
			start:     "",
			fragments: []string{"", ""},
			expected:  "",
		},
		{
			start:     "",
			fragments: []string{"bar", "baz"},
			expected:  "bar/baz",
		},
		{
			start:     "",
			fragments: []string{"/bar", "baz"},
			expected:  "/bar/baz",
		},
		{
			start:     "/foo/",
			fragments: []string{"/bar/", "/baz/"},
			expected:  "/foo/bar/baz",
		},
		{
			start:     "/foo//bar",
			fragments: []string{"", "/baz/"},
			expected:  "/foo/bar/baz",
		},
	}

	for i, test := range tests {
		ps := fs.Path(test.start)
		pp := fs.Path(test.start)

		for _, f := range test.fragments {
			ps = ps.Join(f)
			pp = pp.JoinP(fs.Path(f))
		}

		if actual := ps.String(); actual != test.expected {
			t.Errorf("Case %d, error joining strings: expected '%s', received '%s'", i, test.expected, actual)
		}

		if actual := pp.String(); actual != test.expected {
			t.Errorf("Case %d, error joining paths: expected '%s', received '%s'", i, test.expected, actual)
		}
	}
}

func TestJoinOrder(t *testing.T) {
	const expected = "foo/bar/baz"

	a := fs.Path("foo/")
	b := fs.Path("bar/")
	c := fs.Path("baz/")

	ab := a.JoinP(b)
	bc := b.JoinP(c)

	s1 := ab.JoinP(c).String()
	s2 := a.JoinP(bc).String()

	if s1 != s2 {
		t.Errorf("Joining order should not matter for subpaths; '%s' not equal to '%s'", s1, s2)
	}

	if s1 != expected {
		t.Errorf("Error joining paths, expected '%s', received '%s'", expected, s1)
	}
}

func TestEmpty(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{path: "/foo", expected: false},
		{path: "foo", expected: false},
		{path: "/", expected: false},
		{path: "//", expected: false},
		{path: "", expected: true},
		{path: " ", expected: true},
		{path: "       ", expected: true},
		{path: "   /   ", expected: false},
	}

	for i, test := range tests {
		path := fs.Path(test.path)

		if path.Empty() != test.expected {
			t.Errorf("Case %d, error testing empty path: expected '%v', received '%v'", i, test.expected, path.Empty())
		}
	}
}

func TestBasename(t *testing.T) {
	tests := []struct {
		path     fs.Path
		expected fs.Path
	}{
		{path: "/foo", expected: "foo"},
		{path: "foo.txt", expected: "foo.txt"},
		{path: "/tao/euhaus/teste.go", expected: "teste.go"},
	}

	for i, test := range tests {
		received := test.path.Basename()
		if received != test.expected.String() {
			t.Errorf("Case %d, error testing basename: expected '%v', received '%v'", i, test.expected, received)
		}
	}
}

func TestExt(t *testing.T) {
	tests := []struct {
		path     fs.Path
		expected string
	}{
		{
			path:     "/rafilx/../rafilxtest.txt",
			expected: ".txt",
		},
		{
			path:     "/rafilx/tstgo.go",
			expected: ".go",
		},
		{
			path:     "/testdata/hu3.tar.gz",
			expected: ".gz",
		},
		{
			path:     "/kk/hue/../../testsasasa",
			expected: "",
		},
		{
			path:     "",
			expected: "",
		},
	}

	for i, test := range tests {
		if received := test.path.Ext(); received != test.expected {
			t.Errorf("Case %d, error testing Ext: expected '%v', received '%v'", i, test.expected, received)
		}
	}
}

func TestParent(t *testing.T) {
	tests := []struct {
		path     fs.Path
		expected fs.Path
	}{
		{
			path:     "/rafilx/rafilxtest.txt",
			expected: "/rafilx",
		},
		{
			path:     "/kk/hue/testsasasa",
			expected: "/kk/hue",
		},
		{
			path:     "",
			expected: ".",
		},
	}

	for i, test := range tests {
		if received := test.path.Parent(); received != test.expected {
			t.Errorf("Case %d, error testing Parent: expected '%v', received '%v'", i, test.expected, received)
		}
	}
}

func TestDirname(t *testing.T) {
	tests := []struct {
		path     fs.Path
		expected string
	}{
		{
			path:     "/rafilx/rafilxtest.txt",
			expected: "/rafilx",
		},
		{
			path:     "/kk/hue/testsasasa",
			expected: "/kk/hue",
		},
		{
			path:     "",
			expected: ".",
		},
	}

	for i, test := range tests {
		if received := fs.Dirname(test.path.String()); received != test.expected {
			t.Errorf("Case %d, error testing Dirname: expected '%v', received '%v'", i, test.expected, received)
		}
	}
}

func TestClean(t *testing.T) {
	tests := []struct {
		path     fs.Path
		expected fs.Path
	}{
		{path: "/foo/../a.txt", expected: "/a.txt"},
		{path: "/foo/bar/../a.txt", expected: "/foo/a.txt"},
		{path: "/foo/bar/../../a.txt", expected: "/a.txt"},
		{path: "", expected: "."},
	}

	for i, test := range tests {
		if received := test.path.Clean(); received != test.expected {
			t.Errorf("Case %d, error testing clean: expected '%v', received '%v'", i, test.expected, received)
		}
		if received := fs.Clean(test.path.String()); received != test.expected.String() {
			t.Errorf("Case %d, error testing clean: expected '%v', received '%v'", i, test.expected, received)
		}
	}
}

func TestWalkErrDirExists(t *testing.T) {
	WithTempDir(func(root string) {
		basePath := fs.Path(root)

		tests := []struct {
			path     fs.Path
			expected error
		}{
			{path: basePath.Join("foo"), expected: fs.ErrDirDoesNotExist},
			{path: basePath.Join("bar"), expected: fs.ErrDirDoesNotExist},
		}

		for i, test := range tests {
			received := test.path.Walk(fs.WalkBoth, func(path fs.Path, isDirectory bool) error { return nil })
			if received != test.expected {
				t.Errorf("Case %d, error testing walk Directory Doesn't Exists: expected '%v' received '%v' ", i, test.expected, received)
			}
		}
	})
}

func TestWalkErrFunc(t *testing.T) {
	WithTempDir(func(root string) {
		basePath := fs.Path(root)
		p1 := basePath.Join("test")
		p2 := basePath.Join("testsdsadaskdljsa")
		p3 := basePath.Join("mytest")

		tests := []struct {
			path     fs.Path
			expected error
		}{
			{path: p1, expected: fs.ErrDirDoesNotExist},
			{path: p2, expected: fs.ErrFileDoesNotExist},
			{path: p3, expected: fs.ErrPathIsEmpty},
			{path: p2, expected: fs.ErrNotFound},
		}

		for i, test := range tests {
			if err := test.path.MkdirAll(); err != nil {
				t.Errorf("Case %d, error testing walk Creating dir: '%v' ", i, err)
			}
			err := test.expected
			received := basePath.Walk(fs.WalkBoth, func(path fs.Path, isDirectory bool) error {
				return err
			})
			if received != test.expected {
				t.Errorf("Case %d, error testing walk Err: expected '%v' received '%v' ", i, test.expected, received)
			}
		}
	})
}

func TestWalkFile(t *testing.T) {
	WithTempDir(func(dir string) {
		root := fs.Path(dir)

		if err := createLogDir(root); err != nil {
			t.Errorf("walk error creating tree %v", err)
		}

		// expected files
		found := make(map[string]bool)
		found[root.Join("log").Join("a.log").String()] = false
		found[root.Join("log").Join("b.log").String()] = false
		found[root.Join("log").Join("c.log").String()] = false

		// mark the ones that we found
		if err := root.Walk(fs.WalkFiles, func(path fs.Path, isDirectory bool) error {
			key := path.String()

			if duplicated, ok := found[key]; duplicated {
				t.Errorf("Duplicated path: '%s'", key)
			} else if !ok {
				t.Errorf("Path should not exist: '%s'", key)
			}

			found[key] = true
			return nil
		}); err != nil {
			t.Errorf("Error walking into files: %v", err)
		}

		// check if some of them were not found
		for k, ok := range found {
			if !ok {
				t.Errorf("Path '%s' was not found", k)
			}
		}
	})
}

func createLogDir(root fs.Path) error {
	paths := []struct {
		path    fs.Path
		content string
	}{
		{path: root.Join("dir")},
		{path: root.Join("log").Join("a.log"), content: "a.log"},
		{path: root.Join("log").Join("b.log"), content: "b.log"},
		{path: root.Join("log").Join("c.log"), content: "c.log"},
	}

	for _, p := range paths {
		if p.content == "" {
			if err := p.path.MkdirAll(); err != nil {
				return err
			}

			continue
		}

		file, err := p.path.Create()
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := file.Write([]byte(p.content)); err != nil {
			return err
		}
	}

	return nil
}

func TestWalkDirs(t *testing.T) {
	WithTempDir(func(dir string) {
		root := fs.Path(dir)

		// create test content to walk
		if err := createLogDir(root); err != nil {
			t.Errorf("Error creating test dir: %v", err)
			return
		}

		// expected directories
		found := make(map[string]bool)
		found[root.Join("dir").String()] = false
		found[root.Join("log").String()] = false

		// mark the ones that we found
		if err := root.Walk(fs.WalkDirs, func(path fs.Path, isDirectory bool) error {
			key := path.String()

			if duplicated, ok := found[key]; duplicated {
				t.Errorf("Duplicated path: '%s'", key)
			} else if !ok {
				t.Errorf("Path should not exist: '%s'", key)
			}

			found[key] = true
			return nil
		}); err != nil {
			t.Errorf("Error walking into directory: %v", err)
		}

		// check if some of them were not found
		for k, ok := range found {
			if !ok {
				t.Errorf("Path '%s' was not found", k)
			}
		}
	})
}

func TestAbs(t *testing.T) {
	tests := []struct {
		path     fs.Path
		expected fs.Path
	}{
		{path: "/foo/../a.txt", expected: "/a.txt"},
		{path: "/foo/bar/../../a.txt", expected: "/a.txt"},
		{path: "/foo/bar/../a.txt", expected: "/foo/a.txt"},
		{path: "/foo/bar/../../../", expected: "/"},
	}

	for i, test := range tests {
		received := test.path.Abs()
		if received != test.expected {
			t.Errorf("Case %d, error testing abs: expected '%v', received '%v'", i, test.expected, received)
		}
		receivedString := fs.Abs(test.path.String())
		if receivedString != test.expected.String() {
			t.Errorf("Case %d, error testing abs: expected '%v', received '%v'", i, test.expected, received)
		}
	}
}

func TestCopyToPathErrors(t *testing.T) {
	WithTempDir(func(root string) {
		filePath := fs.Path(root).Join("foo/a.txt")

		file, err := filePath.Create()
		if err != nil {
			t.Errorf("Error creating a file/directory to test: %v", err)
		}
		file.Close()

		tests := []struct {
			source      fs.Path
			destination fs.Path
			expected    error
		}{
			{
				source:      fs.Path(root).Join("testnull.txt"),
				destination: fs.Path(root).Join("testdata/testnull.txt"),
				expected:    fs.ErrNotFound,
			},
			{
				source:      filePath.Parent(),
				destination: filePath,
				expected:    fs.ErrPathIsDirectoryDestFile,
			},
		}

		for i, test := range tests {
			if err := test.source.CopyTo(test.destination); err != test.expected {
				t.Errorf("Case %d, error testing copy: expected '%v', received '%v'", i, test.expected, err)
			}

		}
	})
}

func TestCopyToPathFiles(t *testing.T) {
	WithTempDir(func(root string) {
		basePath := fs.Path(root)

		tests := []struct {
			source      fs.Path
			content     []byte
			destination fs.Path
			expected    error
		}{
			{
				source:      basePath.Join("first.txt"),
				content:     []byte(uuid.New().String()),
				destination: basePath.Join("copydir/first.txt"),
				expected:    nil,
			},
			{
				source:      basePath.Join("second.txt"),
				content:     []byte(uuid.New().String()),
				destination: basePath.Join("/copydir/second.txt"),
				expected:    nil,
			},
		}

		for i, test := range tests {
			file, err := test.source.Create()
			if err != nil {
				t.Errorf("Case %d, error testing copy: writing a source file: '%v'", i, err)
			}
			defer file.Close()

			if _, err := file.Write(test.content); err != nil {
				t.Errorf("Case %d, error writing to file: %v", i, err)
			}

			if err := test.source.CopyTo(test.destination); err != test.expected {
				t.Errorf("Case %d, error testing copy: expected '%v', received '%v'", i, test.expected, err)
			}

			src, err := test.source.ReadAll()
			if err != nil {
				t.Errorf("Case %d, error reading from source: %v", i, err)
			}

			dst, err := test.destination.ReadAll()
			if err != nil {
				t.Errorf("case %d, error reading from destination: %v", i, err)
			}

			if !bytes.Equal(src, dst) {
				t.Errorf("Case %d, error testing files source and destination aren't equals", i)
			}
		}
	})
}

func TestCopyToPathFileToDir(t *testing.T) {
	WithTempDir(func(root string) {
		basePath := fs.Path(root)

		tests := []struct {
			source      fs.Path
			content     []byte
			destination fs.Path
		}{
			{
				source:      basePath.Join("a.txt"),
				content:     []byte(uuid.New().String()),
				destination: basePath.Join("/copydir/"),
			},
			{
				source:      basePath.Join("b.txt"),
				content:     []byte(uuid.New().String()),
				destination: basePath.Join("/anotheragain/"),
			},
		}

		for i, test := range tests {
			fileAp, err := test.source.Create()
			if err != nil {
				t.Errorf("Case %d, error testing writing a source file: '%v'", i, err)
				continue
			}
			defer fileAp.Close()

			if _, err := fileAp.Write(test.content); err != nil {
				t.Errorf("Case %d, error writing to file: %v", i, err)
			}

			if err := test.destination.MkdirAll(); err != nil {
				t.Errorf("Error creating directory: %v", err)
				continue
			}

			if err := test.source.CopyTo(test.destination); err != nil {
				t.Errorf("Case %d, error testing copy: %v", i, err)
			}

			src, err := test.source.ReadAll()
			if err != nil {
				t.Errorf("Case %d, error reading from source: %v", i, err)
			}

			dst, err := test.destination.Join(test.source.Basename()).ReadAll()
			if err != nil {
				t.Errorf("Case %d, error reading from destination: %v", i, err)
			}

			if !bytes.Equal(src, dst) {
				t.Errorf("Case %d, error testing files source and destination aren't equals", i)
			}
		}
	})
}

func TestCopyToPathDirToDir(t *testing.T) {
	WithTempDir(func(root string) {
		src := fs.Path(root).Join("src")
		dst := fs.Path(root).Join("dst")

		// create the source content
		if err := createTreeCopyDirToDir(src.String()); err != nil {
			t.Errorf("Error creating tree %v", err)
		}

		// copy all to destination
		if err := src.CopyTo(dst); err != nil {
			t.Errorf("Copy dir to dir error %v", err)
		}

		// mount a list with everything that should exist and the expected content
		data := make(map[string]string)
		if err := src.Walk(fs.WalkBoth, func(path fs.Path, isDirectory bool) error {
			destination := strings.Replace(path.String(), "src", "dst", 1)

			if isDirectory {
				data[destination] = ""
				return nil
			}

			content, err := path.ReadAll()
			if err != nil {
				return err
			}

			data[destination] = string(content)
			return nil
		}); err != nil {
			t.Errorf("Failed to list source content: %v", err)
		}

		// walk into the copy and see if everything matches
		if err := dst.Walk(fs.WalkBoth, func(path fs.Path, isDirectory bool) error {
			key := path.String()

			expected, ok := data[key]
			if !ok {
				return fmt.Errorf("destination '%v' should not exist", path)
			}

			if isDirectory {
				delete(data, key)
				return nil
			}

			b, err := path.ReadAll()
			if err != nil {
				return err
			}

			actual := string(b)

			if expected != actual {
				return fmt.Errorf("content mismatch on '%v'. Expected: '%s', received '%s'", path, expected, actual)
			}

			delete(data, key)
			return nil
		}); err != nil {
			t.Error(err)
		}

		if len(data) > 0 {
			t.Errorf("One or more paths were not created: %v", data)
		}
	})
}

func createTreeCopyDirToDir(root string) error {
	// to test the copy recursive between 2 directories it's necessary to create a tree directory
	/*
		├── testdir
		│   ├── dir1
		│   │   └── text.txt
		│   ├── another
		│   │   └── txt.go
		│   ├── dir2
		│   │   └── a.txt
		│   ├── empty
		│   ├── dir
		│	│	├── log
		│	│	│   ├── a.c
		│	│	│   ├── b.c
		│	│	│   ├── e.c
		│	│	│   ├── d.c
		------------------------
		├── dirtopastetestdir
		│   ├── dir1
		│   │   └── text.txt
		.....
	*/
	basePath := fs.Path(root)
	log := basePath.Join("dir").Join("log")

	content := []struct {
		path fs.Path
		file bool
	}{
		{path: basePath.Join("dir1").Join("text.txt"), file: true},
		{path: basePath.Join("another").Join("txt.go"), file: true},
		{path: basePath.Join("dir1").Join("text.txt"), file: true},
		{path: log.Join("a.c"), file: true},
		{path: log.Join("b.c"), file: true},
		{path: log.Join("c.c"), file: true},
		{path: basePath.Join("empty"), file: false},
	}

	for _, c := range content {
		if c.file {
			file, err := c.path.Create()
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := file.Write([]byte(uuid.New().String())); err != nil {
				return err
			}
		} else {
			if err := c.path.MkdirAll(); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestCreateOverwrite(t *testing.T) {
	WithTempDir(func(dir string) {
		filename := fs.Path(dir).Join("file.txt")
		f, err := filename.Create()
		if err != nil {
			t.Error(err)
			return
		}

		_, _ = f.Write([]byte("foo bar baz"))
		f.Close()

		f, err = filename.Create()
		if err != nil {
			t.Error(err)
			return
		}
		_, _ = f.Write([]byte("ok"))
		f.Close()

		b, err := filename.ReadAll()
		if err != nil {
			t.Error(err)
			return
		}

		str := string(b)
		if str != "ok" {
			t.Errorf("Unexpected value: '%s'", str)
		}
	})
}

func TestCreateAppendErr(t *testing.T) {
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
			{path: "", expected: fs.ErrPathIsEmpty},
			{path: root.Join("does"), expected: fs.ErrPathIsDirectory},
		}

		for i, test := range tests {
			path := test.path

			if _, err := path.Create(); err != test.expected {
				t.Errorf("Case %d, error testing create: expected '%v', received '%v'", i, test.expected, err)
			}

			if _, err := fs.Create(path.String()); err != test.expected {
				t.Errorf("Case %d, error testing open: expected '%v', received '%v'", i, test.expected, err)
			}

			if _, err := path.Append(); err != test.expected {
				t.Errorf("Case %d, error testing create: expected '%v', received '%v'", i, test.expected, err)
			}

			if _, err := fs.Append(path.String()); err != test.expected {
				t.Errorf("Case %d, error testing open: expected '%v', received '%v'", i, test.expected, err)
			}
		}
	})
}
