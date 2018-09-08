package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/goki0524/gopackage/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

// templは１つのテンプレートを表す
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTPはHTTPリクエストを処理する
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションアドレス")
	flag.Parse() // フラグを解析
	// Gomniauthのセットアップ
	gomniauth.SetSecurityKey("55dfbg7iu2nb4uywevihjw4tuiyub34noilk")
	gomniauth.WithProviders(
		// TODO: facebookとgithubも追加する
		// ("クライアントID", "秘密の値", "リダイレクト先")
		google.New("42313837065-6h3dc1dfpthfa94bgln3i02oi1gumdfu.apps.googleusercontent.com", "A9XTv_XEUnExMjJnUct-Y_es", "http://localhost:8080/auth/callback/google"),
	)
	// アバターを取得する方法は２つある。UseAuthAvatar or UseGravatar
	// r := newRoom(UseAuthAvatar)
	r := newRoom(UseGravatar)
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)
	// チャットルーム開始
	go r.run()
	// Webサーバーを起動
	log.Println("Webサーバーを起動します。ポート:", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
