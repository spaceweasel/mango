# Mango [![Build Status](https://travis-ci.org/spaceweasel/mango.svg?branch=master)](https://travis-ci.org/spaceweasel/mango) [![Coverage Status](http://codecov.io/github/spaceweasel/mango/coverage.svg?branch=master)](http://codecov.io/github/spaceweasel/mango?branch=master) [![GoDoc](http://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/github.com/spaceweasel/mango) [![MIT](https://img.shields.io/npm/l/express.svg)](https://github.com/spaceweasel/mango/blob/master/LICENSE)

Mango is esentially a routing package designed to simplify the development of web service code in Golang. The Router object implements the standard library's http.Handler interface, so it can be used with the http.ListenAndServe method.

Mango uses a context per request approach which enables simplified handlers, as
much of the boiler plate work is done for you. The Context object takes care
of tasks such as serialization/deserialization, respecting the Content-Type
and Accept headers to ensure responses match the request. You can add your own custom content-type encoders if required.

A radix-tree based routing system enables better response times and greater flexibility in the way routes are structured and added to the system.

Hooks and other mechanisms exist to enable customization in accordance with your specific application, such as authentication, database repository injection.

Mango includes many features to speed up your webservice development, including simple CORS setup, a customizable validation system for your routes and models (with several validators built in), plus an easy to use *test browser* to enable  end-to-end simulation testing.

### A *Hello World* example:

```go
package main

import (
  "net/http" 	

  "github.com/spaceweasel/mango"
)

func main() {
  // get a new router instance
  r := mango.NewRouter()

  // register a GET handler function
  r.Get("/hello", hello)

  // assign the router as the main handler
  http.ListenAndServe(":8080", r)
}

// hello handler function
func hello(c *mango.Context) {
  c.RespondWith("Hello world!")
}
```
### Documentation

* [Home](https://github.com/spaceweasel/mango/wiki)
* [Getting Started](https://github.com/spaceweasel/mango/wiki/getting-started)
* [Routing](https://github.com/spaceweasel/mango/wiki/routing)  
  * [Registration Methods](https://github.com/spaceweasel/mango/wiki/registration-methods)
  * [Routing Patterns](https://github.com/spaceweasel/mango/wiki/routing-patterns)
  * [Handler Functions](https://github.com/spaceweasel/mango/wiki/handler-functions)
  * [Organizing Handlers](https://github.com/spaceweasel/mango/wiki/organizing-handlers)
* [The Context Object](https://github.com/spaceweasel/mango/wiki/context-object)
  * [Response Helpers](https://github.com/spaceweasel/mango/wiki/response-helpers)
  * [Identity and Authenticated](https://github.com/spaceweasel/mango/wiki/identity-auth)
  * [Model Binding](https://github.com/spaceweasel/mango/wiki/model-binding)
  * [Model Validation](https://github.com/spaceweasel/mango/wiki/model-validation)
* [Logging](https://github.com/spaceweasel/mango/wiki/logging)
* [CORS](https://github.com/spaceweasel/mango/wiki/cors)
* [PreHooks and PostHooks](https://github.com/spaceweasel/mango/wiki/pre-post-hooks)
* [Handler Wrapping](https://github.com/spaceweasel/mango/wiki/handler-wrapping)
* [Encoders](https://github.com/spaceweasel/mango/wiki/encoders)
* [Validators](https://github.com/spaceweasel/mango/wiki/validators)
* [Testing (Mango Browser)](https://github.com/spaceweasel/mango/wiki/testing)
