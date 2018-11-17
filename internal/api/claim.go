package main

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type (
	ClaimModifier func(j *jwt.StandardClaims)

	ClaimFactory interface {
		Build(...ClaimModifier) *jwt.StandardClaims
		// Build
		SetOverride(ClaimModifier)
	}

	claimFactory struct {
		// secret string
		defaults  jwt.StandardClaims
		modifiers []ClaimModifier
		override  ClaimModifier
	}
)

func NewClaimFactory(defaults jwt.StandardClaims, modifiers ...ClaimModifier) *claimFactory {
	return &claimFactory{
		defaults:  defaults,
		modifiers: modifiers,
	}
}

func NewProductionClaimFactory(aud, iss string) *claimFactory {
	return &claimFactory{
		defaults: jwt.StandardClaims{
			Audience: aud,
			Issuer:   iss,
		},
	}
}

func (c *claimFactory) Build(extras ...ClaimModifier) *jwt.StandardClaims {
	result := c.defaults
	for _, modifier := range append(c.modifiers, extras...) {
		modifier(&result)
	}
	if c.override != nil {
		c.override(&result)
	}
	return &result
}

func (c *claimFactory) SetOverride(override ClaimModifier) {
	c.override = override
}

func makeAudienceModifier(aud string) ClaimModifier {
	return func(j *jwt.StandardClaims) {
		j.Audience = aud
	}
}
func makeExpireAtModifier(now time.Time, expiresIn time.Duration) ClaimModifier {
	return func(j *jwt.StandardClaims) {
		j.ExpiresAt = now.Add(expiresIn).Unix()
	}
}
func makeIssuedAtModifier(now time.Time) ClaimModifier {
	return func(j *jwt.StandardClaims) {
		j.IssuedAt = now.Unix()
	}
}
func makeIssuerModifier(iss string) ClaimModifier {
	return func(j *jwt.StandardClaims) {
		j.Issuer = iss
	}
}

func makeSubjectModifier(sub string) ClaimModifier {
	return func(j *jwt.StandardClaims) {
		j.Subject = sub
	}
}
