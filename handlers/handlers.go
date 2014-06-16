package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/slogsdon/docker-build-service/build"
	"math/rand"
	"net/http"
	"time"
)

var r *rand.Rand

func init() {
	// Create and seed the generator.
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func Compile(w http.ResponseWriter, r *http.Request) {
	app, lang, code := getVars(r)

	build.Create(app, lang, code)

	buildresp, _ := build.Build(app)
	runresp, _ := build.Run(app)

	fullresp := build.FullResp{AppId: app, Build: buildresp, Run: runresp}

	json, err := json.Marshal(fullresp)
	if err == nil {
		fmt.Fprint(w, string(json))
	}
}

func getVars(r *http.Request) (app, lang string, code []byte) {
	if r.FormValue("lang") != "" {
		lang = r.FormValue("lang")
	} else {
		panic("no language given")
	}

	if r.FormValue("app_id") != "" {
		app = r.FormValue("app_id")
	} else {
		app = getId()
	}

	code = []byte(r.FormValue("code"))
	return
}

func getId() string {
	bytes := make([]byte, 16)
	for i := 0; i < 16; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return base64.StdEncoding.EncodeToString(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
