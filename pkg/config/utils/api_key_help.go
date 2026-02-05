// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package utils

import (
	"net/url"
	"strings"

	pkgconfigmodel "github.com/DataDog/datadog-agent/pkg/config/model"
)

const apiKeyDocsURL = "https://docs.datadoghq.com/account_management/api-app-keys/"

// apiKeyHelpURLs holds the org-settings and validate URLs derived from config.
type apiKeyHelpURLs struct {
	orgSettings string
	validate    string
	docs        string
}

func getAPIKeyHelpURLs(c pkgconfigmodel.Reader) apiKeyHelpURLs {
	infra := GetInfraEndpoint(c)
	infra = strings.TrimSuffix(infra, ".")
	u, err := url.Parse(infra)
	if err != nil {
		return apiKeyHelpURLs{docs: apiKeyDocsURL}
	}
	base := strings.TrimSuffix(u.String(), "/")
	orgSettings := base + "/organization-settings/api-keys"
	// Validate URL: app. -> api. in host (e.g. app.datadoghq.com -> api.datadoghq.com/api/v1/validate)
	validateHost := strings.TrimSuffix(u.Host, ".")
	validateURL := base + "/api/v1/validate"
	if u.Host != "" {
		if idx := strings.Index(validateHost, "."); idx > 0 && validateHost[:idx] == "app" {
			validateHost = "api." + validateHost[idx+1:]
			validateURL = u.Scheme + "://" + validateHost + "/api/v1/validate"
		}
	}
	return apiKeyHelpURLs{orgSettings: orgSettings, validate: validateURL, docs: apiKeyDocsURL}
}

// APIKeyInvalidHelpMessage returns a concise, actionable message for invalid API key (e.g. 403) errors.
// It uses the configured site (GetInfraEndpoint) so org-settings and validate URLs match the user's Datadog site.
func APIKeyInvalidHelpMessage(c pkgconfigmodel.Reader) string {
	return APIKeyInvalidHelpMessageForEndpoint(c, "")
}

// APIKeyInvalidHelpMessageForEndpoint returns the same message as APIKeyInvalidHelpMessage, and when
// endpoint is non-empty appends which endpoint failed so users know what data may be missing (e.g. metrics, logs).
func APIKeyInvalidHelpMessageForEndpoint(c pkgconfigmodel.Reader, endpoint string) string {
	urls := getAPIKeyHelpURLs(c)
	msg := "Verify API key in Organization Settings → API Keys (" + urls.orgSettings + "). " +
		"Ensure the `site` (or `dd_url`) in your config matches the Datadog site where your API key was created (e.g. datadoghq.com vs datadoghq.eu). " +
		"Check datadog.yaml api_key and DD_API_KEY. " +
		"Test: curl -X POST " + urls.validate + " -H 'DD-API-KEY: <key>'. " +
		"Docs: " + urls.docs
	if endpoint != "" {
		msg += " This failure occurred when sending to: " + endpoint + "."
	}
	return msg
}

// APIKeyMissingHelpMessage returns a concise, actionable message for missing API key errors.
// Same site-specific links as APIKeyInvalidHelpMessage, plus guidance to set the key.
func APIKeyMissingHelpMessage(c pkgconfigmodel.Reader) string {
	urls := getAPIKeyHelpURLs(c)
	return "Specify an API key in datadog.yaml (api_key) or set DD_API_KEY. " +
		"Verify in Organization Settings → API Keys (" + urls.orgSettings + "). " +
		"Ensure `site` (or `dd_url`) matches the site where your key was created. " +
		"Docs: " + urls.docs
}
