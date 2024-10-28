package review

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"searchgpt/config"
)

const (
	configFile = "config.yml"

	queryLimit = 1000

	urlAccount = "/accounts/"
	urlChanges = "/changes/"
	urlNumber  = "&n="
	urlOption  = "&o="
	urlPrefix  = "/a"
	urlQuery   = "?q="
	urlStart   = "&start="
)

type Review interface {
	Account(string) (string, error)
	Query(string, int, int) ([]interface{}, error)
}

type Config struct {
	Gerrit []Gerrit `yaml:"gerrit"`
}

type Gerrit struct {
	Url  string `yaml:"url"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
}

type review struct {
	cfg *Config
}

func New() Review {
	var cfg Config

	file, err := config.ConfigFile.ReadFile(configFile)
	if err != nil {
		return nil
	}

	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil
	}

	return &review{
		cfg: &cfg,
	}
}

func (r *review) Account(name string) (string, error) {
	var buf []interface{}

	data, err := r.get(r.urlAccount(name, []string{"DETAILS"}))
	if err != nil {
		return "", nil
	}

	if err := json.Unmarshal(data[4:], &buf); err != nil {
		return "", nil
	}

	return buf[0].(map[string]interface{})["email"].(string), nil
}

func (r *review) Query(search string, start, count int) ([]interface{}, error) {
	helper := func(search string, start int) []interface{} {
		buf, err := r.get(r.urlQuery(search, []string{"CURRENT_REVISION", "DETAILED_ACCOUNTS"}, start))
		if err != nil {
			return nil
		}
		ret, err := r.unmarshalList(buf)
		if err != nil {
			return nil
		}
		return ret
	}

	buf := helper(search, start)

	if len(buf) == 0 {
		return []interface{}{}, nil
	}

	if len(buf) >= count {
		return buf[:count], nil
	}

	more, ok := buf[len(buf)-1].(map[string]interface{})["_more_changes"].(bool)
	if !ok {
		more = false
	}

	if !more {
		return buf, nil
	}

	if b, err := r.Query(search, start+len(buf), count-len(buf)); err == nil {
		buf = append(buf, b...)
	}

	return buf, nil
}

func (r *review) unmarshalList(data []byte) ([]interface{}, error) {
	var buf []interface{}

	if err := json.Unmarshal(data[4:], &buf); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}

	if len(buf) == 0 {
		return nil, errors.New("failed to match")
	}

	return buf, nil
}

func (r *review) urlAccount(name string, option []string) string {
	account := urlQuery + url.PathEscape(name) + urlOption + strings.Join(option, urlOption) +
		urlNumber + strconv.Itoa(1)

	buf := r.cfg.Gerrit[0].Url + urlAccount + account
	if r.cfg.Gerrit[0].User != "" && r.cfg.Gerrit[0].Pass != "" {
		buf = r.cfg.Gerrit[0].Url + urlPrefix + urlAccount + account
	}

	return buf
}

func (r *review) urlQuery(search string, option []string, start int) string {
	query := urlQuery + url.PathEscape(search) +
		urlOption + strings.Join(option, urlOption) +
		urlStart + strconv.Itoa(start) +
		urlNumber + strconv.Itoa(queryLimit)

	buf := r.cfg.Gerrit[0].Url + urlChanges + query
	if r.cfg.Gerrit[0].User != "" && r.cfg.Gerrit[0].Pass != "" {
		buf = r.cfg.Gerrit[0].Url + urlPrefix + urlChanges + query
	}

	return buf
}

func (r *review) get(_url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, _url, http.NoBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request")
	}

	if r.cfg.Gerrit[0].User != "" && r.cfg.Gerrit[0].Pass != "" {
		req.SetBasicAuth(r.cfg.Gerrit[0].User, r.cfg.Gerrit[0].Pass)
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to do")
	}

	defer func() {
		_ = rsp.Body.Close()
	}()

	if rsp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid status")
	}

	data, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read")
	}

	return data, nil
}
