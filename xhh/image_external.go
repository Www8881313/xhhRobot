package xhh

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
	"xhhrobot/config"
)

func UploadToExternalImageHost(imageBytes []byte, sourcePath string, dryRun bool) (XHHCOSUploadPlan, error) {
	cfg := config.ConfigStruct.Image
	if strings.TrimSpace(cfg.ExternalDir) == "" || strings.TrimSpace(cfg.ExternalBaseUrl) == "" {
		return XHHCOSUploadPlan{}, errors.New("missing external image host config: image.externalDir and image.externalBaseUrl are required")
	}

	filename := externalImageFilename(imageBytes, sourcePath)
	publicURL, err := joinExternalImageURL(cfg.ExternalBaseUrl, filename)
	if err != nil {
		return XHHCOSUploadPlan{}, err
	}
	plan := XHHCOSUploadPlan{
		Key:       filename,
		UploadURL: filepath.Join(cfg.ExternalDir, filename),
		CDNURL:    publicURL,
		Size:      len(imageBytes),
		DryRun:    dryRun,
	}
	if dryRun {
		return plan, nil
	}

	if err := os.MkdirAll(cfg.ExternalDir, 0755); err != nil {
		return plan, err
	}
	if err := os.WriteFile(plan.UploadURL, imageBytes, 0644); err != nil {
		return plan, err
	}
	plan.Uploaded = true
	return plan, nil
}

func externalImageFilename(imageBytes []byte, sourcePath string) string {
	ext := imageExtFromBytesOrPath(imageBytes, sourcePath)
	sum := sha256.Sum256(append(imageBytes, []byte(time.Now().Format(time.RFC3339Nano))...))
	return hex.EncodeToString(sum[:])[:32] + ext
}

func joinExternalImageURL(baseURL, filename string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil {
		return "", err
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return "", errors.New("image.externalBaseUrl must be an absolute URL")
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/") + "/" + filename
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String(), nil
}
