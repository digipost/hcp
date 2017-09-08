package hcp

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const validPassword = "12345asdfasdf"

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
		res.Header().Set(xHcpErrorMessage, errorMsg)
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
		res.Header().Set(xHcpErrorMessage, errorMsg)
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
	}))

	hcp := &HCP{URL: ts.URL}

	err := hcp.DeleteUserAccount("mwhite")

	assert.Empty(t, err)
}

func TestDeleteUserAccountFailure(t *testing.T) {

	const errorMessage = "some error message"

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.Header().Set(xHcpErrorMessage, errorMessage)
		res.WriteHeader(http.StatusBadRequest)

	}))

	hcp := &HCP{URL: ts.URL}

	err := hcp.DeleteUserAccount("mwhite")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), errorMessage)
}

func TestDeleteUserAccountShouldTargetCorrectEndpoint(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.WriteHeader(http.StatusOK)
		assert.Contains(t, r.URL.Path, "mwhite")
		assert.Equal(t, http.MethodDelete, r.Method)
	}))

	hcp := &HCP{URL: ts.URL}

	hcp.DeleteUserAccount("mwhite")

}

// User Account - Exists

func TestUserAccountExistsSuccessFound(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.WriteHeader(http.StatusOK)
	}))

	hcp := &HCP{URL: ts.URL}

	exists, err := hcp.UserAccountExists("mwhite")
	assert.True(t, exists)
	assert.Empty(t, err)
}

func TestUserAccountExistsSuccessNotFound(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.WriteHeader(http.StatusNotFound)
	}))

	hcp := &HCP{URL: ts.URL}

	exists, err := hcp.UserAccountExists("mwhite")
	assert.False(t, exists)
	assert.Empty(t, err)
}

func TestUserAccountExistsFailure(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.WriteHeader(http.StatusBadGateway)
	}))

	hcp := &HCP{URL: ts.URL}

	exists, err := hcp.UserAccountExists("mwhite")
	assert.False(t, exists)
	assert.Error(t, err)
}

func TestUserAccountExistsShouldTargetCorrectEndpoint(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "white")
		assert.Equal(t, r.Method, http.MethodHead)
		res.WriteHeader(http.StatusOK)
	}))

	hcp := &HCP{URL: ts.URL}

	hcp.UserAccountExists("mwhite")

}
