package gofiberfirebaseauth

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

type ExtractorFun func(ctx *fiber.Ctx) string

func NewMultiExtractor(extractors ...ExtractorFun) ExtractorFun {
	return func(ctx *fiber.Ctx) (token string) {
		for _, e := range extractors {
			if token = e(ctx); token != "" {
				break
			}
		}
		return
	}
}

func stripPrefix(str string, prefixes ...string) string {
	for _, s := range prefixes {
		str = strings.TrimPrefix(str, s)
	}
	return str
}

func NewHeaderExtractor(stripLeft ...string) ExtractorFun {
	return func(ctx *fiber.Ctx) string {
		return stripPrefix(ctx.Get(fiber.HeaderAuthorization), stripLeft...)
	}
}

func NewCookieExtractor(cookieKey string, stripLeft ...string) ExtractorFun {
	return func(ctx *fiber.Ctx) string {
		return stripPrefix(ctx.Cookies(cookieKey), stripLeft...)
	}
}
