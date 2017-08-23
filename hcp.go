package hcp

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
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

	if xml, marshalError := userAccount.marshal(); marshalError != nil {
		return false, marshalError
	} else {

		url := hcp.URL + "/userAccounts?password=" + password

		if request, newRequestError := http.NewRequest(http.MethodPut, url, bytes.NewReader(xml)); newRequestError != nil {
			return false, newRequestError
		} else {

			request.Header.Set("Authorization", "HCP "+hcp.authenticationToken())
			request.Header.Set("Content-Type", "application/xml")

			client := createClient(true)

			if response, doRequestErr := client.Do(request); doRequestErr != nil {
				return false, doRequestErr
			} else {
				return response.StatusCode == 200, nil
			}

		}

	}

}

func createClient(insecure bool) http.Client {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}
	return http.Client{Transport: tr}

}
