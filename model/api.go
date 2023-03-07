package model

import (
	"crypto/tls"
	"github.com/wonderivan/logger"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type ProxmoxAPI struct{}
type JumpServerAPI struct{}

func (p ProxmoxAPI) Post(method string, msgText string) {
	url := pveConfig.ApiUrl + method
	_, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(msgText))
	if err != nil {
		logger.Info(err)
	}
}

func (p ProxmoxAPI) Get(method string) []byte {
	var body []byte
	var resp *http.Response
	url := pveConfig.ApiUrl + method
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", pveConfig.Token)
Loop1:
	for {
		resp, err = client.Do(req)
		if resp.StatusCode == 200 {
			break Loop1
		} else {
			logger.Info(err)
		}
		time.Sleep(5 * time.Second)
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		logger.Info(err)
	}
	err = resp.Body.Close()
	return body
}

func (p JumpServerAPI) Post(method string, msgText string) (int, []byte) {
	var body []byte
	var resp *http.Response
	var req *http.Request
	var err error
	url := jumpServerConfig.ApiUrl + method
	client := &http.Client{}
	gmtFmt := "Mon, 02 Jan 2006 15:04:05 GMT"
	req, err = http.NewRequest("POST", url, strings.NewReader(msgText))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Date", time.Now().Format(gmtFmt))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-JMS-ORG", "00000000-0000-0000-0000-000000000002")
	var sigAuth SigAuth
	err = sigAuth.Sign(req)
	if err != nil {
		logger.Info(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		logger.Info(err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Info(err)
	}
	return resp.StatusCode, body
}

func (p JumpServerAPI) Get(method string) []byte {
	var body []byte
	url := jumpServerConfig.ApiUrl + method
	gmtFmt := "Mon, 02 Jan 2006 15:04:05 GMT"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Date", time.Now().Format(gmtFmt))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-JMS-ORG", "00000000-0000-0000-0000-000000000002")
	if err != nil {
		log.Fatal(err)
	}
	var sigAuth SigAuth
	err = sigAuth.Sign(req)
	if err != nil {
		logger.Info(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Info(err)
	}
	//	defer resp.Body.Close()
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			logger.Info(err)
		}
	}()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		logger.Info(err)
	}
	return body
}
