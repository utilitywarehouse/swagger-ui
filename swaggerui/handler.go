//go:generate sh -c "go-bindata -pkg=internal -o internal/files.go `find swagger-ui -type d`"

package swaggerui

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/utilitywarehouse/swagger-ui/assetfs"
	"github.com/utilitywarehouse/swagger-ui/swaggerui/internal"
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
func FileHandler() http.Handler {
	data := swaggerFile()
	return &handler{modTime: time.Now(), body: data}
}

// UIHandler returns an HTTP handler that serves the
// swagger UI
func UIHandler() http.Handler {
	as := &assetfs.AssetStore{
		Names: internal.AssetNames,
		Data:  internal.Asset,
		Info:  internal.AssetInfo,
	}
	fs, err := assetfs.New(as)
	if err != nil {
		panic(fmt.Sprintf("failed to create static fs: %v", err))
	}
	return http.FileServer(http.FileSystem(fs))
}

type handler struct {
	modTime time.Time
	body    io.ReadSeeker
}

func (f *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeContent(w, r, specFile, f.modTime, f.body)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, basePath, http.StatusSeeOther)
}

func swaggerFile() io.ReadSeeker {
	e, err := os.Executable()
	if err != nil {
		panic(err)
	}

	d := filepath.Dir(e)
	swaggerPath := fmt.Sprintf("%v/swagger.json", d)
	_, err = os.Stat(swaggerPath)
	if err != nil {
		log.Println("swagger.json not found in project root:", swaggerPath)
		swaggerPath = "/tmp/swagger.json"
		_, tmpErr := os.Stat(swaggerPath)
		if tmpErr != nil {
			log.Println("swagger.json not found in /tmp/swagger.json")
			swaggerPath = "../swagger.json"
			_, parentErr := os.Stat(swaggerPath)
			if parentErr != nil {
				log.Println("swagger.json not found in ../swagger.json")
				swaggerPath = os.Getenv("SWAGGER_JSON_PATH")
				_, envErr := os.Stat(swaggerPath)
				if envErr != nil {
					panic("Can not find swagger.json file")
				}
			}
		}
	}

	fi, err := os.Open(swaggerPath)
	if err != nil {
		panic(fmt.Sprintf("unable to open swagger.json: %v", err))
	}
	return fi
}
