package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
			WebsiteURL: util.SanitizeFileName(websiteURL),
		}
	}

	_, err = config.DB.InsertURLs(linksArr)
	if err != nil {
		panic(err)
	}
}

func WaybackDownloadFile(config *config.Config, link data.Link) error {
	keepFullFilePath := (*config.Options.Options)[data.OptionsMap[data.OPTION_COPY_FULL_PATH]]

	fileName := util.SanitizeFileName(filepath.Base(link.Original))
	if len(fileName) > 50 {
		fileName = fileName[:50]
	}
	ext := strings.ToLower(strings.ReplaceAll(filepath.Ext(fileName), ".", ""))

	resp, err := http.Get(fmt.Sprintf(`https://web.archive.org/web/%soe_/%s`, link.Timestamp, link.Original))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("request failed")
	}

	defaultLocation := util.RemoveSlashFromString(link.Mimetype)
	if keepFullFilePath {
		u, err := url.Parse(link.Original)
		if err != nil {
			return err
		}
		urlPath := u.Path
		dirPath := path.Dir(urlPath)
		cleaned := strings.TrimPrefix(dirPath, "/")

		defaultLocation = cleaned
	}

	filePath := path.Join(data.MAIN_PATH, link.WebsiteURL, defaultLocation, fileName)

	dirPath := path.Dir(filePath)
	extSuffix := filepath.Ext(fileName)
	baseName := strings.TrimSuffix(fileName, extSuffix)
	counter := 1
	for {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			break
		}
		fileName = fmt.Sprintf("%s_%d%s", baseName, counter, extSuffix)
		filePath = path.Join(dirPath, fileName)
		counter++
	}

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
