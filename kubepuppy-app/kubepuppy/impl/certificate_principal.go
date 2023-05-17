package impl

import "crypto/x509"

type CertificatePrincipal struct {
	Certificate x509.Certificate
	Subjects    []Subject
}

var _ Principal = (*CertificatePrincipal)(nil)

func NewCertificatePrincipal(cert x509.Certificate) (*CertificatePrincipal, error) {

	subjects := make([]Subject, 0)

	subjects = append(subjects, NewUserSubject(cert.Subject.CommonName))

	for _, organization := range cert.Subject.Organization {
		subjects = append(subjects, NewGroupSubject(organization))
	}

	return &CertificatePrincipal{Certificate: cert, Subjects: subjects}, nil
}

func (c CertificatePrincipal) Matches(target Subject) bool {

	for _, subject := range c.Subjects {
		if subject.Matches(target) {
			return true
		}
	}

	return false
}
