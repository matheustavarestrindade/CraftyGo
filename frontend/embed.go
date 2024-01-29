package frontend

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:generate npm i
//go:generate npm run build
//go:embed all:build
var files embed.FS

func SvelteKitHandler(r *gin.Engine) {
	fsys, err := fs.Sub(files, "build")
	if err != nil {
		log.Fatal(err)
	}
	filesystem := http.FS(fsys)
	r.StaticFS("/", filesystem)
}

