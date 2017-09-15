package hcp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

func (hcp *HCP) CreateUserAccount(userAccount *UserAccount, password string) error {

	if req, reqErr := hcp.createPutRequest("/userAccounts", userAccount); reqErr != nil {
		return reqErr
	} else {

		if err := validatePassword(password); err != nil {
			return err
		}
		query := req.URL.Query()
		query.Set("password", password)
		req.URL.RawQuery = query.Encode()

		if res, doReqErr := hcp.getClient().Do(req); doReqErr != nil {
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

}

func (hcp *HCP) UpdateUserAccountExceptPassword(userAccount *UserAccount) error {

	if req, reqErr := hcp.createPostRequest("/userAccounts/"+userAccount.Username, userAccount); reqErr != nil {
		return reqErr
	} else {

		if res, doReqErr := hcp.getClient().Do(req); doReqErr != nil {
			return doReqErr
		} else {
			if res.StatusCode != http.StatusOK {
				return fmt.Errorf("Failed to update HCP user account except password for username: %s. Status code: %d, HCP error message: %s",
					userAccount.Username,
					res.StatusCode,
					hcpErrorMessage(res))
			} else {
				return nil
			}

		}
	}

}

func (hcp *HCP) UpdateUserAccount(userAccount *UserAccount, password string) error {

	if req, reqErr := hcp.createPostRequest("/userAccounts/"+userAccount.Username, userAccount); reqErr != nil {
		return reqErr
	} else {

		if err := validatePassword(password); err != nil {
			return err
		}
		query := req.URL.Query()
		query.Set("password", password)
		req.URL.RawQuery = query.Encode()

		if res, doReqErr := hcp.getClient().Do(req); doReqErr != nil {
			return doReqErr
		} else {
			if res.StatusCode != http.StatusOK {
				return fmt.Errorf("Failed to update HCP user account for username: %s. Status code: %d, HCP error message: %s",
					userAccount.Username,
					res.StatusCode,
					hcpErrorMessage(res))
			} else {
				return nil
			}

		}
	}

}

func (hcp *HCP) ReadUserAccount(username string) (*UserAccount, error) {

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

func (hcp *HCP) DeleteUserAccount(username string) error {
	if req, reqErr := hcp.createDeleteRequest("/userAccounts/" + username); reqErr != nil {
		return reqErr
	} else {

		if res, doReqErr := hcp.getClient().Do(req); doReqErr != nil {
			return doReqErr
		} else {
			if res.StatusCode != http.StatusOK {
				return fmt.Errorf("Failed to delete HCP user account for username: %s. Status code: %d, HCP error message: %s",
					username,
					res.StatusCode,
					hcpErrorMessage(res))
			} else {
				return nil
			}
		}

	}
}

func (hcp *HCP) UserAccountExists(username string) (bool, error) {
	if req, reqErr := hcp.createHeadRequest("/userAccounts/" + username); reqErr != nil {
		return false, reqErr
	} else {

		if res, doReqErr := hcp.getClient().Do(req); doReqErr != nil {
			return false, doReqErr
		} else {
			if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound {
				return false, fmt.Errorf("Failed to check existence of HCP user account for username: %s. Status code: %d, HCP error message: %s",
					username,
					res.StatusCode,
					hcpErrorMessage(res))
			}
			return res.StatusCode == http.StatusOK, nil
		}

	}
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
