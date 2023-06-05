# Go Fiber Firebase Auth Middleware

> **Note**: This is a fork of the original [gofiber-firebaseauth](https://github.com/sacsand/gofiber-firebaseauth)
> middlware for the RALF project.
> The API changed in this fork.

Authenticate your endpoints with [Firebase Authentication](https://github.com/LeafyCode/express-firebase-auth/).

gofiber-firebaseauth is inspired by npm
package [express-firebase-auth](https://github.com/LeafyCode/express-firebase-auth/) .

```sh
$ go get -u github.com/gofiber/fiber/v2
$ go get github.com/ralf-life/gofiber-firebaseauth
```

# Features

- Authenticate the user using Firebase before running the function.
- Ability to skip authentication on public API endpoints.

# Usage

## Configuration

```go
app.Use(gofiberfirebaseauth.New(firebaseApp, gofiberfirebaseauth.Config{
	// ...
}))
```

> **Note**: You can find all available configuration parameters [here](config.go).

## Accessing the User

```go
func Handler(ctx *fiber.Ctx) error {
    currentUser := ctx.Locals("user").(gofiberfirebaseauth.User)
	return ctx.SendString(fmt.Sprintf("Hello, %s!", currentUser.UserID))
}
```

**Customizing the User object**

You can specify a `TokenCallback` function in the config, to create a custom user object.

```go
cfg := gofiberfirebaseauth.Config{
    TokenCallback: func(ctx *fiber.Ctx, token *auth.Token) error {
        ctx.Locals("raw-token", token)
        return nil
    }
}
```

---

## Extracting the Token

By default, the middleware checks the `Authorization` header for tokens.

If you have a custom token format (e.g. a `Bearer ...` prefix), or if you want to read token from cookies,
you can specify a custom `TokenExtractor` in the config.

**Extract Token from Headers**

```go
cfg := gofiberfirebaseauth.Config{
	TokenExtractor: gofiberfirebaseauth.NewHeaderExtractor("Bearer ")
}
```

**Extract Token from Cookies**

```go
cfg := gofiberfirebaseauth.Config{
    TokenExtractor: gofiberfirebaseauth.NewCookieExtractor("my-cookie-name", "Bearer ")
}
```

**Extract Token from Header or Cookie**

```go
cfg := gofiberfirebaseauth.Config{
    TokenExtractor: gofiberfirebaseauth.NewMultiExtractor(
        gofiberfirebaseauth.NewHeaderExtractor("Bearer "),
        gofiberfirebaseauth.NewCookieExtractor("my-cookie-name", "Bearer "),
    )
}
```

## License

[MIT licensed](./LICENSE).
