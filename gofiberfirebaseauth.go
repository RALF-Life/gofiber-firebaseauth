// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber
// Special thanks to : https://github.com/LeafyCode/express-firebase-auth

package gofiberfirebaseauth

import (
	"errors"
	firebase "firebase.google.com/go"

	"github.com/gofiber/fiber/v2"
)

var ErrTokenMissingInHeader = errors.New("token missing")

// New - Signature Function
func New(app *firebase.App, config Config) fiber.Handler {
	cfg := configDefault(app, config)

	// Return authed handler
	return func(c *fiber.Ctx) error {

		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// 1) Get token from header
		tokenCandidate := cfg.TokenExtractor(c)
		if tokenCandidate == "" {
			return cfg.ErrorHandler(c, ErrTokenMissingInHeader)
		}

		// 2) Validate the Token Candidate
		token, err := cfg.Authorizer(tokenCandidate)
		if err != nil {
			return cfg.ErrorHandler(c, err)
		}

		// 3) If Token Candidate valid, return SuccessHandler
		if token == nil {
			return cfg.ErrorHandler(c, err)
		}

		// Do something with the token
		if err = cfg.TokenCallback(c, token); err != nil {
			return cfg.ErrorHandler(c, err)
		}

		return cfg.SuccessHandler(c)
	}
}
