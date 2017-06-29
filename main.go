package main

import (
	"net/http"
	"log"
	"sync"
	"html/template"
	"path/filepath"
	"flag"
	"trace"
	"os"
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
	t.temp1.Execute(w, r) // 戻り値はチェックするべき Request情報を添付する
}

func main() {

	//ポート番号
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()// フラグを解釈する

	r := newRoom()
	r.traer = trace.New(os.Stdout)

	//ルート
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	go r.run()

	//webサーバを開始します
	log.Println("Webサーバを起動します。ポート番号: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
