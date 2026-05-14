package manifest

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

const manifestURL = "https://launchermeta.mojang.com/mc/game/version_manifest.json"

type manifest struct {
	Versions []Version `json:"versions"`
}

type Version struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

type VersionDetails struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Downloads struct {
		Client struct {
			URL  string `json:"url"`
			Size int    `json:"size"`
			SHA1 string `json:"sha1"`
		} `json:"client"`
	} `json:"downloads"`
}

func GetVersionList() ([]Version, error) {
	_, b, err := fasthttp.Get(nil, manifestURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch manifest")
	}

	var m manifest
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, errors.Wrap(err, "failed to parse manifest")
	}

	return m.Versions, nil
}

func FilterVersionList(versions []Version, releaseOnly bool, newerThan string) []Version {
	var filtered []Version
	for _, v := range versions {
		if releaseOnly && v.Type != "release" {
			continue
		}
		if newerThan != "" && v.ID == newerThan {
			break
		}
		filtered = append(filtered, v)
	}
	return filtered
}

func (v Version) GetDetails() (*VersionDetails, error) {
	_, b, err := fasthttp.Get(nil, v.URL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch version details")
	}

	var vd VersionDetails
	if err := json.Unmarshal(b, &vd); err != nil {
		return nil, errors.Wrap(err, "failed to parse version details")
	}

	return &vd, nil
}

func (vd *VersionDetails) FetchClientJar() ([]byte, error) {
	_, b, err := fasthttp.Get(nil, vd.Downloads.Client.URL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch client jar")
	}
	return b, nil
}

func (v Version) FetchClientJar() ([]byte, error) {
	vd, err := v.GetDetails()
	if err != nil {
		return nil, err
	}
	return vd.FetchClientJar()
}
