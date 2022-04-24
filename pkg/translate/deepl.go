package translate

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gopkg.in/ini.v1"
)

type Deepl struct {
	apiKey   string
	lang     string
	waitTime int // milli second
}

func NewDeepl(sec *ini.Section) *Deepl {
	return &Deepl{
		apiKey:   sec.Key("key").String(),
		lang:     sec.Key("lang").String(),
		waitTime: 100, // TODO: to be configurable
	}
}

const (
	endpoint = "https://api-free.deepl.com/v2/translate"
)

type translations struct {
	DetectedSourceLanguage string // `json:"detected_source_language"`
	Text                   string //`json:"text"`
}

type deeplRes struct {
	Translations []translations // `json:"translations"`
}

var lastRunTime = time.Now()

func (d Deepl) Run(text string) (string, error) {
	if lastRunTime.After(time.Now().Add(time.Duration(d.waitTime) * time.Millisecond)) {
		time.Sleep(time.Duration(d.waitTime))
	}
	lastRunTime = time.Now()

	// TODO: run with multi thread
	body := url.Values{}
	body.Add("auth_key", d.apiKey)
	body.Add("text", text)
	body.Add("target_lang", d.lang)

	res, err := http.Post(endpoint, "application/x-www-form-urlencoded", strings.NewReader(body.Encode()))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var bytes []byte
	if bytes, err = ioutil.ReadAll(res.Body); err != nil {
		return "", err
	}

	var dr deeplRes
	if err = json.Unmarshal(bytes, &dr); err != nil {
		return "", err
	}
	return dr.Translations[0].Text, nil
}
