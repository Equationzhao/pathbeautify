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

var replacer = strings.NewReplacer(
	"\\", string(filepath.Separator),
	"/", string(filepath.Separator),
)

// CleanSeparator convert '//' and '\' to os-specific separator
func CleanSeparator(path string) string {
	return replacer.Replace(path)
}

// Beautify is alias of Transform
func Beautify(path string) string {
	return Transform(path)
}

// Transform path
//
//	~ -> $HOME
//	~/a/b/c -> $HOME/a/b/c
//	... -> ../..
//	.... -> ../../..
//	..../../.../a/b/c -> ../../../../../../a/b/c
func Transform(path string) (transformed string) {

	switch path {
	case "./":
		path = "."
	case ".", "..":
	case "...":
		path = filepath.Join("..", "..")
	case "....":
		path = filepath.Join("..", "..", "..")
	case "":
	case "~":
		path = GetUserHomeDir()
	default:
		path = CleanSeparator(path)
		// if strings.HasPrefix(path, "/") || strings.HasPrefix(path, "\\") {
		// 	return path
		// }

		// ~/a/b/c
		if strings.HasPrefix(path, "~") {
			home := GetUserHomeDir()
			path = home + (path)[1:]
		}

		names := strings.Split(path, string(filepath.Separator))

		for i := range names {
			names[i] = clean(names[i])
		}
		path = strings.Join(names, string(filepath.Separator))
	}
	return path
}

func clean(path string) string {
	// start from 3, aka ...
	matchDots := true
	times := -1
	for _, dot := range path {
		if dot != '.' {
			if !IsPathSeparator(dot) {
				matchDots = false
			}
			break
		}
		times++
	}

	if matchDots {
		path = cleanDot(path, times)
	}
	return path
}

func cleanDot(path string, times int) string {
	if times == 0 {
		return path
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
	const parent = ".."
	buffer := bytebufferpool.Get()
	for i := 0; i < times; i++ {
		_, _ = buffer.WriteString(parent)
		if i != times-1 {
			_ = buffer.WriteByte(filepath.Separator)
		}
	}
	if times+2 < len(path) {
		_, _ = buffer.WriteString((path)[times+2:])
	}

	path = buffer.String()
	bytebufferpool.Put(buffer)
	return path
}

func IsPathSeparator(c rune) bool {
	return c == '\\' || c == '/'
}
