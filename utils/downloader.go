package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

// ProgressTracker implements io.Writer to count bytes and track progress.
type DownloadProgressTracker struct {
	Total   int64
	Current int64
}

// Write intercepts bytes to update progress and display to terminal.
func (pt *DownloadProgressTracker) Write(p []byte) (int, error) {
	n := len(p)
	pt.Current += int64(n)
	if pt.Total > 0 {
		percent := float64(pt.Current) / float64(pt.Total) * 100
		fmt.Printf("\r [+] Downloading... %.2f%% (%d/%d bytes)", percent, pt.Current, pt.Total)
	}
	return n, nil
}

func Download(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, _ := os.Create(dest)
	defer out.Close()

	size, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)

	// TeeReader pipes data to file and tracker concurrently
	src := io.TeeReader(resp.Body, &DownloadProgressTracker{Total: size})
	io.Copy(out, src)
	fmt.Println("\n [+] Download finished.")

	return nil
}
