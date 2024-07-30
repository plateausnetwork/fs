# Rhizom File System

A simple package for handling file system problems easily

## Building from source

To build from source, you will need the following prerequisites:

- Go 1.13 or greater;
- Git

### Downloading the code

First, clone the project:

```bash
git clone git@github.com:plateausnetwork/fs.git /your/directory/of/choice/rhizom
cd /your/directory/of/choice/rhizom
```

### Testing

To run the tests, try `go test`.

### Using as a library

```go
import (
  "github.com/plateausnetwork/fs"
)

func myFunc() {
  path := fs.Path{"my/directory/path/not/exists/something"}

  // creates all directories that doesn't exists
  if err := path.MkdirAll(); err != nil {
    panic(err)
  }
}

```

## License

For more details about our license model, please take a look at the [LICENSE](LICENSE) file.

**2020**, Rhizom Platform.
