/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cert

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"io/ioutil"
)

// RevocationPolicy is a strategy for querying CRLs and evaluating if a certificate
// has been revoked.
type RevocationPolicy interface {
	VerifyCertificate(cert *x509.Certificate) error
}

// LoadCRLFile parses a revocation list from a file and uses it to statically mark
// certificates as valid or revoked. The signatures of the CRL and the CRL disturbation
// points of passed certificates are not evaluated.
func LoadCRLFile(file string) (RevocationPolicy, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("read file: %v", err)
	}
	certList, err := x509.ParseCRL(data)
	if err != nil {
		return nil, fmt.Errorf("parse crl: %v", err)
	}
	return &localCRL{certList}, nil
}

// localCRL is a static, in memory revocation list.
type localCRL struct {
	certList *pkix.CertificateList
}

func (c *localCRL) VerifyCertificate(cert *x509.Certificate) error {
	for _, r := range c.certList.TBSCertList.RevokedCertificates {
		if cert.SerialNumber.Cmp(r.SerialNumber) == 0 {
			return errors.New("certificate revoked")
		}
	}
	return nil
}
