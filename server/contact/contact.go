package contact

import (
	"database/sql"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// contactMessage represents the structure of an incoming contact message.
type contactMessage struct {
	Id		int `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

// GetMessages gets all the contact messages.
func GetMessages() func(context *gin.Context) {
	return func(context *gin.Context) {
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Methods", "GET")
		
		db, _ := sql.Open("sqlite3", "./storage/bgguna.db")
		
		var messages = []contactMessage{}
		log.Infof("Fetching contact messages...")
		rows, err := db.Query("SELECT * FROM contact")
		if err != nil {
			log.Errorf("Error preparing to fetch all contact messages.", err)
			context.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
		}

		for rows.Next() {
			var msg contactMessage
			rows.Scan(&msg.Id, &msg.Name, &msg.Email, &msg.Phone, &msg.Message)
			messages = append(messages, msg)
		}

		log.Infof("Fetched all contact messages.")
		context.JSON(http.StatusOK, messages)
	}
}

// HandleNewMsg handles incoming messages from the Contact tab.
func HandleNewMsg() func(context *gin.Context) {
	return func(context *gin.Context) {
		context.Header("Access-Control-Allow-Origin", "*")
    	context.Header("Access-Control-Allow-Methods", "POST")

		db, _ := sql.Open("sqlite3", "./storage/bgguna.db")
		message := contactMessage{}
		rawContextData, err := context.GetRawData()
		if err != nil {
			log.Errorf("Failed to process request.", err)
			context.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
		}

		err = json.Unmarshal(rawContextData, &message)
		if err != nil {
			log.Errorf("Error processing request.", err)
			context.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
		}

		log.Printf("Received contact message: %+v", message)
		statement, err := db.Prepare("INSERT INTO contact (name, email, phone, message) VALUES (?, ?, ?, ?)")
		if err != nil {
			log.Errorf("Error preparing to store contact message.", err)
			context.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
		}

		_, err = statement.Exec(message.Name, message.Email, message.Phone, message.Message)
		if err != nil {
			log.Errorf("Error storing contact message.", err)
			context.JSON(http.StatusBadRequest, gin.H{"status": "fail"})
		}

		log.Infof("Contact message stored.")
		context.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}
