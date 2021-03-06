package hcp

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestShoulBeAbleToUnmarshalUserAccountXML(t *testing.T) {

	data, _ := ioutil.ReadFile("xml/userAccount/request/CreateUserAccount.xml")

	userAccount := &UserAccount{}
	unmarshal(data, userAccount)

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

	xml, _ := (&UserAccount{Roles: []string{COMPLIANCE, MONITOR, ADMINISTRATOR}}).marshal()
	assert.Equal(t, "<userAccount><forcePasswordChange>false</forcePasswordChange><enabled>false</enabled><roles><role>COMPLIANCE</role><role>MONITOR</role><role>ADMINISTRATOR</role></roles></userAccount>", string(xml))

}

func TestShouldBeAbleToUnmarshalNamespaceXML(t *testing.T) {

	data, _ := ioutil.ReadFile("xml/namespace/request/CreateNamespace.xml")

	ns := &Namespace{}
	unmarshal(data, ns)

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

func TestShouldOmitVersioningSettingOnNamespaceXML(t *testing.T) {
	xml, _ := (&Namespace{Name: "name"}).marshal()
	assert.Equal(t,
		"<namespace>"+
			"<name>name</name>"+
			"<enterpriseMode>false</enterpriseMode>"+
			"<searchEnabled>false</searchEnabled>"+
			"<indexingEnabled>false</indexingEnabled>"+
			"<customMetadataIndexingEnabled>false</customMetadataIndexingEnabled>"+
			"<replicationEnabled>false</replicationEnabled>"+
			"<readFromReplica>false</readFromReplica>"+
			"<serviceRemoteSystemRequests>false</serviceRemoteSystemRequests>"+
			"<tags></tags>"+
			"</namespace>",
		string(xml))

}
