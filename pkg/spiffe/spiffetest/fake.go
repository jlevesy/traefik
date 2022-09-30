package spiffetest

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net/url"
	"time"

	"github.com/spiffe/go-spiffe/v2/bundle/x509bundle"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/svid/x509svid"
)

// X509Source allows retrieving staticly an SVID and its associated bundle.
type X509Source struct {
	Bundle *x509bundle.Bundle
	SVID   *x509svid.SVID
}

func (s *X509Source) GetX509BundleForTrustDomain(trustDomain spiffeid.TrustDomain) (*x509bundle.Bundle, error) {
	return s.Bundle, nil
}

func (s *X509Source) GetX509SVID() (*x509svid.SVID, error) {
	return s.SVID, nil
}

// PKI simulates a SPIFFE aware PKI and allows generating multiple valid SVIDs.
type PKI struct {
	caPrivateKey *rsa.PrivateKey

	bundle *x509bundle.Bundle
}

func NewPKI(trustDomain spiffeid.TrustDomain) (*PKI, error) {
	caPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	caTemplate := x509.Certificate{
		SerialNumber: big.NewInt(2000),
		Subject: pkix.Name{
			Organization: []string{"spiffe"},
		},
		URIs:         []*url.URL{spiffeid.RequireFromPath(trustDomain, "/ca").URL()},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour),
		SubjectKeyId: []byte("ca"),
		KeyUsage: x509.KeyUsageCertSign |
			x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		PublicKey:             caPrivateKey.Public(),
	}
	if err != nil {
		return nil, err
	}

	caCertDER, err := x509.CreateCertificate(
		rand.Reader,
		&caTemplate,
		&caTemplate,
		caPrivateKey.Public(),
		caPrivateKey,
	)
	if err != nil {
		return nil, err
	}

	bundle, err := x509bundle.ParseRaw(
		trustDomain,
		caCertDER,
	)
	if err != nil {
		return nil, err
	}

	return &PKI{
		bundle:       bundle,
		caPrivateKey: caPrivateKey,
	}, nil
}

func (f *PKI) Bundle() *x509bundle.Bundle { return f.bundle }

func (f *PKI) GenSVID(id spiffeid.ID) (*x509svid.SVID, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(200001),
		URIs:         []*url.URL{id.URL()},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour),
		SubjectKeyId: []byte("svid"),
		KeyUsage: x509.KeyUsageKeyEncipherment |
			x509.KeyUsageKeyAgreement |
			x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		BasicConstraintsValid: true,
		PublicKey:             privateKey.PublicKey,
	}

	certDER, err := x509.CreateCertificate(
		rand.Reader,
		&template,
		f.bundle.X509Authorities()[0],
		privateKey.Public(),
		f.caPrivateKey,
	)
	if err != nil {
		return nil, err
	}

	keyPKCS8, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	return x509svid.ParseRaw(certDER, keyPKCS8)
}
