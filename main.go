package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"strconv"

	"github.com/ory/viper"
)

var t *template.Template

type EmailMessage struct {
	From string
	Body []byte
	To   []string
}

type EmailCredentials struct {
	Username, Password, Server string
	Port                       int
}

type User struct{
	Name,Surname,Birthday  string
} 

type Config struct {
	From      string        `mapstructure:"FROM"`
	Password      string        `mapstructure:"PASSWORD"`
	To      []string        `mapstructure:"TO"`
	Name   string        `mapstructure:"NAME"`
	Surname   string        `mapstructure:"SURNAME"`
	Birthday   string        `mapstructure:"BIRTHDAY"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func main() {

	config,err:=LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}
	
	u :=&User{
		Name: config.Name,
		Surname: config.Surname,
		Birthday: config.Birthday,
	}
	t = template.Must(template.ParseFiles("templates/s.html"))

	var body bytes.Buffer
	t.Execute(&body, u)

	fmt.Printf("%s", body)
	message := &EmailMessage{
		From:    config.From,
		To:      config.To,
		Body:    []byte(
						"Subject: Email with SMTP package!\r\n" +
						"MIME: MIME-version: 1.0\r\n" +
						"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
    					"\r\n"+body.String(),
					),
	}		
	authCreds := &EmailCredentials{
		Username: config.From,
		Password: config.Password,
		Server:   "smtp.gmail.com",
		Port:     587,
	}

	auth := smtp.PlainAuth("",
		authCreds.Username,
		authCreds.Password,
		authCreds.Server,
	)

	err=smtp.SendMail(authCreds.Server+":"+strconv.Itoa(authCreds.Port),
		auth,
		message.From,
		message.To,
		message.Body)
    if err != nil {
      log.Println(err)
      return
    }
}
