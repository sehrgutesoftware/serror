// package serror handles serialization and deserialization of errors.
//
// It can be used to restore the original error type so that it passes eg.
// the `errors.Is()` check after it has been sent over the wire.
//
// This is useful for service interfaces that can be used both as part of
// a monolith and as a microservice.
package serror
