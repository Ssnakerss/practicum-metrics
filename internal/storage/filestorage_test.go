package storage

import "testing"

func TestFileStorage_New(t *testing.T) {
	f := FileStorage{}
	f.New("file name")
}
