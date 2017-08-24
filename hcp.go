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

func (hcp *HCP) CreateUserAccount(userAccount *UserAccount, password string) error {

	if req, reqErr := hcp.createRequest(http.MethodPut, "/userAccounts?password="+url.QueryEscape(password), userAccount); reqErr != nil {
		return reqErr
	} else {

		if res, doReqErr := hcp.getClient().Do(req); doReqErr != nil {
			return doReqErr
		} else {
			if res.StatusCode != http.StatusOK {
				return fmt.Errorf("Failed to create user account for username: %s. Status code: %s, HCP error message: %s",
					userAccount.Username,
					res.StatusCode,
					hcpErrorMessage(res))
			} else {
				return nil
			}

		}

	}

}

func (hcp *HCP) CreateNamespace(namespace *Namespace) error {

	if req, reqErr := hcp.createRequest(http.MethodPut, "/namespaces", namespace); reqErr != nil {
		return reqErr
	} else {

		if res, doReqErr := hcp.getClient().Do(req); doReqErr != nil {
			return doReqErr
		} else {
			if res.StatusCode != http.StatusOK {
				return fmt.Errorf("Failed to create namespace: %s. Status code: %s, HCP error message: %s",
					namespace.Name,
					res.StatusCode,
					hcpErrorMessage(res))
			} else {
				return nil
			}
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

func hcpErrorMessage(response *http.Response) string {
	return response.Header.Get("X-HCP-ErrorMessage")
}

func (hcp *HCP) getClient() *http.Client {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: hcp.Insecure},
	}
	return &http.Client{Transport: tr}

}
