package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"waybackdownloader/cmd/data"
	"waybackdownloader/cmd/data/config"
	"waybackdownloader/cmd/util"

	"github.com/h2non/filetype"
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

func WaybackDownloadFile(config *config.Config, link data.Link) error {
	fileName := util.SanitizeFileName(filepath.Base(link.Original))
	ext := strings.ToLower(strings.ReplaceAll(filepath.Ext(fileName), ".", ""))

	resp, err := http.Get(fmt.Sprintf(`https://web.archive.org/web/%sif_/%s`, link.Timestamp, link.Original))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("request failed")
	}

	filePath := path.Join(data.MAIN_PATH, link.WebsiteURL, util.RemoveSlashFromString(link.Mimetype), fileName)

	dirPath := path.Dir(filePath)

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	if !filetype.IsSupported(ext) {
		file.Seek(0, 0)
		head := make([]byte, 261)
		_, err := file.Read(head)
		if err != nil && err != io.EOF {
			return err
		}

		kind, err := filetype.Match(head)
		if err != nil {
			return err
		}
		if kind != filetype.Unknown {
			newFileName := fileName + "." + kind.Extension
			newFilePath := path.Join(dirPath, newFileName)

			file.Close()
			err = os.Rename(filePath, newFilePath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
