package fsutil

type PathInfo []BlockInfo

func ValidateFilePath(path string) bool {
	return len(path) > 0
}
