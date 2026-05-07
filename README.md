# Route Pattern
## Static
// match /users
app.Get("/users")

## Param
// match /users/abc
// match /users/123
app.Get("/users/:id")
app.Get("/users/{id}")

regex is not supported.

## Catch-All (wildcard)
// match /users/abc
// match /users/abc/def/ghi
app.Get("/users/*others")
app.Get("/users/{others...}")

// this is not supported
app.Get("/users/*")
