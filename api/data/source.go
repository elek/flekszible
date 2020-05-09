package data

import (
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"regexp"
	"strings"
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
	Dir string
}

func (source *LocalSource) GetPath(manager *SourceCacheManager, relativeDir string) (string, error) {
	return path.Join(source.Dir, relativeDir), nil

}
func (source *LocalSource) ToString() (string, string) {
	return "local dir", source.Dir
}

func LocalSourcesFromEnv() []Source {
	sources := make([]Source, 0)
	if os.Getenv("FLEKSZIBLE_PATH") != "" {
		for _, path := range strings.Split(os.Getenv("FLEKSZIBLE_PATH"), ",") {
			sources = append(sources, &LocalSource{Dir: path})
		}
	}
	return sources
}

type RemoteSource struct {
	Url      string
	CacheDir string
}

func (source *RemoteSource) ToString() (string, string) {
	return "remote dir", source.Url
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
	return DownloaderPlugin.Download(source.Url, destinationDir, manager.RootPath)
}

func (source *RemoteSource) GetPath(manager *SourceCacheManager, relativeDir string) (string, error) {
	err := source.EnsureDownloaded(manager)
	if err != nil {
		return "", err
	}
	baseDir := path.Join(manager.GetCacheDir(cleanUrl(source.Url)), relativeDir)
	subDir := path.Join(baseDir, "flekszible")
	if _, err := os.Stat(subDir); !os.IsNotExist(err) {
		return subDir, nil
	}

	return baseDir, nil
}
