# RoadRunner GLOW

Simple http request Orchestrate for RoadRunner.

## Installation

### QuickBuild

Add to `.build.json` package `github.com/cjp2600/roadrunner-glow` and register it as `rr.Container.Register(glow.ID, &glow.Service{})`

After it build RR using QuickBuild.

Example of final file:
```json
{
  "packages": [
    "github.com/spiral/roadrunner/service/env",
    "github.com/spiral/roadrunner/service/http",
    "github.com/spiral/roadrunner/service/rpc",
    "github.com/spiral/roadrunner/service/static",
    "github.com/cjp2600/roadrunner-glow"
  ],
  "commands": [
    "github.com/spiral/roadrunner/cmd/rr/http"
  ],
  "register": [
    "rr.Container.Register(env.ID, &env.Service{})",
    "rr.Container.Register(rpc.ID, &rpc.Service{})",
    "rr.Container.Register(http.ID, &http.Service{})",
    "rr.Container.Register(static.ID, &static.Service{})",
    "rr.Container.Register(glow.ID, &glow.Service{})"
  ]
}
```

### Manual

1. Add dependency by running `go get github.com/cjp2600/roadrunner-glow`

2. Add to `rr/main.go` import `github.com/cjp2600/roadrunner-glow`

3. Add to `rr/main.go` line `rr.Container.Register(glow.ID, &glow.Service{})` after `rr.Container.Register(http.ID, &http.Service{})`

Final file should look like this:
```go
package main

import (
	rr "github.com/spiral/roadrunner/cmd/rr/cmd"

	// services (plugins)
	"github.com/spiral/roadrunner/service/env"
	"github.com/spiral/roadrunner/service/http"
	"github.com/spiral/roadrunner/service/rpc"
	"github.com/spiral/roadrunner/service/static"
	"github.com/cjp2600/roadrunner-glow"

	// additional commands and debug handlers
	_ "github.com/spiral/roadrunner/cmd/rr/http"
)

func main() {
	rr.Container.Register(env.ID, &env.Service{})
	rr.Container.Register(rpc.ID, &rpc.Service{})
	rr.Container.Register(http.ID, &http.Service{})
	rr.Container.Register(static.ID, &static.Service{})

    // register custom services
	rr.Container.Register(glow.ID, &glow.Service{})

	// you can register additional commands using cmd.CLI
	rr.Execute()
}
```

### PHP USAGE [PHP-GLOW](https://github.com/cjp2600/php-glow)


