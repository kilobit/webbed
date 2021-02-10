WebbEd
======

Simple utilities for building Web Applications in Golang.

Status: In-Development

```
package main

import . "kilobit.ca/go/webbed"
import "kilobit.ca/go/server"
import "kilobit.ca/go/methods"
import "kilobit.ca/go/routes"
import "kilobit.ca/go/limits"

func main() {

	lg := log.New(os.Stdout, "- ", log.LstdFlags|log.Lmsgprefix)

	ctx := context.WithValue(context.Background(), SvcLoggerKey, lg)
	ctx = context.WithValue(ctx, SvcAppNameKey, "myapp")

	srv := MyService(ctx, ":8088")

	done := srv.Start()

	<-done
}

func MyService(ctx context.Context, addr string) *Server {

	lg, ok := ctx.Value(SvcLoggerKey).(*log.Logger)

	app := StringFromCtx(ctx, SvcAppNameKey)

	mh := methods.New(nil)

	rs := routes.New(http.NotFoundHandler())
	rs.Add("/" + app + "/fs/temp", mh.Get(http.FileServer("/tmp")))
	rs.Add("/" + app + "/api/",    mh.Post(MyHandler()))

	lh := limits.New(rs)

	srv := server.New(ctx, addr, lh, ServerOptLogger(lg))

	return srv
}

```

Sub-Packages include:

- server
- logger
- methods
- routes
- limits
- forms

Something else you need?  Submit an issue or better yet, a PR!

Features
--------

- Simple server management with cancellation / shutdown.
- Routing based on ShiftPath and Tries.
- Apply rate and other limits to inbound requests.
- Handle form data easily.
- Log requests and custom parameters in any format, including JSON.
- Load server parameters from the environment and propagate them
  through contexts.
- Use the parts you like, ignore the rest.

Installation
------------

```
go get kilobit.ca/go/webbed
go test -v
```

Building
--------

```
go get kilobit.ca/go/webbed
go build
```

Contribute
----------

Please help!  Submit pull requests through
[Github](https://github.com/kilobit/webbed).

Support
-------

Please submit issues through
[Github](https://github.com/kilobit/webbed).

License
-------

See LICENSE.

--  
Created: Aug 11, 2020  
By: Christian Saunders <cps@kilobit.ca>  
Copyright 2021 Kilobit Labs Inc.  
