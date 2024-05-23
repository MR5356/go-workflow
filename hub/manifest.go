package hub

import (
	"errors"
	"runtime"
)

var (
	ErrNotSupportedMediaType = errors.New("not supported media type")
	ErrNotSupportedPlatform  = errors.New("not supported platform")
)

const (
	ManifestMediaType     = "application/plugin.workflow.manifest.v1+json"
	ManifestListMediaType = "application/plugin.workflow.manifest.list.v1+json"
)

type Manifest struct {
	MediaType string      `json:"mediaType"`
	Size      int64       `json:"size"`
	Digest    string      `json:"digest"`
	Platform  *Platform   `json:"platform"`
	Manifests []*Manifest `json:"manifests"`
}

type Platform struct {
	Architecture string `json:"architecture"`
	OS           string `json:"os"`
}

func (m *Manifest) GetDigest() (string, error) {
	if m.MediaType == ManifestMediaType {
	}
	switch m.MediaType {
	case ManifestMediaType:
		return m.Digest, nil
	case ManifestListMediaType:
		os := runtime.GOOS
		arch := runtime.GOARCH
		for _, manifest := range m.Manifests {
			if manifest.Platform.OS == os && manifest.Platform.Architecture == arch {
				return manifest.Digest, nil
			}
		}
		return "", ErrNotSupportedPlatform
	default:
		return "", ErrNotSupportedMediaType
	}
}
