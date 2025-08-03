// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package azauth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/cockroachdb/errors"
	"github.com/kopexa-grc/x/vault"
	"github.com/rs/zerolog/log"
)

// TokenResolverFn is a function type that returns an Azure token credential and any error encountered.
// It's used to implement different token acquisition strategies that can be chained together.
type TokenResolverFn (func() (azcore.TokenCredential, error))

// WithStaticToken creates a TokenResolverFn that returns a pre-existing static token.
//
// This is useful for testing scenarios or when you already have a token from another source.
//
// Parameters:
//   - t: The token credential to return.
//
// Returns:
//   - A TokenResolverFn that returns the provided token credential.
func WithStaticToken(t azcore.TokenCredential) TokenResolverFn {
	return func() (azcore.TokenCredential, error) {
		return t, nil
	}
}

// WithCliCredentials creates a TokenResolverFn that uses the Azure CLI for authentication.
//
// This method acquires tokens using the credentials from Azure CLI.
//
// Parameters:
//   - opts: Options for configuring the Azure CLI credential.
//
// Returns:
//   - A TokenResolverFn that returns a CLI-based token credential.
func WithCliCredentials(opts *azidentity.AzureCLICredentialOptions) TokenResolverFn {
	return func() (azcore.TokenCredential, error) {
		return azidentity.NewAzureCLICredential(opts)
	}
}

// WithEnvCredentials creates a TokenResolverFn that uses environment variables for authentication.
//
// This method looks for standard Azure environment variables such as AZURE_TENANT_ID,
// AZURE_CLIENT_ID, and AZURE_CLIENT_SECRET or AZURE_CLIENT_CERTIFICATE_PATH.
//
// Parameters:
//   - opts: Options for configuring the environment credential.
//
// Returns:
//   - A TokenResolverFn that returns an environment-based token credential.
func WithEnvCredentials(opts *azidentity.EnvironmentCredentialOptions) TokenResolverFn {
	return func() (azcore.TokenCredential, error) {
		return azidentity.NewEnvironmentCredential(opts)
	}
}

// WithRetryableManagedIdentityCredentials creates a TokenResolverFn that uses Azure Managed Identity with retry capability.
//
// This function enhances the standard Managed Identity authentication by implementing retry logic and
// configurable timeouts. It addresses reliability issues when managed identity authentication
// encounters transient failures or timing out in certain environments.
//
// Parameters:
//   - timeout: Maximum duration to wait for each token acquisition attempt.
//   - attempts: Number of retry attempts before giving up.
//   - opts: Standard Azure Managed Identity credential options.
//
// Returns:
//   - A TokenResolverFn that returns a retryable managed identity credential.
func WithRetryableManagedIdentityCredentials(timeout time.Duration, attempts int, opts *azidentity.ManagedIdentityCredentialOptions) TokenResolverFn {
	return func() (azcore.TokenCredential, error) {
		mic, err := azidentity.NewManagedIdentityCredential(opts)
		if err != nil {
			return nil, err
		}
		return &retryableManagedIdentityCredential{mic: *mic, timeout: timeout, attempts: attempts}, nil
	}
}

// WithWorkloadIdentityCredentials creates a TokenResolverFn that uses Azure Workload Identity for authentication.
//
// Workload Identity allows applications running in Kubernetes to access Azure resources using
// Kubernetes service accounts federated with Azure AD.
//
// Parameters:
//   - opts: Options for configuring the workload identity credential.
//
// Returns:
//   - A TokenResolverFn that returns a workload identity-based token credential.
func WithWorkloadIdentityCredentials(opts *azidentity.WorkloadIdentityCredentialOptions) TokenResolverFn {
	return func() (azcore.TokenCredential, error) {
		return azidentity.NewWorkloadIdentityCredential(opts)
	}
}

