package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"waybackdownloader/cmd/data"
	"waybackdownloader/cmd/data/config"
	"waybackdownloader/cmd/util"
)

func WaybackLinksCollectionSave(config *config.Config, websiteURL string) {
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
		panic(err)
	}

	var requestData [][]string
	if err := json.Unmarshal(body, &requestData); err != nil {
		panic(err)
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
			WebsiteURL: util.RemoveSlashFromString(websiteURL),
		}
	}

	_, err = config.DB.InsertURLs(linksArr)
	if err != nil {
		panic(err)
	}
}

func WaybackDownloadFile(config *config.Config, link data.Link) (succeess bool) {
	fileName := filepath.Base(link.Original)

	resp, err := http.Get(fmt.Sprintf(`https://web.archive.org/web/%sif_/%s`, link.Timestamp, link.Original))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	filePath := path.Join(data.MAIN_PATH, link.WebsiteURL, util.RemoveSlashFromString(link.Mimetype), fileName)

	dirPath := path.Dir(filePath)

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		panic(err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		panic(err)
	}

	err = config.DB.UpdateURL(link)
	if err != nil {
		panic(err)
	}

	return true
}
