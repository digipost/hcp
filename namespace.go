package hcp

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func (hcp *HCP) CreateNamespace(namespace *Namespace) error {

	if req, reqErr := hcp.createPutRequest("/namespaces", namespace); reqErr != nil {
		return reqErr
	} else {

		if res, doReqErr := hcp.getClient().Do(req); doReqErr != nil {
			return doReqErr
		} else {
			if res.StatusCode != http.StatusOK {
				return fmt.Errorf("Failed to create HCP namespace: %s. Status code: %d, HCP error message: %s",
					namespace.Name,
					res.StatusCode,
					hcpErrorMessage(res))
			} else {
				return nil
			}
		}

	}

}

func (hcp *HCP) UpdateNamespace(namespace *Namespace) error {

	if req, reqErr := hcp.createPostRequest("/namespaces/"+namespace.Name, namespace); reqErr != nil {
		return reqErr
	} else {

		if res, doReqErr := hcp.getClient().Do(req); doReqErr != nil {
			return doReqErr
		} else {
			if res.StatusCode != http.StatusOK {
				return fmt.Errorf("Failed to update HCP namespace: %s. Status code: %d, HCP error message: %s",
					namespace.Name,
					res.StatusCode,
					hcpErrorMessage(res))
			} else {
				return nil
			}

		}
	}

}

func (hcp *HCP) Namespace(name string) (*Namespace, error) {

	if req, reqErr := hcp.createGetRequest("/namespaces/" + name); reqErr != nil {
		return nil, reqErr
	} else {

		if res, doReqErr := hcp.getClient().Do(req); doReqErr != nil {
			return nil, doReqErr
		} else {
			if res.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("Failed to retrieve HCP namespace: %s. Status code: %d, HCP error message: %s",
					name,
					res.StatusCode,
					hcpErrorMessage(res))
			} else {

				if bytes, readErr := ioutil.ReadAll(res.Body); readErr != nil {
					return nil, readErr
				} else {
					namespace := &Namespace{}
					if unmarshalErr := unmarshal(bytes, namespace); unmarshalErr != nil {
						return nil, unmarshalErr
					} else {
						return namespace, nil
					}
				}

			}
		}

	}
}

func (hcp *HCP) DeleteNamespace(name string) error {
	if req, reqErr := hcp.createDeleteRequest("/namespaces/" + name); reqErr != nil {
		return reqErr
	} else {

		if res, doReqErr := hcp.getClient().Do(req); doReqErr != nil {
			return doReqErr
		} else {
			if res.StatusCode != http.StatusOK {
				return fmt.Errorf("Failed to delete HCP namespace: %s. Status code: %d, HCP error message: %s",
					name,
					res.StatusCode,
					hcpErrorMessage(res))
			} else {
				return nil
			}
		}

	}
}

func (hcp *HCP) NamespaceExists(name string) (bool, error) {
	if req, reqErr := hcp.createHeadRequest("/namespaces/" + name); reqErr != nil {
		return false, reqErr
	} else {

		if res, doReqErr := hcp.getClient().Do(req); doReqErr != nil {
			return false, doReqErr
		} else {
			if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound {
				return false, fmt.Errorf("Failed to check existence of HCP namespace: %s. Status code: %d, HCP error message: %s",
					name,
					res.StatusCode,
					hcpErrorMessage(res))
			}
			return res.StatusCode == http.StatusOK, nil
		}

	}
}
