package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"waybackdownloader/cmd/data"
)

func WaybackLinksCollectionSave(websiteURL string, folderPath string) {
	resp, err := http.Get(fmt.Sprintf(`https://web.archive.org/web/timemap/json?url=%s&matchType=prefix&output=json&collapse=urlkey`, websiteURL))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("Failed to fetch data. HTTP Status: %d\n", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Error reading response body: %v\n", err))
	}

	var requestData [][]string
	if err := json.Unmarshal(body, &requestData); err != nil {
		panic(fmt.Sprintf("Error parsing JSON: %v\n", err))
	}

	file, err := os.Create(path.Join(folderPath, data.LINKS_FILE_NAME))
	if err != nil {
		panic(fmt.Sprintf("Error creating file: %v\n", err))
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(requestData); err != nil {
		panic(fmt.Sprintf("Error writing to file: %v\n", err))
	}
}

func WaybackDownloadFile(folderURL, mimeType, originalURL, timestamp string) (succeess bool) {
	fileName := filepath.Base(originalURL)

	resp, err := http.Get(fmt.Sprintf(`https://web.archive.org/web/%s/%s`, timestamp, originalURL))
	if err != nil {
		panic(fmt.Sprintf("Error fetching the file: %v\n", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	parts := strings.Split(mimeType, "/")
	if len(parts) > 0 {
		mimeType = parts[0]
	}

	file, err := os.Create(path.Join(data.MAIN_PATH, folderURL, mimeType, fileName))
	if err != nil {
		panic(fmt.Sprintf("Error creating file: %v\n", err))
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Error saving file: %v\n", err))
	}

	return true
}
