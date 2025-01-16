package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

func WaybackCollectionSave(websiteURL string, folderPath string) {
	resp, err := http.Get(fmt.Sprintf(`https://web.archive.org/web/timemap/json?url=%s&matchType=prefix&output=json`, websiteURL))
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

	var data [][]string
	if err := json.Unmarshal(body, &data); err != nil {
		panic(fmt.Sprintf("Error parsing JSON: %v\n", err))
	}

	fileName := "links.json"
	file, err := os.Create(path.Join(folderPath, fileName))
	if err != nil {
		panic(fmt.Sprintf("Error creating file: %v\n", err))
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		panic(fmt.Sprintf("Error writing to file: %v\n", err))
	}
}
