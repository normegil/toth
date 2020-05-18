// +build webui

package router

import (
	"github.com/markbates/pkger"
	"net/http"
)

func init() {
	static = http.FileServer(pkger.Dir("/../ui/web/dist"))
}
