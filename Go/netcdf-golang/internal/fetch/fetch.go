package fetch

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// GetNWM take a URL and loads the URL into a
// local scratch or virtual file system
func GetNWM(dir, dataurl string) (string, error) {
	u, err := url.Parse(dataurl)
	if err != nil {
		panic(err)
	}

	base := filepath.Base(u.Path)
	fp := fmt.Sprintf("%s/%s", dir, base)
	err = downloadFile(fp, dataurl)

	return fp, err
}

func downloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
