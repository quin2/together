package main

import (
	"github.com/gin-gonic/gin"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/lib/pq"

	"encoding/json"
	"net/http"
	"os"
	"fmt"
)

var (
	db = openDB()
)

//how are entry IDs handled?
type User struct {
	gorm.Model
	Name      string `json:"name"`
	EmailAddr string `json:"email"`
	RelationshipID uint `json:"relID"`
}

type Relationship struct {
	gorm.Model
	Users      []User      `json:"users"`
	Scrapbooks []Scrapbook `json:"scrapbook"`
	Countdowns []Countdown `json:"countdowns"`
	Albums     []Album     `json:"albums"`
}

type Scrapbook struct {
	gorm.Model
	Title  string   `json:"title"`
	Photos pq.StringArray `gorm:"type:varchar(64)[]" json:"photos"`
	Song   string   `json:"song"`
	RelationshipID uint `json:"relID"`
}

type Countdown struct {
	gorm.Model
	TimeTo   string    `json:"timeTo"`
	Messages []Message `json:"messages"`
	RelationshipID uint `json:"relID"`
}

type Album struct {
	gorm.Model
	Title  string   `json:"title"`
	Photos pq.StringArray `gorm:"type:varchar(64)[]" json:"photos"`
	RelationshipID uint `json:"relID"`
}

type Message struct {
	gorm.Model
	Sender int    `json:"sender"`
	Text   string `json:"text"`
	CountdownID uint `json:"countdownID"`
}

func main() {
	//begin router logic
	router := gin.Default()

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			users := v1.Group("/users")
			{
				//users.GET("/me", getMe)         //get self from middleware(todo)
				users.GET("/user/:id", getUser) //get user by ID
			}

			/*
			rships := v1.Group("/rships")
			{
				rships.POST("/join", joinRship)   //join relationship
				rships.POST("/leave", leaveRship) //leave relationship
			}
			*/
		}
	}

	router.Run()
}

func openDB() *gorm.DB { //find gorm data type
	//host := "localhost" //fix this haha
	pass := os.Getenv("PGPASS")
	user := os.Getenv("PGUSER")
	dbName := "rships" //might have to change later for zeet
	port := os.Getenv("PGPORT")

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable", user, pass, dbName, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil{
		panic("failed to connect to DB")
	}

	// Migrate the schema
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Relationship{})
	db.AutoMigrate(&Scrapbook{}) 
	db.AutoMigrate(&Countdown{}) //problem
	db.AutoMigrate(&Album{}) //problem
	db.AutoMigrate(&Message{}) //problem

	return db
}

func getUser(c *gin.Context) {
	var id = c.Param("id")

	var user User
	db.First(&user, id)

	res, bERR := json.Marshal(user)
	if bERR != nil {
		fmt.Println("json marshal error")
	}

	c.JSON(http.StatusOK, res)
}


