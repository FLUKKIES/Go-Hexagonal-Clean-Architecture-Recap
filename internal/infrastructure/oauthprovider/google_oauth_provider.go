package oauthprovider

import (
	"context"
	"encoding/json"
	"io"

	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/ports"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleOAuthProvider struct {
	config *oauth2.Config
}

func NewGoogleOAuthProvider(clientID, clientSecret, redirectURL string) ports.IOAuthProvider {
	return &googleOAuthProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (g *googleOAuthProvider) ProviderName() string { return "google" }

func (g *googleOAuthProvider) GetAuthURL(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

func (g *googleOAuthProvider) GetUserProfile(code string) (*ports.OAuthUserProfile, error) {
	// แลก Code เป็น Token
	token, err := g.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	// เรียก Google UserInfo API
	client := g.config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info struct {
		Sub        string `json:"sub"`   // Google User ID
		Email      string `json:"email"`
		GivenName  string `json:"given_name"`
		FamilyName string `json:"family_name"`
		Picture    string `json:"picture"`
	}
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}

	return &ports.OAuthUserProfile{
		ProviderID: info.Sub,
		Email:      info.Email,
		FirstName:  info.GivenName,
		LastName:   info.FamilyName,
		ProfileUrl: info.Picture,
	}, nil
}
