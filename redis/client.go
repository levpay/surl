package redis

import (
	"errors"
	"log"
	"math/rand"
	u "net/url"

	"github.com/go-redis/redis"
	"github.com/rs/xid"
)

var (
	ErrKeyNotFound = errors.New("Redis does not contain key")
	ErrInvalidURL  = errors.New("Not a valid url")
)

type Client struct {
	cli *redis.Client
}

//Custom config
type Config struct {
	Addr       string
	Password   string
	DB         int
	MaxRetries int
}

//create redis client
func NewClient(config *Config) (*Client, error) {

	cli := redis.NewClient(&redis.Options{
		Addr:       config.Addr,
		Password:   config.Password,
		DB:         0, //DEFAULT
		MaxRetries: config.MaxRetries,
	})

	client := &Client{cli}
	pong, err := cli.Ping().Result()

	log.Println("pong:", pong)
	log.Println("error:", err)

	if err != nil {
		return nil, err
	}

	return client, nil
}

// find value pair by key
func (client *Client) Find(id string) (string, error) {

	cli := client.cli
	url, err := cli.Get(id).Result()

	// does not contain key
	if err == redis.Nil {
		return "", ErrKeyNotFound
	}

	// error
	if err != nil {
		return "", err
	}

	//found
	return url, nil
}

// set key-value pair
func (client *Client) Set(url string) (string, error) {

	//check validity of url
	_, err := u.ParseRequestURI(url)
	if err != nil {
		return "", ErrInvalidURL
	}

	cli := client.cli
	guid := xid.New()
	//decode value, shorten url
	val := guid.String()

	//set key-value to redis client
	err = cli.Set(val, url, 0).Err() //set no expire-time

	if err != nil {
		return "", err
	}

	return val, nil
}

//generate 6 byte slug - deprecate change to xid
func generateSlug() string {
	slug := make([]byte, 6)
	//base58
	var base = []byte("123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ")
	//generate slug
	for i := range slug {
		slug[i] = base[rand.Intn(len(base))] //base 58
	}
	key := string(slug[:])

	log.Println("generated key:", key)
	return key
}
