package hcp

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

const xHcpErrorMessage = "X-HCP-ErrorMessage"

type HCP struct {
	// "https://finance.hcp.example.com:9090/mapi/tenants/finance
	URL      string
	Insecure bool
	Username string
	Password string
	client   *http.Client
}

func (hcp *HCP) authenticationToken() string {
	username := base64.StdEncoding.EncodeToString([]byte(hcp.Username))
	h := md5.New()
	io.WriteString(h, hcp.Password)
	password := fmt.Sprintf("%x", h.Sum(nil))
	return username + ":" + password
}

func (hcp *HCP) createHeadRequest(urlStr string) (*http.Request, error) {
	return hcp.createAuthorizedRequest(http.MethodHead, urlStr, nil)
}

func (hcp *HCP) createGetRequest(urlStr string) (*http.Request, error) {
	return hcp.createAuthorizedRequest(http.MethodGet, urlStr, nil)
}

func (hcp *HCP) createDeleteRequest(urlStr string) (*http.Request, error) {
	return hcp.createAuthorizedRequest(http.MethodDelete, urlStr, nil)
}

func (hcp *HCP) createPutRequest(urlStr string, xml Payload) (*http.Request, error) {
	return hcp.createRequestWithPayload(http.MethodPut, urlStr, xml)
}

func (hcp *HCP) createPostRequest(urlStr string, xml Payload) (*http.Request, error) {
	return hcp.createRequestWithPayload(http.MethodPost, urlStr, xml)
}

func (hcp *HCP) createRequestWithPayload(method string, urlStr string, xml Payload) (*http.Request, error) {

	if reader, error := xml.Reader(); error != nil {
		return nil, error
	} else {
		return hcp.createAuthorizedRequest(method, urlStr, reader)
	}

}

func (hcp *HCP) createAuthorizedRequest(method string, urlStr string, body io.Reader) (*http.Request, error) {

	if req, err := http.NewRequest(method, hcp.URL+urlStr, body); err != nil {
		return nil, err
	} else {
		req.Header.Set("Authorization", "HCP "+hcp.authenticationToken())
		req.Header.Set("Content-Type", "application/xml")
		return req, nil
	}

}

func hcpErrorMessage(response *http.Response) string {
	return response.Header.Get(xHcpErrorMessage)
}

func (hcp *HCP) getClient() *http.Client {
	if hcp.client == nil {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: hcp.Insecure},
		}
		hcp.client = &http.Client{Transport: tr}

	}
	return hcp.client
}
