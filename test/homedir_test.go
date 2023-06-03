package test

import (
	"os"
	"testing"

	"github.com/Equationzhao/pathbeautify"
	"github.com/mitchellh/go-homedir"
)

func BenchmarkGetUserHomeDir(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := pathbeautify.GetUserHomeDir()
		_ = s
	}
}

func BenchmarkOSGetUserHomeDir(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s, _ := os.UserHomeDir()
		_ = s
	}
}

func BenchmarkGo_Homedir(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s, _ := homedir.Dir()
		_ = s
	}
}
