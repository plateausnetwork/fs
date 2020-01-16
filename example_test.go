package fs_test

import (
	"fmt"

	"github.com/rhizomplatform/fs"
)

func ExamplePath_joining() {
	user := "username"
	pkg := "mypackage"

	home := fs.Path("/usr/home").Join(user)
	gosrc := fs.Path("go").Join("src")
	github := fs.Path("github.com")

	fmt.Println(home.JoinP(gosrc).JoinP(github).Join(user).Join(pkg))
	// Output: /usr/home/username/go/src/github.com/username/mypackage
}

func ExamplePath_sanitize() {
	root := fs.Path("/")
	home := fs.Path("home/")
	username := fs.Path("/username/")
	file := fs.Path("/myfile.txt")

	fmt.Println(root.JoinP(home).JoinP(username).JoinP(file))
	// Output: /home/username/myfile.txt
}
