package discord

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
)

const (
	authURL    = "https://discord.com/oauth2/authorize"
	tokenURL   = "https://discord.com/api/oauth2/token"
	profileURL = "https://discord.com/api/users/@me"
	name       = "Discord"
)

var _ idp.Provider = (*Provider)(nil)

// Provider is the [idp.Provider] implementation for Discord.
// Discord does not implement OpenID Connect, so it is built on top of the
// generic [oauth.Provider] with a mapper hardcoded to Discord's user object shape.
type Provider struct {
	*oauth.Provider
}

// New creates a Discord provider using the [oauth.Provider] (OAuth 2.0 generic provider).
func New(clientID, clientSecret, callbackURL string, scopes []string, httpClient *http.Client, options ...oauth.ProviderOpts) (*Provider, error) {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  callbackURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		Scopes: scopes,
	}
	rp, err := oauth.New(
		config,
		name,
		profileURL,
		func() idp.User {
			return new(User)
		},
		httpClient,
		options...,
	)
	if err != nil {
		return nil, err
	}
	return &Provider{Provider: rp}, nil
}

// User represents the authenticated Discord user and implements the [idp.User] interface.
// https://discord.com/developers/docs/resources/user#user-object
type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	GlobalName    string `json:"global_name"`
	Avatar        string `json:"avatar"`
	Email         string `json:"email"`
	Verified      bool   `json:"verified"`
	Locale        string `json:"locale"`
}

// GetID is an implementation of the [idp.User] interface.
func (u *User) GetID() string {
	return u.ID
}

// GetFirstName is an implementation of the [idp.User] interface.
// It returns an empty string because Discord does not provide the user's first name.
func (u *User) GetFirstName() string {
	return ""
}

// GetLastName is an implementation of the [idp.User] interface.
// It returns an empty string because Discord does not provide the user's last name.
func (u *User) GetLastName() string {
	return ""
}

// GetDisplayName is an implementation of the [idp.User] interface.
// It prefers Discord's "global display name" and falls back to the username.
func (u *User) GetDisplayName() string {
	if u.GlobalName != "" {
		return u.GlobalName
	}
	return u.Username
}

// GetNickname is an implementation of the [idp.User] interface, returning the Discord username.
func (u *User) GetNickname() string {
	return u.Username
}

// GetPreferredUsername is an implementation of the [idp.User] interface, returning the Discord username.
func (u *User) GetPreferredUsername() string {
	return u.Username
}

// GetEmail is an implementation of the [idp.User] interface.
func (u *User) GetEmail() domain.EmailAddress {
	return domain.EmailAddress(u.Email)
}

// IsEmailVerified is an implementation of the [idp.User] interface.
func (u *User) IsEmailVerified() bool {
	return u.Verified
}

// GetPhone is an implementation of the [idp.User] interface.
// It returns an empty string because Discord does not provide the user's phone number.
func (u *User) GetPhone() domain.PhoneNumber {
	return ""
}

// IsPhoneVerified is an implementation of the [idp.User] interface.
func (u *User) IsPhoneVerified() bool {
	return false
}

// GetPreferredLanguage is an implementation of the [idp.User] interface.
func (u *User) GetPreferredLanguage() language.Tag {
	tag, err := language.Parse(u.Locale)
	if err != nil {
		return language.Und
	}
	return tag
}

// GetProfile is an implementation of the [idp.User] interface.
// It returns an empty string because Discord does not provide a public profile URL.
func (u *User) GetProfile() string {
	return ""
}

// GetAvatarURL is an implementation of the [idp.User] interface.
// Discord only returns an avatar hash, so the CDN URL has to be constructed.
// https://discord.com/developers/docs/reference#image-formatting
func (u *User) GetAvatarURL() string {
	if u.Avatar == "" || u.ID == "" {
		return ""
	}
	return fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", u.ID, u.Avatar)
}
