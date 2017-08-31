package hcp

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

const x_HCP_ERROR_MESSAGE = "X-HCP-ErrorMessage"

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

/** User Account methods */

func (hcp *HCP) CreateUserAccount(userAccount *UserAccount, password string) error {
	return hcp.createOrUpdateUserAccount(true, userAccount, password)
}

func (hcp *HCP) UpdateUserAccount(userAccount *UserAccount, password string) error {
	return hcp.createOrUpdateUserAccount(false, userAccount, password)
}

func (hcp *HCP) createOrUpdateUserAccount(create bool, userAccount *UserAccount, password string) error {

	if err := validatePassword(password); err != nil {
		return err
	}

	var request *http.Request
	if create {
		if req, reqErr := hcp.createPutRequest("/userAccounts?password="+url.QueryEscape(password), userAccount); reqErr != nil {
			return reqErr
		} else {
			request = req
		}
	} else {
		if req, reqErr := hcp.createPostRequest("/userAccounts/"+userAccount.Username+"?password="+url.QueryEscape(password), userAccount); reqErr != nil {
			return reqErr
		} else {
			request = req
		}
	}

	if res, doReqErr := hcp.getClient().Do(request); doReqErr != nil {
		return doReqErr
	} else {
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("Failed to create HCP user account for username: %s. Status code: %d, HCP error message: %s",
				userAccount.Username,
				res.StatusCode,
				hcpErrorMessage(res))
		} else {
			return nil
		}

	}

}

func (hcp *HCP) UserAccount(username string) (*UserAccount, error) {

	if req, reqErr := hcp.createGetRequest("/userAccounts/" + username); reqErr != nil {
		return nil, reqErr
	} else {

		if res, doReqErr := hcp.getClient().Do(req); doReqErr != nil {
			return nil, doReqErr
		} else {
			if res.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("Failed to retrieve HCP user account for username: %s. Status code: %d, HCP error message: %s",
					username,
					res.StatusCode,
					hcpErrorMessage(res))
			} else {

				if bytes, readErr := ioutil.ReadAll(res.Body); readErr != nil {
					return nil, readErr
				} else {
					userAccount := &UserAccount{}
					if unmarshalErr := unmarshal(bytes, userAccount); unmarshalErr != nil {
						return nil, unmarshalErr
					} else {
						return userAccount, nil
					}
				}

			}
		}

	}
}

/** Namespace methods **/

func (hcp *HCP) CreateNamespace(namespace *Namespace) error {

	if req, reqErr := hcp.createPutRequest("/namespaces", namespace); reqErr != nil {
		return reqErr
	} else {

		if res, doReqErr := hcp.getClient().Do(req); doReqErr != nil {
			return doReqErr
		} else {
			if res.StatusCode != http.StatusOK {
				return fmt.Errorf("Failed to create HCP namespace: %s. Status code: %d, HCP error message: %s",
					namespace.Name,
					res.StatusCode,
					hcpErrorMessage(res))
			} else {
				return nil
			}
		}

	}

}

// VERBS

func (hcp *HCP) createGetRequest(urlStr string) (*http.Request, error) {
	return hcp.createRequest(http.MethodGet, urlStr, nil)
}

func (hcp *HCP) createDeleteRequest(urlStr string) (*http.Request, error) {
	return hcp.createRequest(http.MethodDelete, urlStr, nil)
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
		return hcp.createRequest(method, urlStr, reader)
	}

}

func (hcp *HCP) createRequest(method string, urlStr string, body io.Reader) (*http.Request, error) {

	if req, err := http.NewRequest(method, hcp.URL+urlStr, body); err != nil {
		return nil, err
	} else {
		req.Header.Set("Authorization", "HCP "+hcp.authenticationToken())
		req.Header.Set("Content-Type", "application/xml")
		return req, nil
	}

}

func hcpErrorMessage(response *http.Response) string {
	return response.Header.Get(x_HCP_ERROR_MESSAGE)
}

/**
* Passwords can be up to 64 characters long, are case sensitive,
* and can contain any valid UTF-8 characters, including white space.
* To be valid, a password must include at least one character from two
* of these three groups: alphabetic, numeric, and other.
 */
func validatePassword(password string) error {
	// it's 2017, man!
	if len(password) < 10 {
		return fmt.Errorf("Password too short")
	} else if len(password) > 64 {
		return fmt.Errorf("Password too long")
	} else if !(alphabetic(password) && numeric(password) ||
		alphabetic(password) && other(password) ||
		numeric(password) && other(password)) {
		return fmt.Errorf("Password must include at least one character from two of these three groups: " +
			"alphabetic, numeric, and other.")
	} else {
		return nil
	}
}

func alphabetic(str string) bool {
	matched, _ := regexp.Match("[[:alpha:]]", []byte(str))
	return matched
}

func numeric(str string) bool {
	matched, _ := regexp.Match("\\d", []byte(str))
	return matched
}

func other(str string) bool {
	return !numeric(str) && !alphabetic(str)
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
