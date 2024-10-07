package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type WordsOut struct {
	Page  string   `json:"page"`
	Input string   `json:"input"`
	Words []string `json:"words"`
}

type OccurenceOut struct {
	Page string         `json:"page"`
	Freq map[string]int `json:"freq"`
}

type LoginRequest struct {
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type database struct {
	words       []string
	password    string
	tokenSecret []byte
}

func (db *database) insert_handler(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("input")
	if input != "" {
		db.words = append(db.words, input)
	}

	output := WordsOut{
		Page:  "words",
		Input: input,
		Words: db.words,
	}

	out, err := json.Marshal(output)
	if err != nil {

		fmt.Println("Error in marhsalling")
		return
	}

	fmt.Fprint(w, string(out))
}

func (db *database) occurence_handler(w http.ResponseWriter, r *http.Request) {
	mp := make(map[string]int)
	for _, v := range db.words {
		if _, ok := mp[v]; !ok {
			mp[v] = 1
		} else {
			mp[v]++
		}
	}

	occurenceout := OccurenceOut{
		Page: "occurence",
		Freq: mp,
	}

	out, err := json.Marshal(occurenceout)
	if err != nil {
		log.Print(err)
	}
	fmt.Fprint(w, string(out))
}

func (db *database) login_handler(w http.ResponseWriter, r *http.Request) {
	var (
		body          []byte
		err           error
		loginrequest  LoginRequest
		signedtoken   string
		loginresponse LoginResponse
	)

	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Not a Post Request")
		return
	}

	if db.password == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No password set for the db")
		return
	}

	if body, err = io.ReadAll(r.Body); err != nil {
		fmt.Fprintf(w, "ReadAll error")
		return
	}

	if err = json.Unmarshal(body, &loginrequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unmarshal error cannot get client password")
		return
	}

	if loginrequest.Password != db.password {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Wrong Password")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nbf": time.Now().Unix(),
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	})

	if signedtoken, err = token.SignedString(db.tokenSecret); err != nil {
		log.Fatalf("signing token error: %v", err)
	}

	loginresponse = LoginResponse{
		Token: signedtoken,
	}

	if err = json.NewEncoder(w).Encode(loginresponse); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatalf("Token Encoding Error: %v", err)
	}
}

func getsecret() []byte {
	b := make([]byte, 30)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("rand read err: %v", err)
	}
	return b
}

func (db *database) middleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if db.password != "" {
			if r.Header.Get("Authorization") == "" {
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprint(w, "Authorization not set")
				return
			}

			tokenstring := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", -1)
			_, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected Signing Method: %v", ok)
				}
				return db.tokenSecret, nil
			})
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprintf(w, "Authorization token invalid: %v", err)
				return
			}
		}
		next(w, r)
	})
}

func main() {
	var password *string = flag.String("password", "", "Password to Protect API")
	flag.Parse()

	db := &database{
		words:       []string{},
		password:    *password,
		tokenSecret: getsecret(),
	}

	http.HandleFunc("/put", db.middleware(db.insert_handler))
	http.HandleFunc("/login", db.middleware(db.login_handler))
	http.HandleFunc("/occur", db.occurence_handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
