package wxpay

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type PayClient struct {

	httpClient *http.Client
}

func NewPayClient(httpClient *http.Client) *PayClient  {

	if httpClient == nil {

		httpClient = http.DefaultClient
		httpClient.Timeout = time.Second * 5
	}

	return &PayClient{
		httpClient:httpClient,
	}
}

func (pc *PayClient) Login(jscode string) (*PayData, error)  {

	loginUrl := "https://api.weixin.qq.com/sns/jscode2session?appid=" + url.QueryEscape(payConfig.AppId()) +
		"&secret=" + url.QueryEscape(payConfig.AppSecret()) +
		"&js_code=" + url.QueryEscape(jscode) +
		"&grant_type=authorization_code"


	httpResp, err := pc.httpClient.Get(loginUrl)
	if err != nil {
		return  nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http.Status: %s", httpResp.Status)
	}

	respData := NewPayData()
	err = respData.FromJson(httpResp.Body)
	if err != nil {
		return nil, err
	}

	return respData, nil
}

func (pc *PayClient) PostXML(url string, pdata *PayData) (*PayData, error)  {

	reqSignType := pdata.Get("sign")
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(pdata.ToXml()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	req = req.WithContext(ctx)
	httpResp, err := pc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http.Status: %s", httpResp.Status)
	}

	respData := NewPayData()
	err = respData.FromXml(httpResp.Body)
	if err != nil {
		return nil, err
	}

	err = respData.CheckSign(payConfig.ApiKey(), reqSignType)
	if err != nil {
		return nil, err
	}

	return respData, nil
}