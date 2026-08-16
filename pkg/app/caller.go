package app

import "context"

// CallerResolver names the principal behind a request, so a document can record
// who wrote it.
//
// An interface injected by the wrapping service, exactly like the project
// resolver: restcol is a library and does not own an authentication scheme, so
// it cannot know how a caller is represented. The service that authenticated
// the request is the only thing that does.
//
// Returning "" means "no principal available", which is recorded as an empty
// writer. That is deliberate: an unattributed document must be distinguishable
// from one written by a principal literally named "anonymous" or "system", and
// inventing a name here would make the two identical forever after.
type CallerResolver interface {
	// Caller returns the principal for this request, or "" if there is none.
	Caller(ctx context.Context) string
}

// CallerResolverFunc adapts a function to CallerResolver.
type CallerResolverFunc func(ctx context.Context) string

func (f CallerResolverFunc) Caller(ctx context.Context) string { return f(ctx) }

// callerFromCtx is the single place a writer is derived.
//
// Nil-safe on purpose. A deployment that has not wired a resolver keeps working
// and records no writer, rather than failing every write - attribution is a
// property worth having, not worth refusing data over. The empty value is then
// honest about what happened.
func (r *RestColServiceServerService) callerFromCtx(ctx context.Context) string {
	if r.callerResolver == nil {
		return ""
	}
	return r.callerResolver.Caller(ctx)
}
