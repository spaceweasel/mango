// Package mango is a routing object which implements the http.Handler interface.
// It has been designed primarily to facilitate the creation of HTTP based APIs
// and REST services.
//
// It uses a context per request approach which enables simplified handlers, as
// much of the boiler plate work is done for you. The Context object takes care
// of tasks such as serialization/deserialization, respecting the Content-Type
// and Accept headers to ensure responses match the request. You can create and
// add your custom content-type encoders to suit your needs.
//
// A radix-tree based routing system enables better response times and
// greater flexibility in the way routes are structured and added to the system.
// Parameter based patterns with constraint rules provide great flexibility, and
// custom rules can be added.
//
// Hooks and other mechanisms exist to enable customisation in accordance with
// your specific application, such as authentication, database repository
// injection.
//
// See the full documentation here https://github.com/spaceweasel/mango/wiki
//
package mango
