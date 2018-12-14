package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/recoilme/slowpoke"
)

const storageFile = ".storage/sites.slowpoke"
const tokenSize = 40

type likeInfo struct {
	Liked bool   `json:"liked"`
	Count int    `json:"count"`
	Token string `json:"token"`
}

func makeToken() []byte {
	token := make([]byte, tokenSize)
	rand.Read(token)
	return token
}

func dumpToken(token []byte) string {
	return base64.StdEncoding.EncodeToString(token)
}

func loadToken(token string) []byte {
	result, _ := base64.StdEncoding.DecodeString(token)
	return result
}

func tokenInTokens(token []byte, tokens [][]byte) bool {
	for _, other := range tokens {
		if bytes.Equal(token, other) {
			return true
		}
	}
	return false
}

func split(buf []byte) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/tokenSize+1)
	for len(buf) >= tokenSize {
		chunk, buf = buf[:tokenSize], buf[tokenSize:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}

// SiteListHandler returns list of available sites
func SiteListHandler(w http.ResponseWriter, r *http.Request) {
	keys, _ := slowpoke.Keys(storageFile, nil, 0, 0, true)
	sep := []byte("\n")
	w.Write(bytes.Join(keys, sep))
}

// SiteAddHandler adds new site
func SiteAddHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slowpoke.Set(storageFile, []byte(vars["site"]), []byte(""))
	w.Write([]byte("OK!"))
}

// LikeInfoHandler returns info about likes for given site and post
func LikeInfoHandler(w http.ResponseWriter, r *http.Request) {
	// get or make token
	cookie, err := r.Cookie("token")
	var dumpedToken string
	var token []byte
	if err != nil {
		fmt.Print(err)
		fmt.Print("\n")
		token = makeToken()
		dumpedToken = dumpToken(token)
	} else {
		dumpedToken = cookie.Value
		token = loadToken(dumpedToken)
	}

	vars := mux.Vars(r)
	values, _ := slowpoke.Get(storageFile, []byte(vars["site"]))
	tokens := split(values)

	// set cookie
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie = &http.Cookie{Name: "token", Value: dumpedToken, Expires: expiration}
	http.SetCookie(w, cookie)

	// make response body
	info := likeInfo{
		Liked: tokenInTokens(token, tokens),
		Count: len(tokens),
		Token: dumpedToken,
	}
	data, _ := json.Marshal(info)
	w.Write(data)
}

func getRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/sites", SiteListHandler).Methods("GET")
	r.HandleFunc("/sites/{site:[a-zA-Z\\.]+}", SiteAddHandler).Methods("POST")
	r.HandleFunc("/likes/{site:[a-zA-Z\\.]+}/{post:[0-9]+}", LikeInfoHandler).Methods("GET")
	return r
}

func main() {
	defer slowpoke.CloseAll()
	r := getRouter()
	fmt.Print("Start!\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}
