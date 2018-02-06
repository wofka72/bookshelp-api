package main

import (
    "fmt"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"

    "github.com/jinzhu/gorm"
    _ "github.com/lib/pq"
)


var db *gorm.DB
var err error


type Book struct {
    ID			int    `json:"id"`
    Title		string `json:"title"`
    Author		string `json:"author"`
    Publisher	string `json:"publisher"`
    PagesCount  int    `json:"pagescount"`
    Description string `json:"description"`
}


func main() {
    const (
        user     = "postgres"
        password = "password"
        dbname   = "books_shelf"
    )

    psqlInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
        user, password, dbname)
    db, err = gorm.Open("postgres", psqlInfo)
    if err != nil {
        fmt.Println(err)
    }
    defer db.Close()
    db.AutoMigrate(&Book{})

    router := gin.Default()

    router.Use( cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"*"},
        AllowHeaders:     []string{"*"},
        ExposeHeaders:    []string{"*"},
        AllowCredentials: true,
        AllowOriginFunc:  func(origin string) bool {
            return true
            // return origin == "https://github.com"
        },
        MaxAge: 12 * time.Hour,
    }))

    // router.Use(cors.Default())

    router.GET("/books/", GetBooks)
    router.GET("/last/", GetLast)
    router.GET("/books/:id", GetBook)
    router.POST("/books", AddBook)
    router.PUT("/books/:id", UpdateBook)
    router.DELETE("/books/:id", DeleteBook)
    router.DELETE("/del_last/", DeleteLast)

    router.Run(":9191")
}


func GetBooks(c *gin.Context) {
    var books []Book
    if err := db.Find(&books).Error; err != nil {
        c.AbortWithStatus(404)
        fmt.Println(err)
    } else {
        c.JSON(200, books)
    }
}

func GetBook(c *gin.Context) {
    id := c.Params.ByName("id")
    var book Book
    if err := db.Where("id = ?", id).First(&book).Error; err != nil {
        c.AbortWithStatus(404)
        fmt.Println(err)
    } else {
        c.JSON(200, book)
    }
}

func GetLast(c *gin.Context) {
    var lastBook Book
    if err := db.Last(&lastBook).Error; err != nil {
        c.AbortWithStatus(404)
        fmt.Println(err)
    } else {
        c.JSON(200, lastBook)
    }
}

func AddBook(c *gin.Context) {
    var book Book
    c.BindJSON(&book)

    db.Create(&book)
    c.JSON(200, book)
}

func UpdateBook(c *gin.Context) {
    var book Book
    id := c.Params.ByName("id")

    if err := db.Where("id = ?", id).First(&book).Error; err != nil {
        c.AbortWithStatus(404)
        fmt.Println(err)
    }
    c.BindJSON(&book)

    db.Save(&book)
    c.JSON(200, book)
}

func DeleteBook(c *gin.Context) {
    id := c.Params.ByName("id")
    var book Book
    db.Where("id = ?", id).Delete(&book)
    c.JSON(200, gin.H{"id #" + id: "deleted"})
}

func DeleteLast(c *gin.Context) {
    var lastBook Book
    if err := db.Last(&lastBook).Error; err != nil {
        c.AbortWithStatus(404)
        fmt.Println(err)
    } else {
        if lastBook.ID == 0 {
            c.JSON(200, "the library is empty")
        } else {
            lastID := strconv.Itoa(int(lastBook.ID))
            db.Where("id = ?", lastID).Delete(&lastBook)
            c.JSON(200, gin.H{"id #" + lastID: "deleted"})
        }
    }
}