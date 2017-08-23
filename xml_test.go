package hcp

import (
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestShoulBeAbleToUnmarshalUserAccountXML(t *testing.T) {

	data, _ := ioutil.ReadFile("xml/userAccount_CreateUser.xml")

	userAccount := UserAccount{}

	xml.Unmarshal(data, &userAccount)

	assert.Equal(t, "mwhite", userAccount.Username)
	assert.Equal(t, "Morgan White", userAccount.FullName)
	assert.Equal(t, "Compliance officer.", userAccount.Description)
	assert.Equal(t, true, userAccount.LocalAuthentication)
	assert.Equal(t, true, userAccount.ForcePasswordChange)
	assert.Equal(t, true, userAccount.Enabled)
	assert.Equal(t, 2, len(userAccount.Roles))
	assert.Equal(t, COMPLIANCE, userAccount.Roles[0])
	assert.Equal(t, MONITOR, userAccount.Roles[1])
	assert.Equal(t, false, userAccount.AllowNamespaceManagement)

}

func TestShouldOmitIfEmptyOnUserAccountXML(t *testing.T) {

	userAccount := &UserAccount{Roles: []string{COMPLIANCE, MONITOR, ADMINISTRATOR}}
	xml, _ := userAccount.marshal()
	assert.Equal(t, "<userAccount><roles><role>COMPLIANCE</role><role>MONITOR</role><role>ADMINISTRATOR</role></roles></userAccount>", string(xml))

}

func TestShouldBeAbleToNamespaceXML(t *testing.T) {

	data, _ := ioutil.ReadFile("xml/namespace_CreateNamespace.xml")

	ns := Namespace{}

	xml.Unmarshal(data, &ns)

	assert.Equal(t, "Accounts-Receivable", ns.Name)
	assert.Equal(t, "Created for the Finance department at Example Company by Lee Green on 2/9/2012.", ns.Description)
	assert.Equal(t, SHA_256, ns.HashScheme)
	assert.True(t, ns.EnterpriseMode)
	assert.Equal(t, "50 GB", ns.HardQuota)
	assert.Equal(t, 75, ns.SoftQuota)
	assert.Equal(t, ALL, ns.OptimizedFor)
	assert.True(t, ns.VersioningSettings.Enabled)
	assert.True(t, ns.VersioningSettings.Prune)
	assert.Equal(t, 10, ns.VersioningSettings.PruneDays)
	assert.True(t, ns.SearchEnabled)
	assert.True(t, ns.IndexingEnabled)
	assert.True(t, ns.CustomMetadataIndexingEnabled)
	assert.True(t, ns.ReplicationEnabled)
	assert.True(t, ns.ReadFromReplica)
	assert.True(t, ns.ServiceRemoteSystemRequests)
	assert.Equal(t, 2, len(ns.Tags))
	assert.Equal(t, "Billing", ns.Tags[0])
	assert.Equal(t, "lgreen", ns.Tags[1])

}
