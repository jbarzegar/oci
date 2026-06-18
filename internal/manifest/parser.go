package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ManifestV2Config struct {
	MediaType string `json:"mediaType"`
	Size      int    `json:"size"`
	Digest    string `json:"digest"`
}

type ManifestV2 struct {
	SchemaVersion int                `json:"schemaVersion"`
	MediaType     string             `json:"mediaType"`
	Config        ManifestV2Config   `json:"config"`
	Layers        []ManifestV2Config `json:"layers"`
}

var (
	ErrManifestInvalid = errors.New("manifest could not be parsed")
)

type ManifestV2Reader struct{}

func NewV2Parser() *ManifestV2Reader {
	return &ManifestV2Reader{}
}

func UnmarshalV2(data []byte) (ManifestV2, error) {
	var m ManifestV2
	if err := json.Unmarshal(data, &m); err != nil {
		fmt.Println(err)
		// unhandled err
		return ManifestV2{}, err
	}

	return m, nil
}

func MarshalV2(m ManifestV2) ([]byte, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}
