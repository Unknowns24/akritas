package githubapp

import (
	"context"
	"crypto/sha256"
	"net/url"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func (uc *UseCase) CompleteManifest(ctx context.Context, code, state string) (portsin.GitHubManifestCallbackResult, error) {
	if len(code) < 20 || len(code) > 255 || len(state) < 32 || len(state) > 512 {
		return portsin.GitHubManifestCallbackResult{}, domain.ErrManifestStateInvalid
	}
	digest := sha256.Sum256([]byte(state))
	registration, err := uc.store.ConsumeConversionState(ctx, digest[:], uc.now().UTC())
	if err != nil {
		return portsin.GitHubManifestCallbackResult{}, err
	}
	conversion, err := uc.gateway.ExchangeManifest(ctx, code)
	if err != nil {
		return portsin.GitHubManifestCallbackResult{}, err
	}
	defer wipe(conversion.PrivateKey)
	defer wipe(conversion.WebhookSecret)
	installationState, err := uc.newState()
	if err != nil {
		return portsin.GitHubManifestCallbackResult{}, domain.ErrIntegrationUnavailable.Wrap(err)
	}
	installationDigest := sha256.Sum256([]byte(installationState))
	registration.InstallationStateDigest = installationDigest[:]
	registration.AppID = &conversion.AppID
	registration.AppSlug = conversion.AppSlug
	registration.AppName = conversion.AppName
	registration.ClientID = conversion.ClientID
	registration.Status = portsout.GitHubAppRegistrationConverted
	registration.UpdatedAt = uc.now().UTC()
	secrets := []portsout.SecretValue{
		{Kind: portsout.SecretKindGitHubPrivateKey, Plaintext: conversion.PrivateKey},
		{Kind: portsout.SecretKindGitHubWebhook, Plaintext: conversion.WebhookSecret},
	}
	if err := uc.store.CompleteConversion(ctx, registration, secrets); err != nil {
		return portsin.GitHubManifestCallbackResult{}, err
	}
	redirect := "https://github.com/apps/" + url.PathEscape(strings.TrimSpace(conversion.AppSlug)) + "/installations/new?state=" + url.QueryEscape(installationState)
	return portsin.GitHubManifestCallbackResult{RedirectURL: redirect}, nil
}
