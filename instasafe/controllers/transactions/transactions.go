package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"instasafe/conn"
	transaction "instasafe/models/transaction"
	"net/http"
	"time"
)

// TransactionCollection statically declared
const TransactionCollection = "transaction"

var (
	errNotExist           = errors.New("No Transactions Found")
	errInvalidID          = errors.New("Invalid ID")
	errInvalidTransaction = errors.New("Invalid Transaction")
	errExpired            = errors.New("Transaction Expired")
	errInsertionFailed    = errors.New("Error In Transaction Insertion")
	errDeletionFailed     = errors.New("Error in Transaction Deletion")
)

// GetAllTransaction Endpoint
func GetAllTransaction(c *gin.Context) {
	// Get DB from Mongo Config
	db := conn.GetMongoDB()
	transactions := transaction.Transactions{}

	query := `{
		timestamp: { // 1 minutes ago (from now)
			$gt: new Date(ISODate().getTime() - 1000 * 60 * 1)
		}
	}`

	err := db.C(TransactionCollection).Find(query).All(&transactions)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errNotExist.Error()})
		return
	}

	// we can also make use of aggregation pipeline to get sum, avergae, max, min from mongo query itself
	// stage_match := bson.M{"$match":bson.M{"timestamp": bson.M{"$gt": $(new Date(ISODate().getTime() - 1000 * 60 * 1))}}
	// stage_project:= bson.M{"$group": bson.M{"_id": null, "average": bson.M{"$avg": $amount}, "max": bson.M{"$max": $amount}, "min": bson.M{"$min": $amount}, "count": bson.M{"$sum": 1}}
	// _ := db.C(TransactionCollection).Pipe([]bson.M{stage_match, stage_project})

	var sum, average, max, min float64
	for i := 0; i < len(transactions); i++ {
		sum += transactions[i].Amount
		if transactions[i].Amount > max {
			max = transactions[i].Amount
		}
		if transactions[i].Amount < min {
			min = transactions[i].Amount
		}
	}
	average = sum / float64(len(transactions))
	c.JSON(http.StatusOK, gin.H{"sum": sum, "average": average, "max": max, "min": min, "count": len(transactions)})
}

// GetTransaction Endpoint
func GetTransaction(c *gin.Context) {
	var id bson.ObjectId = bson.ObjectIdHex(c.Param("id")) // Get Param
	transaction, err := transaction.TransactionInfo(id, TransactionCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errInvalidID.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "user": &transaction})
}

// CreateTransaction Endpoint
func CreateTransaction(c *gin.Context) {
	// Get DB from Mongo Config
	db := conn.GetMongoDB()
	transaction := transaction.Transaction{}
	err := c.Bind(&transaction)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errInvalidTransaction.Error()})
		return
	}

	transaction.ID = bson.NewObjectId()
	transaction.Timestamp = transaction.Timestamp // ISO 8601 format

	currentTimeStr := time.Now().UTC().Format(time.RFC3339) // convert current time to ISO 8601 format
	currentTime, _ := time.Parse("2006-01-02T15:04:05.000Z", currentTimeStr)

	if currentTime.Sub(transaction.Timestamp).Seconds() > 60 {
		c.JSON(204, gin.H{"message": errExpired.Error()})
		return
	}

	if transaction.Timestamp.Sub(currentTime).Seconds() > 0 {
		c.JSON(422, gin.H{"message": errInvalidTransaction.Error()})
		return
	}

	err = db.C(TransactionCollection).Insert(transaction)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errInsertionFailed.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "transaction": &transaction})
}

// DeleteTransactions Endpoint
func DeleteTransactions(c *gin.Context) {
	// Get DB from Mongo Config
	db := conn.GetMongoDB()
	// var id bson.ObjectId = bson.ObjectIdHex(c.Param("id")) // Get Param
	// err := db.C(TransactionCollection).RemoveAll(bson.M{"_id": &id})

	_, err := db.C(TransactionCollection).RemoveAll(bson.M{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errDeletionFailed.Error()})
		return
	}
	c.JSON(204, gin.H{"status": "success", "message": "Transactions deleted successfully"})
}
