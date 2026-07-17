package imageprovider

import (
	"embed"
)

//go:embed assets/*
var assetsFS embed.FS

type ImageProvider struct {
}

func New() *ImageProvider {
	return &ImageProvider{}
}
