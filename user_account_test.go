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

	assert.Empty(t, hcp.CreateUserAccount(&UserAccount{Username: "mwhite"}, validPassword))

}

func TestCreateUserAccountFailure(t *testing.T) {

	const errorMsg = "Some random error"

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.Header().Set(xHcpErrorMessage, errorMsg)
		res.WriteHeader(http.StatusBadRequest)
	}))

	defer ts.Close()

	hcp := &HCP{URL: ts.URL}

	assert.Contains(t, hcp.CreateUserAccount(&UserAccount{Username: "mwhite"}, validPassword).Error(), errorMsg)

}

func TestCreateUserAccountShouldTargetCorrectEndpoint(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/userAccounts", r.URL.Path)
		assert.Equal(t, http.MethodPut, r.Method)
	}))

	defer ts.Close()

	hcp := &HCP{URL: ts.URL}
	hcp.CreateUserAccount(&UserAccount{}, validPassword)

}

// User Account - Update

func TestUpdateUserAccountSuccess(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.WriteHeader(http.StatusOK)
	}))

	defer ts.Close()

	hcp := &HCP{URL: ts.URL}

	assert.Empty(t, hcp.UpdateUserAccount(&UserAccount{Username: "mwhite"}, validPassword))

}

func TestUpdateUserAccountFailure(t *testing.T) {

	const errorMsg = "some errror"
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.Header().Set(xHcpErrorMessage, errorMsg)
		res.WriteHeader(http.StatusBadGateway)
	}))

	defer ts.Close()

	hcp := &HCP{URL: ts.URL}

	err := hcp.UpdateUserAccount(&UserAccount{Username: "mwhite"}, validPassword)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), errorMsg)

}

func TestUpdateUserAccountShouldTargetCorrectEndpoint(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/userAccounts/mwhite", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
	}))
	defer ts.Close()

	hcp := &HCP{URL: ts.URL}
	hcp.UpdateUserAccount(&UserAccount{Username: "mwhite"}, validPassword)

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
		assert.Equal(t, "/userAccounts/mwhite", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
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
		assert.Equal(t, "/userAccounts/mwhite", r.URL.Path)
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
		assert.Equal(t, "/userAccounts/mwhite", r.URL.Path)
		assert.Equal(t, http.MethodHead, r.Method)
		res.WriteHeader(http.StatusOK)
	}))

	hcp := &HCP{URL: ts.URL}

	hcp.UserAccountExists("mwhite")

}
