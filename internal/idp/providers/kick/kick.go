package kick

import (
	"net/http"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
)

// NOTE: Kick's public OAuth2 API is comparatively new and less widely documented
// than Discord's or Twitch's. The endpoint URLs and response shape below reflect
// Kick's public API documentation at the time of writing (https://docs.kick.com) -
// verify them against the live docs before relying on this in production.
const (
	authURL    = "https://id.kick.com/oauth/authorize"
	tokenURL   = "https://id.kick.com/oauth/token"
	profileURL = "https://api.kick.com/public/v1/users"
	name       = "Kick"
)

var _ idp.Provider = (*Provider)(nil)

// Provider is the [idp.Provider] implementation for Kick.
// Kick's public API is a plain OAuth2 API (no OpenID Connect), so it is built on top of the
// generic [oauth.Provider] with a mapper hardcoded to Kick's user object shape.
type Provider struct {
	*oauth.Provider
}

// New creates a Kick provider using the [oauth.Provider] (OAuth 2.0 generic provider).
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

// User represents the authenticated Kick user and implements the [idp.User] interface.
// Kick's "get authenticated user" endpoint wraps the user object in a `data` array.
type User struct {
	Data []struct {
		UserID         int    `json:"user_id"`
		Name           string `json:"name"`
		Email          string `json:"email"`
		ProfilePicture string `json:"profile_picture"`
	} `json:"data"`
}

func (u *User) entry() (userID int, name, email, picture string) {
	if len(u.Data) == 0 {
		return 0, "", "", ""
	}
	d := u.Data[0]
	return d.UserID, d.Name, d.Email, d.ProfilePicture
}

// GetID is an implementation of the [idp.User] interface.
func (u *User) GetID() string {
	id, _, _, _ := u.entry()
	if id == 0 {
		return ""
	}
	return strconv.Itoa(id)
}

// GetFirstName is an implementation of the [idp.User] interface.
// It returns an empty string because Kick does not provide the user's first name.
func (u *User) GetFirstName() string {
	return ""
}

// GetLastName is an implementation of the [idp.User] interface.
// It returns an empty string because Kick does not provide the user's last name.
func (u *User) GetLastName() string {
	return ""
}

// GetDisplayName is an implementation of the [idp.User] interface.
func (u *User) GetDisplayName() string {
	_, name, _, _ := u.entry()
	return name
}

// GetNickname is an implementation of the [idp.User] interface.
func (u *User) GetNickname() string {
	_, name, _, _ := u.entry()
	return name
}

// GetPreferredUsername is an implementation of the [idp.User] interface.
func (u *User) GetPreferredUsername() string {
	_, name, _, _ := u.entry()
	return name
}

// GetEmail is an implementation of the [idp.User] interface.
func (u *User) GetEmail() domain.EmailAddress {
	_, _, email, _ := u.entry()
	return domain.EmailAddress(email)
}

// IsEmailVerified is an implementation of the [idp.User] interface.
// Kick does not expose a separate verified flag, so a present email is treated
// as verified, consistent with how other providers without that flag are handled (e.g. GitHub).
func (u *User) IsEmailVerified() bool {
	_, _, email, _ := u.entry()
	return email != ""
}

// GetPhone is an implementation of the [idp.User] interface.
// It returns an empty string because Kick does not provide the user's phone number.
func (u *User) GetPhone() domain.PhoneNumber {
	return ""
}

// IsPhoneVerified is an implementation of the [idp.User] interface.
func (u *User) IsPhoneVerified() bool {
	return false
}

// GetPreferredLanguage is an implementation of the [idp.User] interface.
// It returns [language.Und] because Kick does not provide the user's language.
func (u *User) GetPreferredLanguage() language.Tag {
	return language.Und
}

// GetProfile is an implementation of the [idp.User] interface.
// It returns an empty string because Kick does not provide a public profile URL.
func (u *User) GetProfile() string {
	return ""
}

// GetAvatarURL is an implementation of the [idp.User] interface.
func (u *User) GetAvatarURL() string {
	_, _, _, picture := u.entry()
	return picture
}
