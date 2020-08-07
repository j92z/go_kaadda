package api

import (
	"github/j92z/go_kaadda/setting"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/mozillazg/request"
	"net/http"
)

type RequestUserInfoResponseBody struct {
	ID     string `json:"Id"`
	Name   string `json:"Display"`
	Parent struct {
		ID   string `json:"Id"`
		Name string `json:"Name"`
	} `json:"Parent"`
}

func RequestUserInfo(submitter string) ([]RequestUserInfoResponseBody, error) {

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true}}}
	req := request.NewRequest(client)
	resp, err := req.Get(setting.StructureServer + "/Structure?Method=ById&Id=" + submitter)
	if err != nil {
		return nil, err
	}
	res, err := resp.Text()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var userInfo []RequestUserInfoResponseBody
	if err := json.Unmarshal([]byte(res), &userInfo); err != nil {
		return nil, err
	}
	if len(userInfo) != 1 {
		return nil, errors.New("用户数据有误")
	}
	return userInfo, nil
}
