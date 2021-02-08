package mp3mirror

import (
	"path"
)

// Mirror path mapper
type Mirror struct {
	Folder string
}

const mark = "(gen)"

// Get mirror path.
// Empty mirror means nearby
func (mirror *Mirror) Get(filepath string) string {
	ext := path.Ext(filepath)
	if ext != ".wav" && ext != ".aiff" {
		return ""
	}
	if len(mirror.Folder) == 0 {
		return filepath[:len(filepath)-len(ext)] + mark + ".mp3" // nearby
	}
	base := path.Base(filepath)
	return mirror.Folder + base[:len(base)-len(ext)] + mark + ".mp3" // mirror to folder
}
