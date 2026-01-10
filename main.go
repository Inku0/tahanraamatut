package main

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/maddalax/htmgo/framework/config"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/service"

	"tahanraamatut/__htmgo"
	"tahanraamatut/internal/api"
)

func main() {
	locator := service.NewLocator()
	service.Set(locator, "readarr", func() *api.ReadarrService {
		return api.NewReadarrService()
	})
	cfg := config.Get()

	h.Start(h.AppOpts{
		ServiceLocator: locator,
		LiveReload:     true,
		Register: func(app *h.App) {
			sub, err := fs.Sub(GetStaticAssets(), "assets/dist")

			if err != nil {
				panic(err)
			}

			http.FileServerFS(sub)

			// change this in htmgo.yml (public_asset_path)
			app.Router.Handle(fmt.Sprintf("%s/*", cfg.PublicAssetPath),
				http.StripPrefix(cfg.PublicAssetPath, http.FileServerFS(sub)))

			__htmgo.Register(app.Router)
		},
	})
}
