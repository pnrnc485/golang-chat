package main

import (
	"net/http"
	"log"
	"sync"
	"html/template"
	"path/filepath"
	"flag"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

// temp1は一つのテンプレートを表します
type templateHandler struct {
	once		sync.Once
	filename	string
	temp1		*template.Template
}

//ServeHTTPはHTTPリクエストを処理します
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.temp1 = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	data := map[string]interface{} {
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		userData := objx.MustFromBase64(authCookie.Value)
		log.Println( userData)
		data["UserData"] = userData
	}

	t.temp1.Execute(w, data) // 戻り値はチェックするべき Request情報を添付する
}

func main() {

	//ポート番号
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()// フラグを解釈する

	//Gomniauthのセットアップ
	gomniauth.SetSecurityKey(securityKey)
	gomniauth.WithProviders(
		google.New(googleClientId, googleSecret,"http://localhost:8080/auth/callback/google"),
	)

	r := newRoom(UsefileSystemAvatar)

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}) )
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploadHandler)
	http.Handle("/room", r)

	//ログアウトするために、クッキーを削除する
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name: "auth",
			Value: "",
			Path: "/",
			MaxAge: -1,//-1にするとブラウザ上のクッキーが即座に削除される
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	go r.run()

	//webサーバを開始します
	log.Println("Webサーバを起動します。ポート番号: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
