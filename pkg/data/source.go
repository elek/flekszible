package data

import (
	"github.com/hashicorp/go-getter"
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
type LocalSource struct {
	RelativeTo string
}

func (source *LocalSource) GetPath(manager *SourceCacheManager, relativeDir string) (string, error) {
	return path.Join(source.RelativeTo, relativeDir), nil

}
func (source *LocalSource) ToString() (string, string) {
	return "current dir", source.RelativeTo
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

type GoGetter struct {
	Url      string
	CacheDir string
}

func (source *GoGetter) ToString() (string, string) {
	return "GoGetter", source.Url
}

func (source *GoGetter) EnsureDownloaded(manager *SourceCacheManager) error {
	setPwd := func(client *getter.Client) error { client.Pwd = manager.RootPath; return nil; }
	return getter.Get(manager.GetCacheDir(cleanUrl(source.Url)), source.Url, setPwd)
}

func (source *GoGetter) GetPath(manager *SourceCacheManager, relativeDir string) (string, error) {
	err := source.EnsureDownloaded(manager)
	if err != nil {
		return "", err
	}
	return manager.GetCacheDir(cleanUrl(source.Url)), nil
}
