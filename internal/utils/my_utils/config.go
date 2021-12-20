package my_utils

import (
	"log"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/rbcervilla/redisstore"
	"gopkg.in/gomail.v2"
)

// connect nats
func ConnectNats() *nats.Conn {

	nc, err := nats.Connect("YOUR_NATS_SERVER:4222")
	if err != nil {
		log.Fatal(err)
	}
	return nc
}

func SendMail(message string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "EMAIL_ADDRESS")
	m.SetHeader("To", "EMAIL_ADDRESS")
	m.SetHeader("Subject", "Dablopipe contact form")
	m.SetBody("text/html", message)

	d := gomail.NewDialer("YOUR_SMTP_SERVER:PORT", 465, "EMAIL_ADDRESS", "YOUR_PASSWORD")

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func GetCookieStore() *redisstore.RedisStore {

	client := redis.NewClient(&redis.Options{
		Addr: "YOUR_REDIS_SERVER:PORT",
	})

	// New default RedisStore
	store, err := redisstore.NewRedisStore(client)
	if err != nil {
		log.Fatal("failed to create redis store: ", err)
	}
	return store
}
