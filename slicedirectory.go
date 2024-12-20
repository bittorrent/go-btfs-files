package files

import (
	"errors"
	"os"
	"sort"
	"time"
)

type fileEntry struct {
	name string
	file Node
}

func (e fileEntry) Name() string {
	return e.name
}

func (e fileEntry) Node() Node {
	return e.file
}

func FileEntry(name string, file Node) DirEntry {
	return fileEntry{
		name: name,
		file: file,
	}
}

type sliceIterator struct {
	files []DirEntry
	n     int
}

func (it *sliceIterator) Name() string {
	return it.files[it.n].Name()
}

func (it *sliceIterator) Node() Node {
	return it.files[it.n].Node()
}

func (it *sliceIterator) Next() bool {
	it.n++
	return it.n < len(it.files)
}

func (it *sliceIterator) Err() error {
	return nil
}

func (it *sliceIterator) BreadthFirstTraversal() {
}

// SliceFile implements Node, and provides simple directory handling.
// It contains children files, and is created from a `[]Node`.
// SliceFiles are always directories, and can't be read from or closed.
type SliceFile struct {
	files []DirEntry
	stat  os.FileInfo
}

func NewMapDirectory(f map[string]Node) Directory {
	ents := make([]DirEntry, 0, len(f))
	for name, nd := range f {
		ents = append(ents, FileEntry(name, nd))
	}
	sort.Slice(ents, func(i, j int) bool {
		return ents[i].Name() < ents[j].Name()
	})

	return NewSliceDirectory(ents)
}

func NewSliceDirectory(files []DirEntry) Directory {
	return &SliceFile{files: files}
}

func (f *SliceFile) Entries() DirIterator {
	return &sliceIterator{files: f.files, n: -1}
}

func (f *SliceFile) Close() error {
	return nil
}

func (f *SliceFile) Length() int {
	return len(f.files)
}

func (f *SliceFile) Mode() os.FileMode {
	if f.stat != nil {
		return f.stat.Mode()
	}
	return 0
}

func (f *SliceFile) ModTime() time.Time {
	if f.stat != nil {
		return f.stat.ModTime()
	}
	return time.Time{}
}

func (f *SliceFile) Size() (int64, error) {
	var size int64

	for _, file := range f.files {
		s, err := file.Node().Size()
		if err != nil {
			return 0, err
		}
		size += s
	}

	return size, nil
}

func (f *SliceFile) SetSize(size int64) error {
	return errors.New("not supported")
}

func (f *SliceFile) IsReedSolomon() bool {
	return false
}

func IsMapDirectory(d Directory) bool {
	if _, ok := d.(*SliceFile); ok {
		return true
	} else {
		return false
	}
}

var _ Directory = &SliceFile{}
var _ DirEntry = fileEntry{}
