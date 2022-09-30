package tls

import "github.com/traefik/traefik/v2/pkg/types"

const certificateHeader = "-----BEGIN CERTIFICATE-----\n"

// +k8s:deepcopy-gen=true

// ClientAuth defines the parameters of the client authentication part of the TLS connection, if any.
type ClientAuth struct {
	CAFiles []FileOrContent `json:"caFiles,omitempty" toml:"caFiles,omitempty" yaml:"caFiles,omitempty"`
	// ClientAuthType defines the client authentication type to apply.
	// The available values are: "NoClientCert", "RequestClientCert", "VerifyClientCertIfGiven" and "RequireAndVerifyClientCert".
	ClientAuthType string `json:"clientAuthType,omitempty" toml:"clientAuthType,omitempty" yaml:"clientAuthType,omitempty" export:"true"`
}

// Supported SPIFFE use cases.
const (
	// Traefik serves its SPIFFE certificate, without requiring any client certificate.
	SpiffeModeTLS = "TLS"
	// Traefik serves its SPIFFE certificate, and requires a SPIFFE client certficate.
	SpiffeModeMTLS = "mTLS"
	// Traefik serves a configured certificate, but requires a SPIFFE client certificate
	SpiffeModeMTLSWebserver = "mTLSWebServer"
)

// +k8s:deepcopy-gen=true

// SpiffeOptions configures SPIFFE for an entry point.
type SpiffeOptions struct {
	Mode        string   `description:"Spiffe mode to use" json:"mode,ommitempty" toml:"mode,omitempty" yaml:"mode,omitempty" export:"true"`
	ClientIDs   []string `description:"Lists all allowed client spiffeIDs. Takes precedences over TrustDomain." json:"clientIDs,omitempty" toml:"clientIDs,omitempty" yaml:"clientIDs,omitempty" export:"true"`
	TrustDomain string   `description:"Specifies an allowed Spiffe trust domain for clients." json:"trustDomain,omitempty" yaml:"trustDomain,omitempty" toml:"trustDomain,omitempty"`
}

func (s *SpiffeOptions) NeedsClientCertValidation() bool {
	return s.Mode != "" && s.Mode == SpiffeModeTLS
}

func (s *SpiffeOptions) NeedsServingSVIDCertificate() bool {
	return s.Mode != "" && s.Mode != SpiffeModeMTLSWebserver
}

// +k8s:deepcopy-gen=true

// Options configures TLS for an entry point.
type Options struct {
	MinVersion               string         `json:"minVersion,omitempty" toml:"minVersion,omitempty" yaml:"minVersion,omitempty" export:"true"`
	MaxVersion               string         `json:"maxVersion,omitempty" toml:"maxVersion,omitempty" yaml:"maxVersion,omitempty" export:"true"`
	CipherSuites             []string       `json:"cipherSuites,omitempty" toml:"cipherSuites,omitempty" yaml:"cipherSuites,omitempty" export:"true"`
	CurvePreferences         []string       `json:"curvePreferences,omitempty" toml:"curvePreferences,omitempty" yaml:"curvePreferences,omitempty" export:"true"`
	ClientAuth               ClientAuth     `json:"clientAuth,omitempty" toml:"clientAuth,omitempty" yaml:"clientAuth,omitempty"`
	SniStrict                bool           `json:"sniStrict,omitempty" toml:"sniStrict,omitempty" yaml:"sniStrict,omitempty" export:"true"`
	PreferServerCipherSuites bool           `json:"preferServerCipherSuites,omitempty" toml:"preferServerCipherSuites,omitempty" yaml:"preferServerCipherSuites,omitempty" export:"true"` // Deprecated: https://github.com/golang/go/issues/45430
	ALPNProtocols            []string       `json:"alpnProtocols,omitempty" toml:"alpnProtocols,omitempty" yaml:"alpnProtocols,omitempty" export:"true"`
	Spiffe                   *SpiffeOptions `json:"spiffe,omitempty" toml:"spiffe,omitempty" yaml:"spiffe,omitempty"`
}

// SetDefaults sets the default values for an Options struct.
func (o *Options) SetDefaults() {
	// ensure http2 enabled
	o.ALPNProtocols = DefaultTLSOptions.ALPNProtocols
}

// +k8s:deepcopy-gen=true

// Store holds the options for a given Store.
type Store struct {
	DefaultCertificate   *Certificate   `json:"defaultCertificate,omitempty" toml:"defaultCertificate,omitempty" yaml:"defaultCertificate,omitempty" export:"true"`
	DefaultGeneratedCert *GeneratedCert `json:"defaultGeneratedCert,omitempty" toml:"defaultGeneratedCert,omitempty" yaml:"defaultGeneratedCert,omitempty" export:"true"`
}

// +k8s:deepcopy-gen=true

// GeneratedCert defines the default generated certificate configuration.
type GeneratedCert struct {
	// Resolver is the name of the resolver that will be used to issue the DefaultCertificate.
	Resolver string `json:"resolver,omitempty" toml:"resolver,omitempty" yaml:"resolver,omitempty" export:"true"`
	// Domain is the domain definition for the DefaultCertificate.
	Domain *types.Domain `json:"domain,omitempty" toml:"domain,omitempty" yaml:"domain,omitempty" export:"true"`
}

// +k8s:deepcopy-gen=true

// CertAndStores allows mapping a TLS certificate to a list of entry points.
type CertAndStores struct {
	Certificate `yaml:",inline" export:"true"`
	Stores      []string `json:"stores,omitempty" toml:"stores,omitempty" yaml:"stores,omitempty" export:"true"`
}
