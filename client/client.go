package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/XANi/rapporter/db"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var HTTPClient = http.Client{Timeout: time.Second * 30}

func SendReport(url string, report db.Report) error {
	if err := report.Validate(); err != nil {
		return err
	}
	if !strings.Contains(url, "/api/v1") {
		url = strings.TrimRight(url, "/")
		url = url + "/api/v1"
	}
	jsonData, err := json.Marshal(&report)

	if err != nil {
		return err
	}

	resp, err := http.Post("https://httpbin.org/post", "application/json",
		bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("err: [%d %s] %s", resp.StatusCode, resp.Status, string(body))
	}
	return nil
}
