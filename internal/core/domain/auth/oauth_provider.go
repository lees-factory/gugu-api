package auth

type Provider string

const (
	ProviderGoogle Provider = "google"
)

type ProviderIdentity struct {
	Provider    Provider
	Subject     string
	Email       string
	Verified    bool
	DisplayName string
}
