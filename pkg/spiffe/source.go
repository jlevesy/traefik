package spiffe

import (
	"github.com/spiffe/go-spiffe/v2/bundle/x509bundle"
	"github.com/spiffe/go-spiffe/v2/svid/x509svid"
)

// SpiffeX509Source allows to retrieve an x509 sVID and an x509 bundle.
type X509Source interface {
	x509svid.Source
	x509bundle.Source
}
