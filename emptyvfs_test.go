// Author: ZHU HAIHUA
// Date: 9/8/16
package vfs_test

import (
	"github.com/kimiazhu/vfs"
	"github.com/kimiazhu/vfs/mapfs"
	"testing"
	"time"
)

func TestNewNameSpace(t *testing.T) {

	// We will mount this filesystem under /fs1
	mount := mapfs.New(map[string]string{"fs1file": "abcdefgh"})

	// Existing process. This should give error on Stat("/")
	t1 := vfs.NameSpace{}
	t1.Bind("/fs1", mount, "/", vfs.BindReplace)

	// using NewNameSpace. This should work fine.
	t2 := vfs.NewNameSpace()
	t2.Bind("/fs1", mount, "/", vfs.BindReplace)

	testcases := map[string][]bool{
		"/":            []bool{false, true},
		"/fs1":         []bool{true, true},
		"/fs1/fs1file": []bool{true, true},
	}

	fss := []vfs.FileSystem{t1, t2}

	for j, fs := range fss {
		for k, v := range testcases {
			_, err := fs.Stat(k)
			result := err == nil
			if result != v[j] {
				t.Errorf("fs: %d, testcase: %s, want: %v, got: %v, err: %s", j, k, v[j], result, err)
			}
		}
	}

	fi, err := t2.Stat("/")
	if err != nil {
		t.Fatal(err)
	}

	if fi.Name() != "/" {
		t.Errorf("t2.Name() : want:%s got:%s", "/", fi.Name())
	}

	if !fi.ModTime().IsZero() {
		t.Errorf("t2.Modime() : want:%v got:%v", time.Time{}, fi.ModTime())
	}
}
