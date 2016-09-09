// Author: ZHU HAIHUA
// Date: 9/9/16
package mapbytefs // import "github.com/kimiazhu/vfs/mapbytefs"
import (
	"github.com/kimiazhu/vfs"
	"strings"
	"os"
	pathpkg "path"
	"io"
	"time"
	"sort"
	"bytes"
)

type mapByteFS map[string][]byte

// New returns a new FileSystem from the provided map.
// Map keys should be forward slash-separated pathnames
// and not contain a leading slash.
// Map values are byte slice
func New(m map[string][]byte) vfs.FileSystem {
	return mapByteFS(m)
}

func (fs mapByteFS) String() string {
	return "mapbytefs"
}

func (fs mapByteFS) Close() error {
	return nil
}

func filename(p string) string {
	return strings.TrimPrefix(p, "/")
}

func (fs mapByteFS) Open(p string) (vfs.ReadSeekCloser, error) {
	b, ok := fs[filename(p)]
	if !ok {
		return nil, os.ErrNotExist
	}
	return nopCloser{bytes.NewReader(b)}, nil
}

func fileInfo(name string, content []byte) os.FileInfo {
	return mapFI{name: pathpkg.Base(name), size: len(content)}
}

func dirInfo(name string) os.FileInfo {
	return mapFI{name: pathpkg.Base(name), dir: true}
}

func (fs mapByteFS) Stat(p string) (os.FileInfo, error) {
	return fs.Lstat(p)
}

// slashdir returns path.Dir(p), but special-cases paths not beginning
// with a slash to be in the root.
func (fs mapByteFS) Lstat(p string) (os.FileInfo, error) {
	b, ok := fs[filename(p)]
	if ok {
		return fileInfo(p, b), nil
	}
	ents, _ := fs.ReadDir(p)
	if len(ents) > 0 {
		return dirInfo(p), nil
	}
	return nil, os.ErrNotExist
}

func slashdir(p string) string {
	d := pathpkg.Dir(p)
	if d == "." {
		return "/"
	}
	if strings.HasPrefix(p, "/") {
		return d
	}
	return "/" + d
}

func (fs mapByteFS) ReadDir(p string) ([]os.FileInfo, error) {
	p = pathpkg.Clean(p)
	var ents []string
	fim := make(map[string]os.FileInfo)
	for fn, b := range fs {
		dir := slashdir(fn)
		isFile := true
		var lastBase string
		for {
			if dir == p {
				base := lastBase
				if isFile {
					base = pathpkg.Base(fn)
				}
				if fim[base] == nil {
					var fi os.FileInfo
					if isFile {
						fi = fileInfo(fn, b)
					} else {
						fi = dirInfo(base)
					}
					ents = append(ents, base)
					fim[base] = fi
				}
			}
			if dir == "/" {
				break
			} else {
				isFile = false
				lastBase = pathpkg.Base(dir)
				dir = pathpkg.Dir(dir)
			}
		}
	}
	if len(ents) == 0 {
		return nil, os.ErrNotExist
	}

	sort.Strings(ents)
	var list []os.FileInfo
	for _, dir := range ents {
		list = append(list, fim[dir])
	}
	return list, nil
}

// mapFI is the map-based implementation of FileInfo.
type mapFI struct {
	name string
	size int
	dir  bool
}

func (fi mapFI) IsDir() bool {
	return fi.dir
}
func (fi mapFI) ModTime() time.Time {
	return time.Time{}
}
func (fi mapFI) Mode() os.FileMode {
	if fi.IsDir() {
		return 0755 | os.ModeDir
	}
	return 0444
}
func (fi mapFI) Name() string {
	return pathpkg.Base(fi.name)
}
func (fi mapFI) Size() int64 {
	return int64(fi.size)
}
func (fi mapFI) Sys() interface{} {
	return nil
}

type nopCloser struct {
	io.ReadSeeker
}

func (nc nopCloser) Close() error {
	return nil
}
