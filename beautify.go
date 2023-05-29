package pathbeautify

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/valyala/bytebufferpool"
)

var syncHomedir sync.Once
var userHomeDir string

func GetUserHomeDir() string {
	syncHomedir.Do(func() {
		userHomeDir, _ = os.UserHomeDir()
	})
	return userHomeDir
}

// Transform ~ to $HOME
// ... -> ../..
// .... -> ../../..
func Transform(path string) (transformed string) {
	switch path {
	case ".", "..":
	case "...":
		path = filepath.Join("..", "..")
	case "....":
		path = filepath.Join("..", "..", "..")
	case "":
	case "~":
		path = GetUserHomeDir()
	default:
		// ~/a/b/c
		if strings.HasPrefix(path, "~") {
			home := GetUserHomeDir()
			path = home + (path)[1:]
		}

		if strings.HasPrefix(path, string(filepath.Separator)) {
			return
		}

		// .....?
		// start from 3, aka ...
		matchDots := true
		times := -1
		for _, dot := range path {
			if dot != '.' {
				if dot != filepath.Separator {
					matchDots = false
				}
				break
			}
			times++
		}

		// case 1
		// .../a/b/c -> times = 2
		// ../../ + a/b/c -> ../../a/b/c
		// case 2
		// ... -> times = 2
		// ../../ + empty -> ../../
		// case 3
		// .../ -> times = 2
		// ../../ + empty -> ../../
		if matchDots {
			const parent = ".."
			buffer := bytebufferpool.Get()
			for i := 0; i < times; i++ {
				_, _ = buffer.WriteString(parent)
				_ = buffer.WriteByte(filepath.Separator)
			}
			if times+2 < len(path) {
				_, _ = buffer.WriteString((path)[times+2:])
			}

			path = buffer.String()
			bytebufferpool.Put(buffer)
		}

	}
	return path
}
