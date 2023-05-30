package pathbeautify

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTransform(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name            string
		args            args
		wantTransformed string
	}{
		{
			name:            "homedir",
			args:            args{path: "~"},
			wantTransformed: GetUserHomeDir(),
		},
		{
			name:            "homedir+path",
			args:            args{path: "~/a/b/c"},
			wantTransformed: filepath.Join(GetUserHomeDir(), "a", "b", "c"),
		},
		{
			name:            "normal path",
			args:            args{path: "a/b/c"},
			wantTransformed: filepath.Join("a", "b", "c"),
		},
		{
			name:            "dots3",
			args:            args{path: "..."},
			wantTransformed: filepath.Join("..", ".."),
		},
		{
			name:            "dots4",
			args:            args{path: "...."},
			wantTransformed: filepath.Join("..", "..", ".."),
		},
		{
			name:            "dots5",
			args:            args{path: "....."},
			wantTransformed: filepath.Join("..", "..", "..", ".."),
		},
		{
			name:            "dots5",
			args:            args{path: "...../a/b/c"},
			wantTransformed: filepath.Join("..", "..", "..", "..", "a", "b", "c"),
		},
		{
			name:            "dots/dots/path",
			args:            args{path: "..../../.../a/b/c"},
			wantTransformed: filepath.Join("..", "..", "..", "..", "..", "..", "a", "b", "c"),
		},
		{
			name:            "dots/dots",
			args:            args{path: "../.../..../.."},
			wantTransformed: filepath.Join("..", "..", "..", "..", "..", "..", ".."),
		},
		{
			name:            "dots/dots",
			args:            args{path: "a/b/.../.../../c"},
			wantTransformed: strings.Join([]string{"a", "b", "..", "..", "..", "..", "..", "c"}, string(filepath.Separator)),
		},
		{
			name:            "root",
			args:            args{path: "/"},
			wantTransformed: string(filepath.Separator),
		},
		{
			name:            "root",
			args:            args{path: "/a/b/c"},
			wantTransformed: filepath.Join("/", "a", "b", "c"),
		},
		{
			name:            "root/dots",
			args:            args{path: "/a/b/.../c"},
			wantTransformed: strings.Join([]string{"", "a", "b", "..", "..", "c"}, string(filepath.Separator)),
		},
		{
			name:            "./path",
			args:            args{path: "./a/b/c"},
			wantTransformed: strings.Join([]string{".", "a", "b", "c"}, string(filepath.Separator)),
		},
		{
			name:            "./",
			args:            args{path: "./"},
			wantTransformed: filepath.Join(".", ""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTransformed := Transform(tt.args.path); gotTransformed != tt.wantTransformed {
				t.Errorf("Transform() = %v, want %v", gotTransformed, tt.wantTransformed)
			}
		})
	}
}

func BenchmarkGetUserHomeDir(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := GetUserHomeDir()
		_ = s
	}
}

func BenchmarkOSGetUserHomeDir(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s, _ := os.UserHomeDir()
		_ = s
	}
}
