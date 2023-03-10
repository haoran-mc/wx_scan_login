package sessions

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/haoran-mc/wx_scan_login/web-back-end/pkg/config"
)

var filesystemStore *sessions.FilesystemStore

var wxLoginSessionName = config.Conf.Session.Name
var wxLoginSessionPath = config.Conf.Session.Path
var wxLoginSessionSecret = config.Conf.Session.Secret

func init() {
	filesystemStore = sessions.NewFilesystemStore(
		wxLoginSessionPath,
		[]byte(wxLoginSessionSecret),
	)
	filesystemStore.Options = &sessions.Options{
		MaxAge: 60 * 5, // 设置过期时间，5 分钟过期
	}
}

func GetSession(r *http.Request) *sessions.Session {
	session, _ := filesystemStore.Get(r, wxLoginSessionName)
	return session
}
