package data

import (
	"github.com/hashicorp/go-getter"
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
func (manager *SourceCacheManager) GetCacheDir(source Source) string {
	cleanUrl := cleanUrl(source.Url)
	cacheDir := path.Join(manager.RootPath, ".cache", cleanUrl)
	return cacheDir
}

func (manager *SourceCacheManager) EnsureDownloaded(source Source) error {
	setPwd := func(client *getter.Client) error { client.Pwd = manager.RootPath; return nil; }
	return getter.Get(manager.GetCacheDir(source), source.Url, setPwd)

}

func cleanUrl(s string) string {
	var re = regexp.MustCompile("[^A-Za-z0-9\\.]")
	return re.ReplaceAllString(s, `_`)
}

type Source struct {
	Url      string
	CacheDir string
}
