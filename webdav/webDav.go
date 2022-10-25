package webdav

import (
	"net/http"
	"os"

	"github.com/happyanran/walnut/common"
	"github.com/happyanran/walnut/model"
	"golang.org/x/net/webdav"
)

func WebDav(svcCtx *common.ServiceContext) {
	if svcCtx.Cfg.WebDavConf.Enable {
		os.MkdirAll(svcCtx.Cfg.WebDavConf.Data, 640)

		myWD := &webdav.Handler{
			FileSystem: webdav.Dir(svcCtx.Cfg.WebDavConf.Data),
			LockSystem: webdav.NewMemLS(),
		}

		http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			username, password, ok := req.BasicAuth()

			if !ok {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			user := model.User{
				UserName: username,
			}

			if err := user.UserFindByName(); err != nil {
				http.Error(w, "Wait a moment", http.StatusInternalServerError)
				return
			}

			if user.ID == 0 || !svcCtx.Utilw.PwdCheck(user.Password, password) {
				http.Error(w, "Wrong username or password", http.StatusUnauthorized)
				return
			}

			//TODO
			//user
			// if *flagReadonly {
			// 	switch req.Method {
			// 	case "PUT", "DELETE", "PROPPATCH", "MKCOL", "COPY", "MOVE":
			// 		http.Error(w, "WebDAV: Read Only!!!", http.StatusForbidden)
			// 		return
			// 	}
			// }

			myWD.ServeHTTP(w, req)
		})

		svcCtx.Log.Info("WebDav available.")

		if svcCtx.Cfg.HttpsConf.Enable {
			svcCtx.Log.Error(http.ListenAndServeTLS(svcCtx.Cfg.WebDavConf.Addr, svcCtx.Cfg.HttpsConf.Certfile, svcCtx.Cfg.HttpsConf.Keyfile, nil))
		} else {
			svcCtx.Log.Error(http.ListenAndServe(svcCtx.Cfg.WebDavConf.Addr, nil))
		}
	}
}
