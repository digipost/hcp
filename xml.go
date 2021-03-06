package hcp

import (
	"bytes"
	"encoding/xml"
)

func unmarshal(in []byte, out Payload) error {
	return xml.Unmarshal(in, out)
}

type Payload interface {
	Reader() (*bytes.Reader, error)
}

/************************
 *      UserAccount
 ************************/

// Roles
const (
	ADMINISTRATOR = "ADMINISTRATOR"
	COMPLIANCE    = "COMPLIANCE"
	MONITOR       = "MONITOR"
	SECURITY      = "SECURITY"
)

type UserAccount struct {
	XMLName                  xml.Name `xml:"userAccount"`
	Username                 string   `xml:"username,omitempty"`
	FullName                 string   `xml:"fullName,omitempty"`
	Description              string   `xml:"description,omitempty"`
	LocalAuthentication      bool     `xml:"localAuthentication,omitempty"`
	ForcePasswordChange      bool     `xml:"forcePasswordChange"`
	Enabled                  bool     `xml:"enabled"`
	Roles                    []string `xml:"roles>role,omitempty"`
	AllowNamespaceManagement bool     `xml:"allowNamespaceManagement,omitempty"`
}

func (uA *UserAccount) marshal() ([]byte, error) {
	return xml.Marshal(uA)
}

func (uA *UserAccount) Reader() (*bytes.Reader, error) {
	if xml, error := uA.marshal(); error != nil {
		return nil, error
	} else {
		return bytes.NewReader(xml), nil
	}
}

/************************
 *      Namespace
 ************************/

// Hashing schemes
const (
	SHA_512   = "SHA-512"
	SHA_384   = "SHA-384"
	SHA_256   = "SHA-256"
	MD5       = "MD5"
	SHA_1     = "SHA-1"
	RIPEMD160 = "RIPEMD160"
)

// Optimized for
const (
	ALL   = "ALL"
	CLOUD = "CLOUD"
)

// Hard quota units
const (
	GB = "GB"
	TB = "TB"
)

// ACL usage
const (
	NOT_ENABLED  = "NOT_ENABLED"
	ENFORCED     = "ENFORCED"
	NOT_ENFORCED = "NOT_ENFORCED"
)

// Owner type
const (
	LOCAL    = "LOCAL"
	EXTERNAL = "EXTERNAL"
)

type Namespace struct {
	XMLName                       xml.Name            `xml:"namespace"`
	Name                          string              `xml:"name,omitempty"`
	Description                   string              `xml:"description,omitempty"`
	HashScheme                    string              `xml:"hashScheme,omitempty"`
	EnterpriseMode                bool                `xml:"enterpriseMode"`
	HardQuota                     string              `xml:"hardQuota,omitempty"`
	SoftQuota                     int                 `xml:"softQuota,omitempty"`
	OptimizedFor                  string              `xml:"optimizedFor,omitempty"`
	VersioningSettings            *VersioningSettings `xml:"versioningSettings,omitempty"`
	SearchEnabled                 bool                `xml:"searchEnabled"`
	IndexingEnabled               bool                `xml:"indexingEnabled"`
	CustomMetadataIndexingEnabled bool                `xml:"customMetadataIndexingEnabled"`
	ReplicationEnabled            bool                `xml:"replicationEnabled"`
	ReadFromReplica               bool                `xml:"readFromReplica"`
	ServiceRemoteSystemRequests   bool                `xml:"serviceRemoteSystemRequests"`
	AclsUsage                     string              `xml:"aclsUsage,omitempty"`
	OwnerType                     string              `xml:"ownerType,omitempty"`
	Owner                         string              `xml:"owner,omitempty"`
	Tags                          []string            `xml:"tags>tag,omitempty"`
}

type VersioningSettings struct {
	Enabled   bool `xml:"enabled"`
	Prune     bool `xml:"prune"`
	PruneDays int  `xml:"pruneDays"`
}

func (ns *Namespace) marshal() ([]byte, error) {
	return xml.Marshal(ns)
}

func (ns *Namespace) Reader() (*bytes.Reader, error) {
	if xml, error := ns.marshal(); error != nil {
		return nil, error
	} else {
		return bytes.NewReader(xml), nil
	}
}

/************************
 * Namespace Protocol Http
 ************************/

type HttpProtocol struct {
	XMLName                    xml.Name    `xml:"httpProtocol"`
	Hs3Enabled                 bool        `xml:"hs3Enabled"`
	Hs3RequiresAuthentication  bool        `xml:"hs3RequiresAuthentication"`
	HttpEnabled                bool        `xml:"httpEnabled"`
	HttpsEnabled               bool        `xml:"httpsEnabled"`
	RestEnabled                bool        `xml:"restEnabled"`
	RestRequiresAuthentication bool        `xml:"restRequiresAuthentication"`
	IpSettings                 *IpSettings `xml:"ipSettings"`
}

type IpSettings struct {
	AllowAddresses     []string `xml:"allowAddresses>ipAddress"`
	DenyAddresses      []string `xml:"denyAddresses>ipAddress"`
	AllowIfInBothLists bool     `xml:"allowIfInBothLists"`
}

func (uA *HttpProtocol) marshal() ([]byte, error) {
	return xml.Marshal(uA)
}

func (uA *HttpProtocol) Reader() (*bytes.Reader, error) {
	if xml, error := uA.marshal(); error != nil {
		return nil, error
	} else {
		return bytes.NewReader(xml), nil
	}
}
