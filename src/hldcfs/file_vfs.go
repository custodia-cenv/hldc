package hldcfs

import (
	"io"
	"syscall"
	"time"
)

func (o *VfsFile) Close() error {
	return nil
}

func (o *VfsFile) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (o *VfsFile) ReadAt(b []byte, off int64) (n int, err error) {
	return 0, nil
}

func (o *VfsFile) Write(b []byte) (n int, err error) {
	return 0, nil
}

func (o *VfsFile) WriteAt(b []byte, off int64) (n int, err error) {
	return 0, nil
}

func (o *VfsFile) ReadFrom(r io.Reader) (n int64, err error) {
	return 0, nil
}

func (o *VfsFile) Seek(offset int64, whence int) (ret int64, err error) {
	return 0, nil
}

func (o *VfsFile) SetDeadline(t time.Time) error {
	return nil
}

func (o *VfsFile) SetWriteDeadline(t time.Time) error {
	return nil
}

func (o *VfsFile) Stat() (*VfsFileInfo, error) {
	return nil, nil
}

func (o *VfsFile) Sync() error {
	return nil
}

func (o *VfsFile) SyscallConn() (syscall.RawConn, error) {
	return nil, nil
}

func (o *VfsFile) Truncate(size int64) error {
	return nil
}

func (o *VfsFile) WriteString(s string) (n int, err error) {
	return 0, nil
}

func (o *VfsFile) WriteTo(w io.Writer) (n int64, err error) {
	return 0, nil
}