// BuildChainedToken creates a ChainedTokenCredential from a sequence of TokenResolverFn functions.
//
// This function builds a credential chain that tries each authentication method in sequence
// until one succeeds. Only credentials that can be created without error are included in the chain.
//
// Parameters:
//   - opts: A variadic list of TokenResolverFn functions representing different authentication methods.
//
// Returns:
//   - A ChainedTokenCredential that will try each credential in sequence.
//   - Any error that occurred during creation of the chain.
func BuildChainedToken(opts ...TokenResolverFn) (*azidentity.ChainedTokenCredential, error) {
	chain := []azcore.TokenCredential{}
	for _, fn := range opts {
		cred, err := fn()
		if err == nil {
			chain = append(chain, cred)
		}
	}
	return azidentity.NewChainedTokenCredential(chain, nil)
}

const (
	DefaultRetryableManagedIdentityTimeout  = 5 * time.Second // Default timeout for retryable managed identity credential
	DefaultRetryableManagedIdentityAttempts = 3               // Default number of retry attempts for managed identity credential
)

// GetDefaultChainedToken creates a default credential chain with common authentication methods.
//
// The default chain tries the following authentication methods in order:
// 1. Azure CLI credentials
// 2. Environment-based credentials
// 3. Retryable managed identity credentials
// 4. Workload identity credentials
//
// Parameters:
//   - options: Configuration options for the default Azure credential chain.
//     If nil, default options will be used.
//
// Returns:
//   - A ChainedTokenCredential containing the default authentication methods.
//   - Any error that occurred during creation of the chain.
func GetDefaultChainedToken(options *azidentity.DefaultAzureCredentialOptions) (*azidentity.ChainedTokenCredential, error) {
	if options == nil {
		options = &azidentity.DefaultAzureCredentialOptions{}
	}
	opts := []TokenResolverFn{
		WithCliCredentials(&azidentity.AzureCLICredentialOptions{AdditionallyAllowedTenants: []string{"*"}}),
		WithEnvCredentials(&azidentity.EnvironmentCredentialOptions{ClientOptions: options.ClientOptions}),
		WithRetryableManagedIdentityCredentials(DefaultRetryableManagedIdentityTimeout, DefaultRetryableManagedIdentityAttempts, &azidentity.ManagedIdentityCredentialOptions{ClientOptions: options.ClientOptions}),
		WithWorkloadIdentityCredentials(&azidentity.WorkloadIdentityCredentialOptions{
			ClientOptions:            options.ClientOptions,
			DisableInstanceDiscovery: options.DisableInstanceDiscovery,
			TenantID:                 options.TenantID,
		}),
	}
	return BuildChainedToken(opts...)
}

// GetTokenFromCredential creates an Azure token credential from a given inventory credential.
//
// This function supports multiple credential types including certificates and passwords.
// If no credential is provided, it will fall back to using the default credential chain.
//
// Parameters:
//   - credential: The inventory credential containing authentication details.
//     If nil, a default credential chain will be used.
//   - tenantId: The Azure Active Directory tenant ID.
//   - clientId: The client (application) ID registered in Azure Active Directory.
//
// Returns:
//   - An Azure token credential that can be used for authentication.
//   - Any error that occurred during credential creation.
func GetTokenFromCredential(credential *vault.Credential, tenantID, clientID string) (azcore.TokenCredential, error) {
	var azCred azcore.TokenCredential
	var err error

	// Fallback to default authorizer if no credentials are specified
	if credential == nil {
		log.Debug().Msg("using default azure token chain resolver")
		azCred, err = GetDefaultChainedToken(&azidentity.DefaultAzureCredentialOptions{})
		if err != nil {
			return nil, errors.Wrap(err, "error creating CLI credentials")
		}
	} else {
		switch credential.Type {
		case vault.CredentialType_pkcs12:
			// Parse and use certificate credentials
			certs, privateKey, err := azidentity.ParseCertificates(credential.Secret, []byte(credential.Password))
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("could not parse provided certificate at %s", credential.PrivateKeyPath))
			}
			azCred, err = azidentity.NewClientCertificateCredential(tenantID, clientID, certs, privateKey, &azidentity.ClientCertificateCredentialOptions{})
			if err != nil {
				return nil, errors.Wrap(err, "error creating credentials from a certificate")
			}
		case vault.CredentialType_password:
			// Use client secret/password credentials
			azCred, err = azidentity.NewClientSecretCredential(tenantID, clientID, string(credential.Secret), &azidentity.ClientSecretCredentialOptions{})
			if err != nil {
				return nil, errors.Wrap(err, "error creating credentials from a secret")
			}
		default:
			return nil, errors.New("invalid secret configuration for microsoft transport: " + credential.Type.String())
		}
	}
	return azCred, nil
}

