package impl

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

type VerifyCertificateResult struct {
	Certificate x509.Certificate
	Verified    bool
}

func (c *Cluster) VerifyCertificate(certData []byte) (*VerifyCertificateResult, error) {
	certDecoded, _ := pem.Decode(certData)

	if certDecoded == nil {
		return nil, errors.New("PEM decoding error")
	}

	certs, err := x509.ParseCertificates(certDecoded.Bytes)

	if err != nil {
		return nil, fmt.Errorf("certificate parsing failed: %v", err)
	}

	if len(certs) == 0 {
		return nil, errors.New("no certificates found")
	}

	intermediateStore := x509.NewCertPool()

	for _, cert := range certs[1:] {
		intermediateStore.AddCert(cert)
	}

	authorityStore := x509.NewCertPool()

	for _, cert := range c.ClientAuthorities {
		authorityStore.AddCert(cert)
	}

	_, err = certs[0].Verify(x509.VerifyOptions{Intermediates: intermediateStore, Roots: authorityStore})

	if err != nil {
		return &VerifyCertificateResult{Certificate: *certs[0], Verified: false}, nil
	}

	return &VerifyCertificateResult{Certificate: *certs[0], Verified: true}, nil
}
