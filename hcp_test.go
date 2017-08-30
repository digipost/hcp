package hcp

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticationToken(t *testing.T) {
	hcp := HCP{Username: "lgreen", Password: "p4ssw0rd"}
	assert.Equal(t, "bGdyZWVu:2a9d119df47ff993b662a8ef36f9ea20", hcp.authenticationToken())

}

/** User Account methods */

// User Account - Create

func TestCreateUserAccountSuccess(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.WriteHeader(http.StatusOK)

	}))

	defer ts.Close()

	hcp := &HCP{
		URL: ts.URL,
	}

	userAccount := &UserAccount{
		Username:                 "mwhite",
		FullName:                 "Morgan White",
		Description:              "Compliance officer.",
		LocalAuthentication:      true,
		ForcePasswordChange:      true,
		Enabled:                  true,
		AllowNamespaceManagement: false,
	}

	assert.Empty(t, hcp.CreateUserAccount(userAccount, "whatever"))

}

func TestCreateUserAccountFails(t *testing.T) {

	const errorMsg = "Some random error"

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.Header().Set(x_HCP_ERROR_MESSAGE, errorMsg)
		res.WriteHeader(http.StatusBadRequest)

	}))

	defer ts.Close()

	hcp := &HCP{
		URL: ts.URL,
	}

	userAccount := &UserAccount{
		Username:                 "mwhite",
		FullName:                 "Morgan White",
		Description:              "Compliance officer.",
		LocalAuthentication:      true,
		ForcePasswordChange:      true,
		Enabled:                  true,
		AllowNamespaceManagement: false,
	}

	assert.Contains(t, hcp.CreateUserAccount(userAccount, "whatever").Error(), errorMsg)

}

// User Account - Read

func TestUserAccountFails(t *testing.T) {

	const errorMsg = "Some random error"

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.Header().Set(x_HCP_ERROR_MESSAGE, errorMsg)
		res.WriteHeader(http.StatusInternalServerError)

	}))

	hcp := &HCP{
		URL: ts.URL,
	}

	_, err := hcp.UserAccount("mwhite")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), errorMsg)

}

/*

// Use for live testing
func TestCreateUserAccounts(t *testing.T) {

	hcp := &HCP{
		URL:      "",
		Insecure: true,
		Username: "",
		Password: "",
	}

	userAccount := &UserAccount{
		Username:                 "mwhite",
		FullName:                 "Morgan White",
		Description:              "Compliance officer.",
		LocalAuthentication:      true,
		ForcePasswordChange:      true,
		Enabled:                  true,
		AllowNamespaceManagement: false,
	}

	if err := hcp.CreateUserAccount(userAccount, ""); err != nil {
		panic(err)
	} else {
		if !created {

			t.Fail()
		}
	}

}
*/
