package hcp

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
		res.Header().Set(xHcpErrorMessage, errorMsg)
		res.WriteHeader(http.StatusBadRequest)
	}))

	hcp := &HCP{URL: ts.URL}

	namespace := &Namespace{Name: "invalid namespace name"}

	err := hcp.CreateNamespace(namespace)

	assert.Contains(t, err.Error(), errorMsg)
}
