// Replaces the swagger-ui-dist default, which ships pointing at Swagger's own
// public demo API and carries a comment saying the line "will be replaced by
// docker/configurator, when it runs in a docker-container". We embed the assets
// instead of running that container, so nothing ever replaced it: every
// deployment served the demo API in place of ours, and on an air-gapped or
// internal-only site it could not even fetch that, so the page rendered an
// error.
//
// The test for this file greps for the demo host by name, so do not write it
// here - not even in a comment.
//
// Everything here is computed IN THE BROWSER, from the URL the page was served
// on. That is deliberate. The server cannot know where it is mounted: a gateway
// that routes /reststore/ to us rewrites the prefix away before the request
// arrives, so by the time any Go handler sees it the path is /swaggerui/ and
// the prefix is gone. The browser is the only party that still knows.
window.onload = function () {
  // "/reststore/" behind a path-prefix gateway, "/" when run directly.
  const mount = window.location.pathname.replace(/swaggerui\/.*$/, "");

  window.ui = SwaggerUIBundle({
    // Relative to the mount, so this file needs no knowledge of the prefix.
    url: mount + "openapiv2/restcol.swagger.json",
    dom_id: "#swagger-ui",
    deepLinking: true,
    presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
    plugins: [SwaggerUIBundle.plugins.DownloadUrl],
    layout: "StandaloneLayout",

    // The spec declares neither `host` nor `basePath`, so swagger-ui builds
    // request URLs as origin + the path as written - https://host/v1/projects/…
    // - which behind a prefix misses by exactly the prefix and 404s every
    // "Try it out". Putting a basePath in the spec would hardcode one
    // deployment's layout into an artifact shared by all of them, so the fix
    // belongs here, where the real mount point is known.
    requestInterceptor: function (req) {
      if (mount === "/" || !mount) {
        return req;
      }
      const prefix = mount.replace(/\/$/, "");
      const u = new URL(req.url, window.location.origin);
      if (u.origin === window.location.origin && !u.pathname.startsWith(prefix + "/")) {
        u.pathname = prefix + u.pathname;
        req.url = u.toString();
      }
      return req;
    },
  });
};
