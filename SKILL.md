---
name: scalar-galaxy-go-sdk
description: "Go SDK for Scalar Galaxy API. Use when writing Go code that calls Scalar Galaxy API with the scalar-galaxy package: installing it, constructing and authenticating the client, and calling API operations."
---

# Scalar Galaxy Go SDK

Generated Go client for Scalar Galaxy API, published as `scalar-galaxy`. Use the generated client instead of hand-writing HTTP requests.

## Install

```sh
go get scalar-galaxy
```

## Client setup and authentication

```go
import (
	"context"
	"fmt"

	sdk "scalar-galaxy"
)

client := sdk.NewClient()
```

Provide credentials using the options below. Environment variables are read automatically when the target runtime supports them:

- `option.WithBearerAuth` (env: `BEARER_AUTH`) — JWT Bearer token authentication
- `option.WithBasicAuthUsername` (env: `BASIC_AUTH_USERNAME`) — Credential for the basicAuth_username client option.
- `option.WithBasicAuthPassword` (env: `BASIC_AUTH_PASSWORD`) — Credential for the basicAuth_password client option.
- `option.WithAPIKeyHeader` (env: `API_KEY_HEADER`) — API key request header
- `option.WithAPIKeyQuery` (env: `API_KEY_QUERY`) — API key query parameter
- `option.WithAPIKeyCookie` (env: `API_KEY_COOKIE`) — API key browser cookie

## Calling operations

```go
package main

import (
	"context"
	"fmt"
	"os"

	sdk "scalar-galaxy"
	"scalar-galaxy/option"
)

func main() {
	client := sdk.NewClient(
		option.WithBearerAuth(os.Getenv("BEARER_AUTH")),
	)

	planet, err := client.Planets.ListAllData(context.Background(), sdk.PlanetListAllDataParams{
		Limit: sdk.F[int64](10),
		Offset: sdk.F[int64](0),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(planet)
}
```

Method names, parameter shapes, and response types are generated from the API description — do not guess them. Look up the exact call signature in [api.md](./api.md) before writing a call.

## Error handling

Non-success responses return generated API errors. Error objects expose status, headers, response body, and request metadata where the target runtime supports it.

```go
planet, err := client.Planets.ListAllData(context.Background(), sdk.PlanetListAllDataParams{
if err != nil {
	var apiErr *sdk.Error
	if errors.As(err, &apiErr) {
		fmt.Println(apiErr.StatusCode, apiErr.RawJSON())
	}
	panic(err)
}

// imports: sdk "scalar-galaxy", "errors", "fmt"
```

## Requirements

- Go 1.22 or newer

## Reference files

- [README.md](./README.md) — full feature tour: client options, request options, retries and timeouts, logging.
- [api.md](./api.md) — complete catalogue of every operation with request and response types.
