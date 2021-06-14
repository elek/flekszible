package data

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

type UpdateMode string

const (
	Always UpdateMode = "always"
	Never  UpdateMode = "never"
	Init   UpdateMode = "init"
)

type SourceCacheManager struct {
	RootPath    string
	doOnceCache map[string]bool
	UpdateMode  UpdateMode
}

func NewSourceCacheManager(root string) SourceCacheManager {
	res := SourceCacheManager{
		RootPath:   root,
		UpdateMode: Init,
	}
	if os.Getenv("FLEKSZIBLE_OFFLINE") == "true" {
		res.UpdateMode = Never
	} else if os.Getenv("FLEKSZIBLE_UPDATE") != "" {
		res.UpdateMode = UpdateMode(os.Getenv("FLEKSZIBLE_UPDATE"))
	}
	return res
}
func (manager *SourceCacheManager) GetCacheDir(id string) string {
	cacheDir := path.Join(manager.RootPath, ".cache", id)
	return cacheDir
}

func (manager *SourceCacheManager) DoOnce(cacheKey string, task func() error) error {
	if manager.doOnceCache == nil {
		manager.doOnceCache = make(map[string]bool)
	}
	if _, exists := manager.doOnceCache[cacheKey]; !exists {
		manager.doOnceCache[cacheKey] = true
		return task()
	}
	return nil
}

func cleanUrl(s string) string {
	var re = regexp.MustCompile("[^A-Za-z0-9\\.]")
	return re.ReplaceAllString(s, `_`)
}

type Source interface {
	GetPath(manager *SourceCacheManager) (string, error)
	ToString() string
}

type EnvSource struct {
	Dir string
}

func (source *EnvSource) GetPath(manager *SourceCacheManager) (string, error) {
	return filepath.Abs(source.Dir)

}
func (source *EnvSource) ToString() string {
	return "[local dir] FLEKSZIBLE_PATH=" + source.Dir
}

type LocalSource struct {
	Dir string
}

func (source *LocalSource) GetPath(manager *SourceCacheManager) (string, error) {
	return filepath.Abs(source.Dir)

}
func (source *LocalSource) ToString() string {
	return "[local dir] dir=" + source.Dir
}

func LocalSourcesFromEnv() []Source {
	sources := make([]Source, 0)
	if os.Getenv("FLEKSZIBLE_PATH") != "" {
		for _, path := range strings.Split(os.Getenv("FLEKSZIBLE_PATH"), ",") {
			sources = append(sources, &EnvSource{Dir: path})
		}
	}
	return sources
}

type RemoteSource struct {
	Url      string
	CacheDir string
}

func (source *RemoteSource) ToString() string {
	return "[remote] url=" + source.Url
}

type Downloader interface {
	Download(url string, destinationDir string, rootPath string) error
}

var DownloaderPlugin Downloader = NodeDownloader{}

type NodeDownloader struct {
}

func (NodeDownloader) Download(url string, destinationDir string, rootPath string) error {
	logrus.Warn("Downloader component has not been registered. Downloading remote resources are disabled.")
	return nil
}

func (source *RemoteSource) EnsureDownloaded(manager *SourceCacheManager) error {
	destinationDir := manager.GetCacheDir(cleanUrl(source.Url))
	if _, err := os.Stat(destinationDir); !os.IsNotExist(err) && manager.UpdateMode == Init {
		return nil
	}
	task := func() error {
		return DownloaderPlugin.Download(source.Url, destinationDir, manager.RootPath)
	}
	return manager.DoOnce(destinationDir, task)
}

func (source *RemoteSource) GetPath(manager *SourceCacheManager) (string, error) {
	err := source.EnsureDownloaded(manager)
	if err != nil {
		return "", err
	}
	return manager.GetCacheDir(cleanUrl(source.Url)), nil
}
