package hcp

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type HCP struct {
	// "https://finance.hcp.example.com:9090/mapi/tenants/finance
	URL      string
	Insecure bool
	Username string
	Password string
}

func (hcp *HCP) authenticationToken() string {
	username := base64.StdEncoding.EncodeToString([]byte(hcp.Username))
	h := md5.New()
	io.WriteString(h, hcp.Password)
	password := fmt.Sprintf("%x", h.Sum(nil))
	return username + ":" + password
}

func (hcp *HCP) CreateUserAccount(userAccount *UserAccount, password string) (bool, error) {

	if req, reqErr := hcp.createRequest(http.MethodPut, "/userAccounts?password="+url.QueryEscape(password), userAccount); reqErr != nil {
		return false, reqErr
	} else {

		if res, doReqErr := hcp.getClient().Do(req); doReqErr != nil {
			return false, doReqErr
		} else {
			return res.StatusCode == 200, nil
		}

	}

}

func (hcp *HCP) CreateNamespace(namespace *Namespace) (bool, error) {

	if req, reqErr := hcp.createRequest(http.MethodPut, "/namespaces",
		namespace); reqErr != nil {
		return false, reqErr
	} else {

		if res, doReqErr := hcp.getClient().Do(req); doReqErr != nil {
			return false, doReqErr
		} else {
			return res.StatusCode == 200, nil
		}

	}

}

func (hcp *HCP) createRequest(method string, urlStr string, xml Request) (*http.Request, error) {

	if reader, error := xml.Reader(); error != nil {
		return nil, error
	} else {

		if request, newRequestError := http.NewRequest(method, hcp.URL+urlStr, reader); newRequestError != nil {
			return nil, newRequestError
		} else {

			request.Header.Set("Authorization", "HCP "+hcp.authenticationToken())
			request.Header.Set("Content-Type", "application/xml")
			return request, nil
		}
	}

}

func (hcp *HCP) getClient() *http.Client {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: hcp.Insecure},
	}
	return &http.Client{Transport: tr}

}
