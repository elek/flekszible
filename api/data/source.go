package data

import (
	"os"
	"path"
	"regexp"
)

type SourceCacheManager struct {
	RootPath string
}

func NewSourceCacheManager(root string) SourceCacheManager {
	return SourceCacheManager{
		RootPath: root,
	}
}
func (manager *SourceCacheManager) GetCacheDir(id string) string {
	cacheDir := path.Join(manager.RootPath, ".cache", id)
	return cacheDir
}

func cleanUrl(s string) string {
	var re = regexp.MustCompile("[^A-Za-z0-9\\.]")
	return re.ReplaceAllString(s, `_`)
}

type Source interface {
	GetPath(manager *SourceCacheManager, relativeDir string) (string, error)
	ToString() (string, string)
}

type CurrentDir struct {
	CurrentDir string
}

func (source *CurrentDir) GetPath(manager *SourceCacheManager, relativeDir string) (string, error) {
	return path.Join(source.CurrentDir, relativeDir), nil

}
func (source *CurrentDir) ToString() (string, string) {
	return "current dir", "."
}

type LocalSource struct {
	BaseDir     string
	RelativeDir string
}

func (source *LocalSource) GetPath(manager *SourceCacheManager, relativeDir string) (string, error) {
	return path.Join(source.BaseDir, source.RelativeDir, relativeDir), nil

}
func (source *LocalSource) ToString() (string, string) {
	return "local dir", source.RelativeDir
}

type EnvSource struct {
}

func (source *EnvSource) GetPath(manager *SourceCacheManager, relativeDir string) (string, error) {
	if os.Getenv("FLEKSZIBLE_PATH") != "" {
		return path.Join(os.Getenv("FLEKSZIBLE_PATH"), relativeDir), nil
	}
	return "", nil
}

func (source *EnvSource) ToString() (string, string) {
	return "$FLEKSZIBLE_PATH", os.Getenv("FLEKSZIBLE_PATH")
}
