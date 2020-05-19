// +build webui

package http

import (
	"github.com/markbates/pkger"
	"net/http"
)

func init() {
	Static = http.FileServer(pkger.Dir("/ui/web/dist/"))
}
