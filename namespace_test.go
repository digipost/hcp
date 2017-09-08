package hcp

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Namespace - Create

func TestCreateNamespaceSuccess(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.WriteHeader(http.StatusOK)
	}))

	hcp := &HCP{URL: ts.URL}

	err := hcp.CreateNamespace(&Namespace{Name: "mynamespace", Description: "My description"})

	assert.Empty(t, err)
}

func TestCreateNamespaceFailure(t *testing.T) {

	const errorMsg = "Invalid namespace name"

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.Header().Set(xHcpErrorMessage, errorMsg)
		res.WriteHeader(http.StatusBadRequest)
	}))

	hcp := &HCP{URL: ts.URL}

	err := hcp.CreateNamespace(&Namespace{Name: "invalid namespace name"})

	assert.Contains(t, err.Error(), errorMsg)
}

func TestCreateNamespaceShouldTargetCorrectEndpoint(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/namespaces", r.URL.Path)
		assert.Equal(t, http.MethodPut, r.Method)
	}))
	defer ts.Close()

	hcp := &HCP{URL: ts.URL}
	hcp.CreateNamespace(&Namespace{})

}

// Namespace - Update

func TestUpdateNamespaceSuccess(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.WriteHeader(http.StatusOK)
	}))

	hcp := &HCP{URL: ts.URL}

	assert.Empty(t, hcp.UpdateNamespace(&Namespace{Name: "mynamespace", Description: "My description"}))
}

func TestUpdateNamespaceFailure(t *testing.T) {

	const errorMsg = "some random error"

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.Header().Set(xHcpErrorMessage, errorMsg)
		res.WriteHeader(http.StatusBadGateway)
	}))

	hcp := &HCP{URL: ts.URL}

	assert.Contains(t, hcp.UpdateNamespace(&Namespace{}).Error(), errorMsg)
}

func TestUpdateNamespaceShouldTargetCorrectEndpoint(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/namespaces/mynamespace", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
	}))
	defer ts.Close()

	hcp := &HCP{URL: ts.URL}
	hcp.UpdateNamespace(&Namespace{Name: "mynamespace"})

}

// Namespace - Read

func TestNamespaceSuccess(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.WriteHeader(http.StatusOK)
		data, _ := ioutil.ReadFile("xml/namespace/response/Namespace.xml")
		res.Write(data)
	}))

	hcp := &HCP{URL: ts.URL}

	namespace, err := hcp.Namespace("Accounts-Receivable")

	assert.Empty(t, err)
	assert.Equal(t, "Accounts-Receivable", namespace.Name)
	assert.Equal(t, SHA_256, namespace.HashScheme)
}

func TestNamespaceFailure(t *testing.T) {

	const errorMsg = "Some random error"

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.Header().Set(xHcpErrorMessage, errorMsg)
		res.WriteHeader(http.StatusInternalServerError)
	}))

	hcp := &HCP{URL: ts.URL}

	_, err := hcp.Namespace("Accounts-Receivable")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), errorMsg)

}

func TestNamespaceShouldTargetCorrectEndpoint(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/namespaces/Accounts-Receivable", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
	}))

	defer ts.Close()

	hcp := &HCP{URL: ts.URL}
	hcp.Namespace("Accounts-Receivable")

}

// Namespace - Delete

func TestDeleteNamespaceSuccess(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.WriteHeader(http.StatusOK)
	}))

	hcp := &HCP{URL: ts.URL}

	assert.Empty(t, hcp.DeleteNamespace("mynamespace"))
}

func TestDeleteNamespaceFailure(t *testing.T) {

	const errorMsg = "some random error"

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.Header().Set(xHcpErrorMessage, errorMsg)
		res.WriteHeader(http.StatusBadGateway)
	}))

	hcp := &HCP{URL: ts.URL}

	assert.Contains(t, hcp.DeleteNamespace("mynamespace").Error(), errorMsg)
}

func TestDeleteNamespaceShouldTargetCorrectEndpoint(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/namespaces/mynamespace", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)
	}))
	defer ts.Close()

	hcp := &HCP{URL: ts.URL}
	hcp.DeleteNamespace("mynamespace")

}

// Namespace - Exists

func TestNamespaceExistsSuccess(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {

		res.WriteHeader(http.StatusOK)
	}))

	hcp := &HCP{URL: ts.URL}

	exists, err := hcp.NamespaceExists("mynamespace")

	assert.Empty(t, err)
	assert.True(t, exists)
}

func TestNamespaceExistsFailure(t *testing.T) {

	const errorMsg = "some random error"

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		res.Header().Set(xHcpErrorMessage, errorMsg)
		res.WriteHeader(http.StatusBadGateway)
	}))

	hcp := &HCP{URL: ts.URL}

	_, err := hcp.NamespaceExists("mynamespace")
	assert.Contains(t, err.Error(), errorMsg)
}

func TestNamespaceExistsShouldTargetCorrectEndpoint(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/namespaces/mynamespace", r.URL.Path)
		assert.Equal(t, http.MethodHead, r.Method)
	}))
	defer ts.Close()

	hcp := &HCP{URL: ts.URL}
	hcp.NamespaceExists("mynamespace")

}
