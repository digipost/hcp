package hcp

import (
	"bytes"
	"encoding/xml"
)

/**
 * UserAccount
 */

// Roles
const (
	ADMINISTRATOR = "ADMINISTRATOR"
	COMPLIANCE    = "COMPLIANCE"
	MONITOR       = "MONITOR"
	SECURITY      = "SECURITY"
)

type Request interface {
	Reader() (*bytes.Reader, error)
}

type UserAccount struct {
	XMLName                  xml.Name `xml:"userAccount"`
	Username                 string   `xml:"username,omitempty"`
	FullName                 string   `xml:"fullName,omitempty"`
	Description              string   `xml:"description,omitempty"`
	LocalAuthentication      bool     `xml:"localAuthentication,omitempty"`
	ForcePasswordChange      bool     `xml:"forcePasswordChange,omitempty"`
	Enabled                  bool     `xml:"enabled,omitempty"`
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

/**
 * Namespace
 */

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
	ALL = "ALL"
)

// Hard quota units
const (
	GB = "GB"
	TB = "TB"
)

type Namespace struct {
	XMLName                       xml.Name           `xml:"namespace"`
	Name                          string             `xml:"name,omitempty"`
	Description                   string             `xml:"description,omitempty"`
	HashScheme                    string             `xml:"hashScheme,omitempty"`
	EnterpriseMode                bool               `xml:"enterpriseMode,omitempty"`
	HardQuota                     string             `xml:"hardQuota,omitempty"`
	SoftQuota                     int                `xml:"softQuota,omitempty"`
	OptimizedFor                  string             `xml:"optimizedFor,omitempty"`
	VersioningSettings            VersioningSettings `xml:"versioningSettings,omitempty"`
	SearchEnabled                 bool               `xml:"searchEnabled,omitempty"`
	IndexingEnabled               bool               `xml:"indexingEnabled,omitempty"`
	CustomMetadataIndexingEnabled bool               `xml:"customMetadataIndexingEnabled,omitempty"`
	ReplicationEnabled            bool               `xml:"replicationEnabled,omitempty"`
	ReadFromReplica               bool               `xml:"readFromReplica,omitempty"`
	ServiceRemoteSystemRequests   bool               `xml:"serviceRemoteSystemRequests,omitempty"`
	Tags                          []string           `xml:"tags>tag,omitempty"`
}

type VersioningSettings struct {
	Enabled   bool `xml:"enabled,omitempty"`
	Prune     bool `xml:"prune,omitempty"`
	PruneDays int  `xml:"pruneDays,omitempty"`
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
