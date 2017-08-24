package hcp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthenticationToken(t *testing.T) {
	hcp := HCP{Username: "lgreen", Password: "p4ssw0rd"}
	assert.Equal(t, "bGdyZWVu:2a9d119df47ff993b662a8ef36f9ea20", hcp.authenticationToken())

}

/*
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

	if created, err := hcp.CreateUserAccount(userAccount, ""); err != nil {
		panic(err)
	} else {
		if !created {

			t.Fail()
		}
	}

}
*/
