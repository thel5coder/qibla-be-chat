package helper

import (
	"path/filepath"
	"qibla-backend-chat/mongomodel"
	"qibla-backend-chat/pkg/str"
)

var (
	// ImageExtentionWhitelist ...
	ImageExtentionWhitelist = []string{".jpg", ".jpeg", ".png"}
	// VideoExtentionWhitelist ...
	VideoExtentionWhitelist = []string{
		".webm", ".mkv", ".flv", ".vob", ".ogv", ".ogg", ".drc", ".gif", ".gifv", ".mng", ".avi", ".mov",
		".qt", ".wmv", ".mp4", ".m4p", ".mpg", ".mp2", ".mpeg", ".mpe", ".mpv", ".m2v", ".m4v", ".svi",
		".3gp", ".3g2", ".mxf", ".roq", ".nsv", ".f4v", ".f4p", ".f4a", ".f4b",
	}
	// AudioExtentionWhitelist ...
	AudioExtentionWhitelist = []string{
		".aiff", ".ape", ".flac", ".mp3", ".m4p", "m4a", "m4b", ".mmf", ".ogg", ".oga", ".mogg",
		".wav", ".wma", ".wv", ".webm",
	}
)

// GetChatFileType ...
func GetChatFileType(fileName string) string {
	ext := filepath.Ext(fileName)

	if str.Contains(ImageExtentionWhitelist, ext) {
		return mongomodel.ChatTypeImage
	} else if str.Contains(VideoExtentionWhitelist, ext) {
		return mongomodel.ChatTypeVideo
	} else if str.Contains(AudioExtentionWhitelist, ext) {
		return mongomodel.ChatTypeAudio
	}

	return mongomodel.ChatTypeFile
}
