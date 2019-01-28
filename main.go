package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	cli "github.com/levpay/surl/redis"
	redis "gopkg.in/redis.v5"
)

var client *cli.Client

type ResponseBody struct {
	URL   string `json:"url"`
	Short string `json:"short"`
}

func init() {
	var err error
	REDIS_URL := os.Getenv("REDIS_URL")
	if REDIS_URL == "" {
		REDIS_URL = "redis://127.0.0.1"
	}
	opt, err := redis.ParseURL(REDIS_URL)
	client, err = cli.NewClient(&cli.Config{
		Addr:     opt.Addr,
		Password: opt.Password,
		DB:       opt.DB,
	})
	if err != nil {
		fmt.Println("err:", err)
		os.Exit(0)
	}
}

func main() {
	//http router config
	router := httprouter.New()
	router.GET("/", handeIndex)
	router.POST("/", handleSet)
	router.GET("/:slug", handleFind)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}

//Serve static files
func handeIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(201)
	w.Write([]byte("Oops, sorry not found!"))
}

//Hande url post request
func handleSet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//parse request body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	resp := &ResponseBody{}
	//struct request json format
	err = json.Unmarshal(b, resp)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	resp.Short, err = client.Set(resp.URL)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	e, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write([]byte(e))
}

//Hande get request
func handleFind(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	//get params
	slug := ps.ByName("slug")
	// get url from  redis client
	url, err := client.Find(slug)
	if err == cli.ErrKeyNotFound {
		http.Error(w, err.Error(), 404)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Println("Found:", url)
	http.Redirect(w, req, url, 301)
}
