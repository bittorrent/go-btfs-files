package files

import (
	"os"
	"strings"
	"time"
)

type Symlink struct {
	Target string

	mtime  time.Time
	stat   os.FileInfo
	reader strings.Reader
}

func NewLinkFile(target string, stat os.FileInfo) File {
	lf := &Symlink{Target: target, stat: stat}
	if stat != nil && stat.ModTime() != (time.Time{}) {
		lf.mtime = stat.ModTime()
	}
	lf.reader.Reset(lf.Target)
	return lf
}

func (lf *Symlink) Close() error {
	return nil
}

func (lf *Symlink) Read(b []byte) (int, error) {
	return lf.reader.Read(b)
}

func (lf *Symlink) Seek(offset int64, whence int) (int64, error) {
	return lf.reader.Seek(offset, whence)
}

func (lf *Symlink) Mode() os.FileMode {
	return os.ModeSymlink | os.ModePerm
}

func (lf *Symlink) ModTime() time.Time {
	return lf.mtime
}

func (lf *Symlink) Size() (int64, error) {
	return lf.reader.Size(), nil
}

func ToSymlink(n Node) *Symlink {
	l, _ := n.(*Symlink)
	return l
}

var _ File = &Symlink{}
