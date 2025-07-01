package datasource

import "os"

type GoogleApiClient struct {
	clientId    string
	secret      string
	redirectUri string
}

func NewGoogleApiClient() *GoogleApiClient {
	return &GoogleApiClient{
		clientId:    getGoogleClientId(),
		secret:      getGoogleClientSecret(),
		redirectUri: getGoogleRedirectUri(),
	}
}

func (c *GoogleApiClient) ClientId() string {
	return c.clientId
}

func (c *GoogleApiClient) ClientSecret() string {
	return c.secret
}

func (c *GoogleApiClient) RedirectURI() string {
	return c.redirectUri
}

func getGoogleClientId() string {
	clientId, ok := os.LookupEnv("GOOGLE_OAUTH_CLIENT_ID")
	if !ok {
		panic("\"GOOGLE_OAUTH_CLIENT_ID\" is no set")
	}
	return clientId
}

func getGoogleClientSecret() string {
	secret, ok := os.LookupEnv("GOOGLE_OAUTH_CLIENT_SECRET")
	if !ok {
		panic("\"GOOGLE_OAUTH_CLIENT_SECRET\" is no set")
	}
	return secret
}

func getGoogleRedirectUri() string {
	uri, ok := os.LookupEnv("GOOGLE_OAUTH_REDIRECT_URI")
	if !ok {
		panic("\"GOOGLE_OAUTH_REDIRECT_URI\" is no set")
	}
	return uri
}
