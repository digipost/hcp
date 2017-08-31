package hcp

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const validPassword = "12345asdfasdf"

func TestAuthenticationToken(t *testing.T) {
	hcp := HCP{Username: "lgreen", Password: "p4ssw0rd"}
	assert.Equal(t, "bGdyZWVu:2a9d119df47ff993b662a8ef36f9ea20", hcp.authenticationToken())
}

func TestPassword(t *testing.T) {
	assert.True(t, alphabetic("asdf"))
	assert.False(t, alphabetic("123"))
	assert.True(t, numeric("123"))
	assert.False(t, numeric("asdf"))
	assert.True(t, other("!$%&"))
	assert.False(t, other("asdf"))
	assert.False(t, other("123"))
	assert.Error(t, validatePassword("1234567890"))
	assert.Error(t, validatePassword("aaaaaaaaaa"))
	assert.Error(t, validatePassword("!!!!!!!!!!"))
	assert.Empty(t, validatePassword("12345asdfa"))
}

/** User Account methods */

// User Account - Create

func TestCreateUserAccountSuccess(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.WriteHeader(http.StatusOK)
	}))

	defer ts.Close()

	hcp := &HCP{URL: ts.URL}
	userAccount := &UserAccount{Username: "mwhite"}

	assert.Empty(t, hcp.CreateUserAccount(userAccount, validPassword))

}

func TestCreateUserAccountFailure(t *testing.T) {

	const errorMsg = "Some random error"

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.Header().Set(x_HCP_ERROR_MESSAGE, errorMsg)
		res.WriteHeader(http.StatusBadRequest)
	}))

	defer ts.Close()

	hcp := &HCP{URL: ts.URL}
	userAccount := &UserAccount{Username: "mwhite"}

	assert.Contains(t, hcp.CreateUserAccount(userAccount, validPassword).Error(), errorMsg)

}

func TestCreateUserAccountShouldTargetCorrectEndpoint(t *testing.T) {

	userAccount := &UserAccount{Username: "username"}

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		assert.NotContains(t, r.URL.Path, userAccount.Username)
		assert.Equal(t, r.Method, http.MethodPut)
	}))
	defer ts.Close()

	hcp := &HCP{URL: ts.URL}
	hcp.CreateUserAccount(userAccount, validPassword)

}

// User Account - Update

func TestUpdateUserAccountSuccess(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.WriteHeader(http.StatusOK)
	}))

	defer ts.Close()

	hcp := &HCP{URL: ts.URL}
	userAccount := &UserAccount{Username: "mwhite"}

	assert.Empty(t, hcp.UpdateUserAccount(userAccount, validPassword))

}

func TestUpdateUserAccountShouldTargetCorrectEndpoint(t *testing.T) {

	userAccount := &UserAccount{Username: "username"}

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, userAccount.Username)
		assert.Equal(t, r.Method, http.MethodPost)
	}))
	defer ts.Close()

	hcp := &HCP{URL: ts.URL}
	hcp.UpdateUserAccount(userAccount, validPassword)

}

// User Account - Read

func TestUserAccountSuccess(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.WriteHeader(http.StatusOK)
		data, _ := ioutil.ReadFile("xml/userAccount/response/UserAccount.xml")
		res.Write(data)
	}))

	hcp := &HCP{URL: ts.URL}

	uA, err := hcp.UserAccount("mwhite")

	assert.Empty(t, err)
	assert.Equal(t, "mwhite", uA.Username)
}

func TestUserAccountFailure(t *testing.T) {

	const errorMsg = "Some random error"

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.Header().Set(x_HCP_ERROR_MESSAGE, errorMsg)
		res.WriteHeader(http.StatusInternalServerError)
	}))

	hcp := &HCP{URL: ts.URL}

	_, err := hcp.UserAccount("mwhite")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), errorMsg)

}

func TestUserAccountShouldTargetCorrectEndpoint(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "mwhite")
		assert.Equal(t, r.Method, http.MethodGet)
	}))
	defer ts.Close()

	hcp := &HCP{URL: ts.URL}
	hcp.UserAccount("mwhite")

}

// User Account - Delete

func TestDeleteUserAccountSuccess(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.WriteHeader(http.StatusOK)
		assert.Contains(t, r.URL.Path, "mwhite")
		assert.Equal(t, http.MethodDelete, r.Method)
	}))

	hcp := &HCP{URL: ts.URL}

	err := hcp.DeleteUserAccount("mwhite")

	assert.Empty(t, err)
}

/** Namespace methods */

// Namespace - Create

func TestNamespaceSuccess(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.WriteHeader(http.StatusOK)
	}))

	hcp := &HCP{URL: ts.URL}

	namespace := &Namespace{Name: "mynamespace", Description: "My description"}

	err := hcp.CreateNamespace(namespace)

	assert.Empty(t, err)
}

func TestNamespaceFailure(t *testing.T) {

	const errorMsg = "Invalid namespace name"

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.Header().Set(x_HCP_ERROR_MESSAGE, errorMsg)
		res.WriteHeader(http.StatusBadRequest)
	}))

	hcp := &HCP{URL: ts.URL}

	namespace := &Namespace{Name: "invalid namespace name"}

	err := hcp.CreateNamespace(namespace)

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
