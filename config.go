// 1.2.1-alpha

package gofiberfirebaseauth

import (
	"context"
	"errors"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gofiber/fiber/v2"
)

var ErrEmailNotVerified = errors.New("email not verified")

type User struct {
	EmailVerified bool
	UserID, Email string
}

// Config defines the config for middleware
type Config struct {
	// Skip Email Check.
	// Optional. Default: nil
	CheckEmailVerified bool

	// Ignore email verification for these routes
	// Optional. Default: nil
	CheckEmailVerifiedIgnoredUrls []string

	// Filter defines a function to skip middleware.
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// Authorizer defines a function which authenticate the Authorization token and return the authenticated token
	// Optional. Default: nil
	Authorizer func(string) (*auth.Token, error)

	// SuccessHandler defines a function which is executed for a valid token.
	// Optional. Default: nil
	SuccessHandler fiber.Handler

	// ErrorHandler defines a function which is executed for an invalid token.
	// It may be used to define a custom JWT error.
	// Optional. Default: nil
	ErrorHandler fiber.ErrorHandler

	// Context key to store user information from the token into context.
	// Optional. Default: "user".
	ContextKey string

	TokenExtractor ExtractorFun
	TokenCallback  func(ctx *fiber.Ctx, token *auth.Token) error
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	},
	SuccessHandler: func(c *fiber.Ctx) error {
		return c.Next()
	},
	ContextKey:     "user",
	TokenExtractor: NewHeaderExtractor(),
}

// Initializer
func configDefault(app *firebase.App, config ...Config) Config {
	// Return default config if nothing provided
	var cfg Config
	if len(config) > 0 {
		cfg = config[0]
	}

	// Set default values
	if cfg.ContextKey == "" {
		cfg.ContextKey = ConfigDefault.ContextKey
	}
	if cfg.Next == nil {
		cfg.Next = ConfigDefault.Next
	}
	if cfg.SuccessHandler == nil {
		cfg.SuccessHandler = ConfigDefault.SuccessHandler
	}
	if cfg.Authorizer == nil {
		cfg.Authorizer = func(tokenCandidate string) (*auth.Token, error) {
			client, err := app.Auth(context.Background())
			if err != nil {
				return nil, err
			}
			token, err := client.VerifyIDToken(context.Background(), tokenCandidate)
			if err != nil {
				return nil, err
			}
			if cfg.CheckEmailVerified && !token.Claims["email_verified"].(bool) {
				return nil, ErrEmailNotVerified
			}
			return token, nil
		}
	}
	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = ConfigDefault.ErrorHandler
	}
	if cfg.TokenExtractor == nil {
		cfg.TokenExtractor = ConfigDefault.TokenExtractor
	}
	if cfg.TokenCallback == nil {
		cfg.TokenCallback = func(ctx *fiber.Ctx, token *auth.Token) error {
			ctx.Locals(cfg.ContextKey, User{
				Email:         token.Claims["email"].(string),
				EmailVerified: token.Claims["email_verified"].(bool),
				UserID:        token.Claims["user_id"].(string),
			})
			return nil
		}
	}

	return cfg
}
