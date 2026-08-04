# Scalar Galaxy Go API

Complete reference of every operation, grouped by resource. See [the README](./README.md) for usage and configuration.

## Contents

- [`Planets`](#planets)
  - [Get all planets](#get-all-planets)
  - [Create a planet](#create-a-planet)
  - [Get a planet](#get-a-planet)
  - [Update a planet](#update-a-planet)
  - [Delete a planet](#delete-a-planet)
  - [Upload an image to a planet](#upload-an-image-to-a-planet)
- [`CelestialBodies`](#celestialbodies)
  - [Create a celestial body](#create-a-celestial-body)
- [`Authentication`](#authentication)
  - [Create a user](#create-a-user)
  - [Get a token](#get-a-token)
  - [Get authenticated user](#get-authenticated-user)

## Setup

```go
import (
	"context"
	"fmt"

	sdk "scalar-galaxy"
)

client := sdk.NewClient()
```

## `Planets`

### Get all planets

It's easy to say you know them all, but do you really? Retrieve all the planets and check whether you missed one.

| Direction | Type |
| --- | --- |
| Request | [`PlanetListAllDataParams`](./planet.go) |
| Response | [`PlanetListAllDataResponse`](./planet.go) |

```go
planet, err := client.Planets.ListAllData(context.Background(), sdk.PlanetListAllDataParams{
	Limit: sdk.F[int64](10),
	Offset: sdk.F[int64](0),
})
if err != nil {
	panic(err)
}
fmt.Println(planet)
```

### Create a planet

Time to play god and create a new planet. What do you think? Ah, don't think too much. What could go wrong anyway?

| Direction | Type |
| --- | --- |
| Request | [`PlanetNewParams`](./planet.go) |
| Response | [`Planet`](./planet.go) |

```go
planet, err := client.Planets.New(context.Background(), sdk.PlanetNewParams{
	Planet: sdk.PlanetParam{
	Name: sdk.F[string]("Mars"),
},
})
if err != nil {
	panic(err)
}
fmt.Println(planet)
```

### Get a planet

You'll better learn a little bit more about the planets. It might come in handy once space travel is available for everyone.

| Direction | Type |
| --- | --- |
| Response | [`Planet`](./planet.go) |

```go
planet, err := client.Planets.Get(context.Background(), 1)
if err != nil {
	panic(err)
}
fmt.Println(planet)
```

### Update a planet

Sometimes you make mistakes, that's fine. No worries, you can update all planets.

| Direction | Type |
| --- | --- |
| Request | [`PlanetUpdateParams`](./planet.go) |
| Response | [`Planet`](./planet.go) |

```go
planet, err := client.Planets.Update(context.Background(), 1, sdk.PlanetUpdateParams{
	Planet: sdk.PlanetParam{
	Name: sdk.F[string]("Mars"),
},
})
if err != nil {
	panic(err)
}
fmt.Println(planet)
```

### Delete a planet

This endpoint was used to delete planets. Unfortunately, that caused a lot of trouble for planets with life. So, this endpoint is now deprecated and should not be used anymore.

```go
err := client.Planets.Delete(context.Background(), 1)
if err != nil {
	panic(err)
}
```

### Upload an image to a planet

Got a crazy good photo of a planet? Share it with the world!

| Direction | Type |
| --- | --- |
| Request | [`PlanetUploadImageParams`](./planet.go) |
| Response | [`PlanetUploadImageResponse`](./planet.go) |

```go
planet, err := client.Planets.UploadImage(context.Background(), 1, sdk.PlanetUploadImageParams{})
if err != nil {
	panic(err)
}
fmt.Println(planet)
```

## `CelestialBodies`

### Create a celestial body

Stars, moons, comets, the occasional rogue asteroid — if it glows or drifts through the void, you can add it here.

| Direction | Type |
| --- | --- |
| Request | [`CelestialBodyNewParams`](./celestialbody.go) |
| Response | [`CelestialBody`](./celestialbody.go) |

```go
celestialBody, err := client.CelestialBodies.New(context.Background(), sdk.CelestialBodyNewParams{
	CelestialBody: sdk.PlanetParam{
	Name: sdk.F[string]("Mars"),
},
})
if err != nil {
	panic(err)
}
fmt.Println(celestialBody)
```

## `Authentication`

### Create a user

Time to create a user account, eh?

| Direction | Type |
| --- | --- |
| Request | [`AuthenticationNewUserParams`](./authentication.go) |
| Response | [`User`](./authentication.go) |

```go
authentication, err := client.Authentication.NewUser(context.Background(), sdk.AuthenticationNewUserParams{
	Email: sdk.F[string]("marc@scalar.com"),
	Password: sdk.F[string]("i-love-scalar"),
	Name: sdk.F[string]("Marc"),
})
if err != nil {
	panic(err)
}
fmt.Println(authentication)
```

### Get a token

Yeah, this is the boring security stuff. Just get your super secret token and move on.

| Direction | Type |
| --- | --- |
| Request | [`AuthenticationNewTokenParams`](./authentication.go) |
| Response | [`AuthenticationNewTokenResponse`](./authentication.go) |

```go
authentication, err := client.Authentication.NewToken(context.Background(), sdk.AuthenticationNewTokenParams{
	Credentials: sdk.CredentialsParam{
	Email: sdk.F[string]("marc@scalar.com"),
	Password: sdk.F[string]("i-love-scalar"),
},
})
if err != nil {
	panic(err)
}
fmt.Println(authentication)
```

### Get authenticated user

Find yourself they say. That's what you can do here.

| Direction | Type |
| --- | --- |
| Response | [`User`](./authentication.go) |

```go
authentication, err := client.Authentication.ListMe(context.Background())
if err != nil {
	panic(err)
}
fmt.Println(authentication)
```
