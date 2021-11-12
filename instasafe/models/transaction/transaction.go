package model

import (
	"gopkg.in/mgo.v2/bson"
	"instasafe/conn"
	"time"
)

// Transaction structure
type Transaction struct {
	ID        bson.ObjectId `bson:"_id"`
	Amount    float64       `bson:"amount"`
	Timestamp time.Time     `bson:"timestamp"`
}

// Transactions list
type Transactions []Transaction

// TransactionInfo model function
func TransactionInfo(id bson.ObjectId, transactionCollection string) (Transaction, error) {
	// Get DB from Mongo Config
	db := conn.GetMongoDB()
	transaction := Transaction{}
	err := db.C(transactionCollection).Find(bson.M{"_id": &id}).One(&transaction)
	return transaction, err
}
