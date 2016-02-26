# Mango [![Build Status](https://travis-ci.org/spaceweasel/mango.svg?branch=master)](https://travis-ci.org/spaceweasel/mango) [![Coverage Status](http://codecov.io/github/spaceweasel/mango/coverage.svg?branch=master)](http://codecov.io/github/spaceweasel/mango?branch=master) [![GoDoc](http://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/github.com/spaceweasel/mango) [![MIT](https://img.shields.io/npm/l/express.svg)](https://github.com/spaceweasel/mango/blob/master/LICENSE)

Mango is a routing package designed to simplify the development of web service code in Golang. The Router object implements the standard library's http.Handler interface, so it can be used with the http.ListenAndServe method.

Mango uses a context per request approach which enables simplified handlers, as
much of the boiler plate work is done for you. The Context object takes care
of tasks such as serialization/deserialization, respecting the Content-Type
and Accept headers to ensure responses match the request. You can add your own custom content-type encoders if required.

A radix-tree based routing system enables better response times and greater flexibility in the way routes are structured and added to the system.

Hooks and other mechanisms exist to enable customization in accordance with your specific application, such as authentication, database repository injection.

Detailed documentation can be found [here](https://github.com/spaceweasel/mango/wiki).  

A *Hello World* example:

```go
package main

import (
  "net/http" 	

  "github.com/spaceweasel/mango"
)

func main() {
  // get a new router instance
  router := mango.NewRouter()

  // register a GET handler function
  router.Get("/hello", hello)

  // assign the router as the main handler
  http.ListenAndServe(":8080", router)
}

// hello handler function
func hello(c *mango.Context) {
  c.RespondWith("Hello world!")
}
```

#### TODOs
- [x] Add methods to allow custom encoders to be added
- [x] Add methods to allow custom route parameter validators to be added
- [ ] Add more documentation
- [x] Add OPTIONS handler for CORS support
- [x] Add OPTIONS methods to test browser
- [ ] Add more route parameter validators
