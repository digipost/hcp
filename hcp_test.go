package hcp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthenticationToken(t *testing.T) {
	hcp := HCP{Username: "lgreen", Password: "p4ssw0rd"}
	assert.Equal(t, "bGdyZWVu:2a9d119df47ff993b662a8ef36f9ea20", hcp.authenticationToken())
}
