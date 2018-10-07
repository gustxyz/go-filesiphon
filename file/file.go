package file

import (
	"os"
	"time"
)

type File struct {
	FName     string `json:"name"`
	FType     string `json:"type"`
	FTime     int64  `json:"time"`
	FSize     int64  `json:"size"`
	CanRename bool   `json:"can_rename,omitempty"`
	CanMove   bool   `json:"can_move_directory,omitempty"`
	CanDelete bool   `json:"can_delete,omitempty"`
}

func (f File) Name() string {
	return f.FName
}
func (f File) Size() int64 {
	return f.FSize
}
func (f File) Mode() os.FileMode {
	return 0
}
func (f File) ModTime() time.Time {
	return time.Now()
}
func (f File) IsDir() bool {
	if f.FType != "directory" {
		return false
	}
	return true
}
func (f File) Sys() interface{} {
	return nil
}
