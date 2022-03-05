package backuper

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"time"
)

func writeBackup(src io.Reader, outPath string) error {
	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer outFile.Close()

	return writeGzip(src, outFile)
}

func writeGzip(src io.Reader, dst io.Writer) error {
	gz, err := gzip.NewWriterLevel(dst, gzip.DefaultCompression)
	if err != nil {
		return fmt.Errorf("create gzip writer: %w", err)
	}
	defer gz.Close()
	defer gz.Flush()

	_, err = io.Copy(gz, src)
	if err != nil {
		return fmt.Errorf("gzip src: %w", err)
	}

	return nil
}

func makeRotateTag(maxRotate int, date time.Time, filename string) string {
	year, month, day := date.Date()
	rotateNumber := (year + int(month) + day) % maxRotate
	return fmt.Sprintf("rotate.%d-%s", rotateNumber, filename)
}
