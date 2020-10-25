package archive

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

// FileInterface ...
type FileInterface interface {
	HasDownload() bool
	Download(uri string) error
	Path() string
	HasPeaks() bool
	PeakPath() string
	HasSegments() string
	SegmentsPath() string
}

// File ...
type File struct {
	id     string
	tmpDir string
}

// NewFile ...
func NewFile(id string) *File {
	return &File{
		id:     id,
		tmpDir: "/tmp/annotation-agent",
	}
}

// Path ...
func (f *File) Path() string {
	return path.Join(f.tmpDir, fmt.Sprintf("%s.flac", f.id))
}

// PeakPath ...
func (f *File) PeakPath() string {
	return path.Join(f.tmpDir, fmt.Sprintf("%s.dat", f.id))
}

// SegmentsPath ...
func (f *File) SegmentsPath() string {
	return path.Join(f.tmpDir, fmt.Sprintf("%s.segments.csv", f.id))
}

// HasPeaks tells if we have a peaks.js dat file
func (f *File) HasPeaks() bool {
	return fileExists(f.PeakPath())
}

// HasSegments ...
func (f *File) HasSegments() bool {
	return fileExists(f.SegmentsPath())
}

// HasDownload ...
func (f *File) HasDownload() bool {
	return fileExists(f.Path())
}

// Download ...
func (f *File) Download(uri string) error {
	path := f.Path()
	log := logrus.
		WithField("path", path).
		WithField("uri", uri)

	if fileExists(path) {
		log.Info("File exists.")
		return nil
	}

	log.Info("Starting download.")
	out, err := os.Create(path)
	if err != nil {
		log.WithError(err).Error()
		return err
	}
	defer out.Close()

	resp, err := http.Get(uri)
	if err != nil {
		log.WithError(err).Error()
		return err
	}
	defer resp.Body.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		log.WithError(err).Error()
		return err
	}
	log.Info("Download successful")
	return nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
