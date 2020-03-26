package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func Unzip(zipFile string, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		fPath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(fPath, os.ModePerm)
			continue
		}
		if err = os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
			return err
		}

		var err = func() error {
			in, err := f.Open()
			if err != nil {
				return err
			}
			defer in.Close()
			out, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer out.Close()

			if _, err = io.Copy(out, in);err != nil {
				return err
			}
			return nil
		}()
		if err != nil {
			return err
		}
	}
	return nil
}