// retryableManagedIdentityCredential implements the azcore.TokenCredential interface
// and adds retry capabilities to the standard ManagedIdentityCredential.
//
// This type is used to overcome the limitations of the standard Azure Managed Identity
// implementation, particularly in scenarios where transient network issues or
// service unavailability might cause authentication to fail.
type retryableManagedIdentityCredential struct {
	mic      azidentity.ManagedIdentityCredential // The underlying managed identity credential
	attempts int                                  // Maximum number of attempts to acquire a token
	timeout  time.Duration                        // Timeout for each token acquisition attempt
}

// GetToken implements the azcore.TokenCredential interface to acquire an access token.
//
// This method attempts to acquire a token multiple times based on the configured
// number of attempts. Each failure is logged for debugging purposes.
//
// Parameters:
//   - ctx: The context for the token request.
//   - opts: Options for the token request, including scopes.
//
// Returns:
//   - An Azure access token if successful.
//   - An error representing all failed attempts if unsuccessful.
func (t *retryableManagedIdentityCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	// Ensure at least one attempt will be made
	if t.attempts < 1 {
		t.attempts = 1
	}

	errs := []error{}
	for i := 0; i < t.attempts; i++ {
		tk, err := t.tryGetToken(ctx, opts)
		if err == nil {
			return tk, nil
		}
		log.Debug().
			Err(err).
			Int("attempt", i+1).
			Int("max_attempts", t.attempts).
			Msg("failed to get managed identity token (may retry)")
		errs = append(errs, err)
	}

	log.Error().
		Int("num_attempts", t.attempts).
		Msg("failed to get managed identity token (max retries reached)")
	return azcore.AccessToken{}, errors.Join(errs...)
}

// tryGetToken attempts to acquire a token from the managed identity endpoint with a timeout.
//
// This is a helper method for GetToken that handles the timeout logic for a single
// token acquisition attempt. It also converts certain errors to more specific types
// to improve error handling and debugging.
//
// Parameters:
//   - ctx: The context for the token request.
//   - opts: Options for the token request, including scopes.
//
// Returns:
//   - An Azure access token if successful.
//   - An error if token acquisition fails.
func (t *retryableManagedIdentityCredential) tryGetToken(ctx context.Context, opts policy.TokenRequestOptions) (tk azcore.AccessToken, err error) {
	// Create a new timeout context for this attempt
	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()

	if t.timeout > 0 {
		c, cancel := context.WithTimeout(ctx, t.timeout)
		defer cancel()
		tk, err = t.mic.GetToken(c, opts)
		if err != nil {
			var authFailedErr *azidentity.AuthenticationFailedError
			// Improve error classification for timeout errors
			if errors.As(err, &authFailedErr) && strings.Contains(err.Error(), "context deadline exceeded") {
				err = azidentity.NewCredentialUnavailableError("managed identity request timed out")
			}
		} else {
			// If successful, we know the managed identity is working, so we can disable the timeout
			// for future calls to improve performance
			t.timeout = 0
		}
	} else {
		// If timeout is disabled (0), use the original context
		tk, err = t.mic.GetToken(ctx, opts)
	}
	return
}
