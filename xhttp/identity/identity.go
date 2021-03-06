package identity

import (
	"context"
	"net/http"
)

// GuestRoleName is default role name for guest
const GuestRoleName = "guest"

// Identity contains information about the identity of an API caller
type Identity interface {
	String() string
	UserID() string
	Role() string
	Name() string
	UserInfo() interface{}
}

// Mapper returns Identity from supplied HTTP request
type Mapper func(*http.Request) (Identity, error)

// NewIdentity returns a new Identity instance with the indicated role
func NewIdentity(role, name, userID string) Identity {
	return identity{
		role:   role,
		name:   name,
		userID: userID,
	}
}

// NewIdentityWithUserInfo returns a new Identity instance with the indicated role and user info
func NewIdentityWithUserInfo(role, name, userID string, userInfo interface{}) Identity {
	return identity{
		role:     role,
		name:     name,
		userID:   userID,
		userInfo: userInfo,
	}
}

type identity struct {
	// name of identity
	// It can be CommonName extracted from certificate
	name string
	// role of identity
	role string
	// ID of the user
	userID string
	// extra user info, specific to the application
	userInfo interface{}
}

// UserID returns user ID
func (c identity) UserID() string {
	return c.userID
}

// Name returns the clients name
func (c identity) Name() string {
	return c.name
}

// Role returns the clients role
func (c identity) Role() string {
	return c.role
}

// UserInfo returns application specific user info
func (c identity) UserInfo() interface{} {
	return c.userInfo
}

// String returns the identity as a single string value
// in the format of role/name
func (c identity) String() string {
	if c.role != c.name {
		return c.role + "/" + c.name
	}
	return c.role
}

// GuestIdentityMapper always returns "guest" for the role
func GuestIdentityMapper(r *http.Request) (Identity, error) {
	var name string
	if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
		name = ClientIPFromRequest(r)
	} else {
		name = r.TLS.PeerCertificates[0].Subject.CommonName
	}
	return NewIdentity(GuestRoleName, name, ""), nil
}

// WithTestIdentity is used in unit tests to set HTTP request identity
func WithTestIdentity(r *http.Request, identity Identity) *http.Request {
	ctx := &RequestContext{
		identity:      identity,
		correlationID: extractCorrelationID(r),
		clientIP:      nodeInfoFactory().LocalIP(),
	}
	c := context.WithValue(r.Context(), keyContext, ctx)
	return r.WithContext(c)
}
