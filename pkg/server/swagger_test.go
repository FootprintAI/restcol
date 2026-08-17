package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// The defect: swagger-ui-dist ships an initializer pointing at Swagger's demo
// API, and we embed the assets rather than running the configurator container
// that is supposed to rewrite it. A test that only asserted "200 OK" would
// have passed throughout - the wrong file was being served perfectly well.
func TestInitializerDoesNotServeThePetstoreDemo(t *testing.T) {
	rr := httptest.NewRecorder()
	swaggerUIHandler().ServeHTTP(rr, httptest.NewRequest(
		http.MethodGet, "/swaggerui/swagger-initializer.js", nil))

	if rr.Code != http.StatusOK {
		t.Fatalf("initializer: got %d, want 200", rr.Code)
	}
	body := rr.Body.String()
	if strings.Contains(body, "petstore.swagger.io") {
		t.Error("still serving the swagger-ui-dist default, which points at " +
			"petstore.swagger.io - the whole point of overriding this file")
	}
	if !strings.Contains(body, "openapiv2/restcol.swagger.json") {
		t.Error("initializer does not load restcol's own spec")
	}
}

// The URL has to be built from the page's location, not hardcoded. A deployment
// mounted at /reststore/ and one mounted at / are served the SAME file, and the
// server cannot tell them apart - the gateway strips the prefix before any Go
// handler sees the request.
func TestInitializerResolvesTheSpecRelativeToItsMount(t *testing.T) {
	rr := httptest.NewRecorder()
	swaggerUIHandler().ServeHTTP(rr, httptest.NewRequest(
		http.MethodGet, "/swaggerui/swagger-initializer.js", nil))
	body := rr.Body.String()

	if !strings.Contains(body, "window.location.pathname") {
		t.Error("the spec URL is not derived from the page location, so it " +
			"cannot be correct under both / and /reststore/")
	}
	// Without this, every "Try it out" goes to origin + /v1/... and misses by
	// exactly the prefix.
	if !strings.Contains(body, "requestInterceptor") {
		t.Error("no requestInterceptor: the spec declares no basePath, so " +
			"requests will not carry the mount prefix")
	}
}

// Everything except the one overridden file must still come from the agent's
// embed - the bundles, the CSS, the page itself. Getting this wrong would 404
// the whole UI while the initializer test above still passed.
//
// Not /swaggerui/index.html: http.FileServer canonicalises that to ./ with a
// 301, which is correct behaviour and would make this assert the wrong thing.
func TestOtherAssetsStillComeFromTheAgent(t *testing.T) {
	for _, path := range []string{
		"/swaggerui/",
		"/swaggerui/swagger-ui-bundle.js",
		"/swaggerui/swagger-ui.css",
	} {
		rr := httptest.NewRecorder()
		swaggerUIHandler().ServeHTTP(rr, httptest.NewRequest(http.MethodGet, path, nil))
		if rr.Code != http.StatusOK {
			t.Errorf("%s: got %d, want 200 - the UI is not being served", path, rr.Code)
		}
	}
}

// The page has to reference the initializer for the override to matter at all.
// If a future swagger-ui-dist inlines its config instead, this fix becomes a
// no-op and everything above still passes.
func TestThePageActuallyLoadsTheInitializer(t *testing.T) {
	rr := httptest.NewRecorder()
	swaggerUIHandler().ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/swaggerui/", nil))
	if !strings.Contains(rr.Body.String(), "swagger-initializer.js") {
		t.Error("index.html does not load swagger-initializer.js, so overriding " +
			"that file configures nothing")
	}
}
