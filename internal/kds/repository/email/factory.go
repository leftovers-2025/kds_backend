package email

import (
	"os"
	"strconv"
)

func EmailFromEnv() *EmailAuth {
	host, ok := os.LookupEnv("KDS_EMAIL_HOST")
	if !ok {
		panic("\"KDS_EMAIL_HOST\" is not set")
	}
	strPort, ok := os.LookupEnv("KDS_EMAIL_PORT")
	if !ok {
		panic("\"KDS_EMAIL_PORT\" is not set")
	}
	port, err := strconv.Atoi(strPort)
	if err != nil {
		panic("\"KDS_EMAIL_PORT\" is invalid")
	}
	address, ok := os.LookupEnv("KDS_EMAIL_ADDRESS")
	if !ok {
		panic("\"KDS_EMAIL_ADDRESS\" is not set")
	}
	password, ok := os.LookupEnv("KDS_EMAIL_PASSWORD")
	if !ok {
		panic("\"KDS_EMAIL_PASSWORD\" is not set")
	}
	auth, err := NewEmailAuth(host, port, address, password)
	if err != nil {
		panic(err.Error())
	}
	return auth
}
