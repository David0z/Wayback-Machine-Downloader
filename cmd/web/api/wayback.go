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
	"waybackdownloader/cmd/data/config"
)

func WaybackLinksCollectionSave(config *config.Config, websiteURL string, folderPath string) {
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

	var linksArr = make([]data.Link, len(requestData))

	for index, valArray := range requestData {
		linksArr[index] = data.Link{
			Urlkey:     valArray[0],
			Timestamp:  valArray[1],
			Original:   valArray[2],
			Mimetype:   valArray[3],
			Statuscode: valArray[4],
			Downloaded: false,
			WebsiteURL: websiteURL,
		}
	}

	_, err = config.DB.InsertURLs(linksArr)
	if err != nil {
		panic(err)
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

	// @TODO edit in database as downloaded

	return true
}
