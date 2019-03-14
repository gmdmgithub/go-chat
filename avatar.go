package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

// ErrNoAvatarURL is the error that is returned when the
// Avatar instance is unable to provide an avatar URL.
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")

// Avatar represents types capable of representing
// user profile pictures.
type Avatar interface {
	// GetAvatarURL gets the avatar URL for the specified client,
	// or returns an error if something goes wrong.
	// ErrNoAvatarURL is returned if the object is unable to get
	// a URL for the specified client.
	GetAvatarURL(c *client) (string, error)
}

// AuthAvatar ...
type AuthAvatar struct{}

// UseAuthAvatar ...
var UseAuthAvatar AuthAvatar

// GetAvatarURL abstract a method of user avatar_url
func (a AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["picture"]; ok { // TODO use something more general for "picture"
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

// GravatarAvatar ...
type GravatarAvatar struct{}

// UseGravatar ...
var UseGravatar GravatarAvatar

// GetAvatarURL ...
func (a GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if email, ok := c.userData["email"]; ok {
		if emailStr, ok := email.(string); ok {
			m := md5.New()
			io.WriteString(m, strings.ToLower(emailStr))
			return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)), nil
		}
	} else {
		log.Printf("email problem %v", c.userData)
	}
	return "", ErrNoAvatarURL
}
