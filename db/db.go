package db

import (
	"log"

	"gopkg.in/mgo.v2"
)

var session *mgo.Session

func InitMongoDBConnection() error {
	mongoURL := "mongodb+srv://italosilva18:costa2013@cluster0.xpywwjs.mongodb.net/?retryWrites=true&w=majority" // Substitua com sua própria URL de conexão
	var err error
	session, err = mgo.Dial(mongoURL)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
		return err
	}
	return nil
}

func GetMongoDBSession() *mgo.Session {
	return session
}

func CloseMongoDBConnection() {
	if session != nil {
		session.Close()
	}
}
