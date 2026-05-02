package file

import "testing"

func TestSanitizeFileName(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "normal", in: "hello.txt", want: "hello.txt"},
		{name: "path traversal", in: "../../etc/passwd", want: "passwd"},
		{name: "invalid chars", in: "my file@name!.txt", want: "my_file_name_.txt"},
		{name: "trim spaces and dots", in: "  .hidden.  ", want: "hidden"},
		{name: "empty", in: "", want: "file"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := SanitizeFileName(tc.in)
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestGetFileExtension(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "lower ext", in: "a.txt", want: ".txt"},
		{name: "upper ext normalized", in: "a.PDF", want: ".pdf"},
		{name: "no ext", in: "README", want: ""},
		{name: "invalid chars removed", in: "a.t!@#xt", want: ".txt"},
		{name: "truncate length", in: "a.abcdefghijklmnop", want: ".abcdefghij"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := GetFileExtension(tc.in)
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}
