package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/julienschmidt/httprouter"
	cli "github.com/levpay/surl/redis"
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
		Addr:       opt.Addr,
		Password:   opt.Password,
		DB:         opt.DB,
		MaxRetries: 3,
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
		log.Println("body, request: ", err)
		http.Error(w, err.Error(), 500)
		return
	}
	resp := &ResponseBody{}
	//struct request json format
	err = json.Unmarshal(b, resp)
	if err != nil {
		log.Println("json to struct: ", err)
		http.Error(w, err.Error(), 500)
		return
	}
	resp.Short, err = client.Find(resp.URL)
	if err == cli.ErrKeyNotFound {
		resp.Short, err = client.Set(resp.URL)
		if err != nil {
			log.Println("set redis client: ", err)
			http.Error(w, err.Error(), 500)
			return
		}
	} else if err != nil {
		log.Println("get redis client: ", err)
		http.Error(w, err.Error(), 500)
		return
	}
	e, err := json.Marshal(resp)
	if err != nil {
		log.Println("struct to json: ", err)
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
		log.Println("get redis client: ", err)
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Println("Found:", url)
	http.Redirect(w, req, url, 301)
}
