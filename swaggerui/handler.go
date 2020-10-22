package swaggerui

import (
	newswaggerui "github.com/utilitywarehouse/swaggerui"
	"net/http"
)

const (

	// basePath is the base path of the embedded swagger-ui.
	basePath = "/swagger-ui/"
	// specFile is the name of the swagger JSON file to serve.
	specFile = "/swagger.json"
)

// Handler returns an HTTP handler that serves the
// swagger-ui under /swagger-ui/
func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", redirect)
	mux.Handle(basePath, UIHandler())
	mux.Handle(specFile, FileHandler())
	return mux
}

// FileHandler returns an HTTP handler that serves the
// swagger.json file
// Deprecated: https://github.com/utilitywarehouse/swaggerui/blob/55a2f44cd9d9ae0d237f99e111eda1059614c4d8/handler.go#L21
func FileHandler() http.Handler {
	return newswaggerui.SwaggerFile()
}

// UIHandler returns an HTTP handler that serves the
// swagger UI
// Deprecated: https://github.com/utilitywarehouse/swaggerui/blob/55a2f44cd9d9ae0d237f99e111eda1059614c4d8/handler.go#L12
func UIHandler() http.Handler {
	return http.StripPrefix("/swagger-ui/", newswaggerui.SwaggerUI())
}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, basePath, http.StatusSeeOther)
}