package shakespeare

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ShakespeareClient interface {
	Translate(text string) (string, error)
}

type shakespeareClient struct {
	baseURL string
	client  *http.Client
}

type shakespeareTranslationResponse struct {
	SuccessInfo struct {
		Total float64 `json:"total"`
	} `json:"success"`

	Contents struct {
		Translated string `json:"translated"`
	} `json:"contents"`
}

func (c *shakespeareClient) Translate(text string) (translated string, err error) {
	u := fmt.Sprintf("%s/translate/shakespeare", c.baseURL)

	resp, err := c.client.PostForm(u, url.Values{"text": []string{text}})
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("expected 200 received %d", resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var parsed shakespeareTranslationResponse
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return
	}

	if parsed.SuccessInfo.Total < 1 || len(parsed.Contents.Translated) == 0 {
		err = fmt.Errorf("translation did not succeed")
		return
	}

	translated = parsed.Contents.Translated
	return
}

func NewClient(baseURL string, client *http.Client) ShakespeareClient {
	return &shakespeareClient{
		baseURL: baseURL,
		client:  client,
	}
}
