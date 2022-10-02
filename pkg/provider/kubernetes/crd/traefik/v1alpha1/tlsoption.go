package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:storageversion

// TLSOption is the CRD implementation of a Traefik TLS Option, allowing to configure some parameters of the TLS connection.
// More info: https://doc.traefik.io/traefik/v2.9/https/tls/#tls-options
type TLSOption struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata"`

	Spec TLSOptionSpec `json:"spec"`
}

// +k8s:deepcopy-gen=true

// TLSOptionSpec defines the desired state of a TLSOption.
type TLSOptionSpec struct {
	// MinVersion defines the minimum TLS version that Traefik will accept.
	// Possible values: VersionTLS10, VersionTLS11, VersionTLS12, VersionTLS13.
	// Default: VersionTLS10.
	MinVersion string `json:"minVersion,omitempty"`
	// MaxVersion defines the maximum TLS version that Traefik will accept.
	// Possible values: VersionTLS10, VersionTLS11, VersionTLS12, VersionTLS13.
	// Default: None.
	MaxVersion string `json:"maxVersion,omitempty"`
	// CipherSuites defines the list of supported cipher suites for TLS versions up to TLS 1.2.
	// More info: https://doc.traefik.io/traefik/v2.9/https/tls/#cipher-suites
	CipherSuites []string `json:"cipherSuites,omitempty"`
	// CurvePreferences defines the preferred elliptic curves in a specific order.
	// More info: https://doc.traefik.io/traefik/v2.9/https/tls/#curve-preferences
	CurvePreferences []string `json:"curvePreferences,omitempty"`
	// ClientAuth defines the server's policy for TLS Client Authentication.
	ClientAuth ClientAuth `json:"clientAuth,omitempty"`
	// SniStrict defines whether Traefik allows connections from clients connections that do not specify a server_name extension.
	SniStrict bool `json:"sniStrict,omitempty"`
	// PreferServerCipherSuites defines whether the server chooses a cipher suite among his own instead of among the client's.
	// It is enabled automatically when minVersion or maxVersion is set.
	// Deprecated: https://github.com/golang/go/issues/45430
	PreferServerCipherSuites bool `json:"preferServerCipherSuites,omitempty"`
	// ALPNProtocols defines the list of supported application level protocols for the TLS handshake, in order of preference.
	// More info: https://doc.traefik.io/traefik/v2.9/https/tls/#alpn-protocols
	ALPNProtocols []string `json:"alpnProtocols,omitempty"`
	// Spiffe configures server side SPIFFE on this router.
	Spiffe SpiffeOptions `json:"spiffe,omitempty"`
}

// +k8s:deepcopy-gen=true

// ClientAuth holds the TLS client authentication configuration.
type ClientAuth struct {
	// SecretNames defines the names of the referenced Kubernetes Secret storing certificate details.
	SecretNames []string `json:"secretNames,omitempty"`
	// ClientAuthType defines the client authentication type to apply.
	// +kubebuilder:validation:Enum=NoClientCert;RequestClientCert;RequireAnyClientCert;VerifyClientCertIfGiven;RequireAndVerifyClientCert
	ClientAuthType string `json:"clientAuthType,omitempty"`
}

// +k8s:deepcopy-gen=true

// SpiffeOptions configures server side SPIFFE.
type SpiffeOptions struct {
	// Mode defines Traefik behavior regarding SPIFFE certificates. There are 4 possible settings:
	// - Empty: SPIFFE is not enabled for this router.
	// - TLS: Traefik serves its SPIFFE ceritifate, but does not require nor validate any client certificate
	// - mTLS: Traefik serves its SPIFFE certificate, and checks the validity of client SPIFFE certificate. Optionally checks the client certificate Spiffe ID if either ClientIDs or TrustDomain are given.
	// - mTLSWebServer: Traefik serves the configured certificate but checks the validation of the client SPIFFE certificate. Optionally checks the client certificate Spiffe ID if either ClientIDs or TrustDomain are given.
	// +kubebuilder:validation:Enum=TLS;mTLS;mTLSWebServer
	Mode string `json:"mode,omitempty"`
	// List of allowed client SpiffeIDs, only applied if mode is mTLS or mTLSWebServer, takes precedence over TrustDomain.
	ClientIDs []string `json:"clientIDs,omitempty"`
	// Allowed SPIFFE trust domain for client certificates, only applied if mode is mTLS or mTLSWebServer.
	TrustDomain string `json:"trustDomain,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TLSOptionList is a collection of TLSOption resources.
type TLSOptionList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata"`

	// Items is the list of TLSOption.
	Items []TLSOption `json:"items"`
}
