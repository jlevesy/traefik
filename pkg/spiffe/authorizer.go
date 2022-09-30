package spiffe

import (
	"fmt"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
)

// BuildAuthorizer returns the correct authorizer based on given constraints.
func BuildAuthorizer(allowedSpiffeIDs []string, allowedTrustDomain string) (tlsconfig.Authorizer, error) {
	switch {
	case len(allowedSpiffeIDs) > 0:
		var (
			spiffeIDs = make([]spiffeid.ID, len(allowedSpiffeIDs))
			err       error
		)

		for i, rawID := range allowedSpiffeIDs {
			spiffeIDs[i], err = spiffeid.FromString(rawID)
			if err != nil {
				return nil, fmt.Errorf("invalid server spiffeID provided: %w", err)
			}
		}

		return tlsconfig.AuthorizeOneOf(spiffeIDs...), nil
	case allowedTrustDomain != "":
		trustDomain, err := spiffeid.TrustDomainFromString(allowedTrustDomain)
		if err != nil {
			return nil, fmt.Errorf("invalid spiffe trust domain provided: %w", err)
		}

		return tlsconfig.AuthorizeMemberOf(trustDomain), nil
	default:
		return tlsconfig.AuthorizeAny(), nil
	}
}
