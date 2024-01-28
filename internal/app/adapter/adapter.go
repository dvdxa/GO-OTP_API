package adapter

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"otp_api/conf"
	"otp_api/internal/app/pkg/logger"
	"time"
)

type Adapter struct {
	log        *logger.Logger
	cfg        *conf.Config
	httpClient *http.Client
}

func NewAdapter(log *logger.Logger, cfg *conf.Config) *Adapter {
	return &Adapter{
		log: log,
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.Adapter.Timeout),
		},
	}
}

func (a *Adapter) SendOTPToClient(account string, otp string) error {

	text := a.cfg.Adapter.Text + otp

	params := url.Values{}
	params.Add("login", a.cfg.Adapter.Login)
	params.Add("account", account)
	params.Add("text", text)

	fullURL := fmt.Sprintf("%s?%s", a.cfg.Adapter.URL, params.Encode())

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		a.log.Error(err)
		return err
	}

	client := http.Client{
		Timeout: time.Duration(a.cfg.Adapter.Timeout),
	}

	a.log.Infof("[URL] %s\n", fullURL)

	resp, err := client.Do(req)
	if err != nil {
		a.log.Error(err)
		if ok := err.(net.Error).Timeout(); ok {
			return errors.New("timeout")
		}
		return err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		a.log.Error(err)
		return err
	}

	a.log.Infof("[ADAPTER RESPONSE BODY] %v\n", string(b))
	return nil
}
