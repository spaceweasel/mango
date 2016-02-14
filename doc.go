// Package mango is a routing object which implements the http.Handler interface.
// It has been designed primarily to facilitate the creation of HTTP based APIs
// and REST services.
//
// It uses a context per request approach which enables simplified handlers, as
// much of the boiler plate work is done for you. The Context object takes care
// of tasks such as serialization/deserialization, respecting the Content-Type
// and Accept headers to ensure responses match the request.
//
// A radix-tree based routing system enables better response times and
// greater flexiblility in the way routes are structured and added to the system.
//
// Hooks and other mechanisms exist to enable customisation in accordance with
// your specific application, such as authentication, database repository
// injection.
//
package mango
