// Author: ZHU HAIHUA
// Date: 9/9/16
package mapbytefs

import (
	"testing"
	"io/ioutil"
	"os"
	"reflect"
)

func TestOpenRoot(t *testing.T) {
	fs := New(map[string][]byte{
		"foo/bar/three.txt": []byte("a"),
		"foo/bar.txt":       []byte("b"),
		"top.txt":           []byte("c"),
		"other-top.txt":     []byte("d"),
	})
	tests := []struct {
		path string
		want []byte
	}{
		{"/foo/bar/three.txt", []byte("a")},
		{"foo/bar/three.txt", []byte("a")},
		{"foo/bar.txt", []byte("b")},
		{"top.txt", []byte("c")},
		{"/top.txt", []byte("c")},
		{"other-top.txt", []byte("d")},
		{"/other-top.txt", []byte("d")},
	}
	for _, tt := range tests {
		rsc, err := fs.Open(tt.path)
		if err != nil {
			t.Errorf("Open(%q) = %v", tt.path, err)
			continue
		}
		slurp, err := ioutil.ReadAll(rsc)
		if err != nil {
			t.Error(err)
		}
		if string(slurp) != string(tt.want) {
			t.Errorf("Read(%q) = %q; want %q", tt.path, tt.want, slurp)
		}
		rsc.Close()
	}

	_, err := fs.Open("/xxxx")
	if !os.IsNotExist(err) {
		t.Errorf("ReadDir /xxxx = %v; want os.IsNotExist error", err)
	}
}

func TestReaddir(t *testing.T) {
	fs := New(map[string][]byte{
		"foo/bar/three.txt": []byte("333"),
		"foo/bar.txt":       []byte("22"),
		"top.txt":           []byte("top.txt file"),
		"other-top.txt":     []byte("other-top.txt file"),
	})
	tests := []struct {
		dir  string
		want []os.FileInfo
	}{
		{
			dir: "/",
			want: []os.FileInfo{
				mapFI{name: "foo", dir: true},
				mapFI{name: "other-top.txt", size: len("other-top.txt file")},
				mapFI{name: "top.txt", size: len("top.txt file")},
			},
		},
		{
			dir: "/foo",
			want: []os.FileInfo{
				mapFI{name: "bar", dir: true},
				mapFI{name: "bar.txt", size: len([]byte("22"))},
			},
		},
		{
			dir: "/foo/",
			want: []os.FileInfo{
				mapFI{name: "bar", dir: true},
				mapFI{name: "bar.txt", size: len([]byte("22"))},
			},
		},
		{
			dir: "/foo/bar",
			want: []os.FileInfo{
				mapFI{name: "three.txt", size: len([]byte("333"))},
			},
		},
	}
	for _, tt := range tests {
		fis, err := fs.ReadDir(tt.dir)
		if err != nil {
			t.Errorf("ReadDir(%q) = %v", tt.dir, err)
			continue
		}
		if !reflect.DeepEqual(fis, tt.want) {
			t.Errorf("ReadDir(%q) = %#v; want %#v", tt.dir, fis, tt.want)
			continue
		}
	}

	_, err := fs.ReadDir("/xxxx")
	if !os.IsNotExist(err) {
		t.Errorf("ReadDir /xxxx = %v; want os.IsNotExist error", err)
	}
}
