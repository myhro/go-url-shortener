package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// URL is the main resource
type URL struct {
	ID   int
	Hash string
	Full string `json:"url"`
}

func details(c *gin.Context) {
	var u URL

	u.Hash = c.Param("hash")
	row := db.QueryRow("SELECT full FROM urls WHERE hash = ?", u.Hash)
	err := row.Scan(&u.Full)
	if err == sql.ErrNoRows {
		c.String(http.StatusNotFound, "Not Found")
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"hash": u.Hash,
		"url":  u.Full,
	})
}

func index(c *gin.Context) {
	c.String(http.StatusOK, "Go URL Shortener")
}

func migrate() {
	f, err := os.Open("schema.sql")
	if err != nil {
		log.Print(err)
		return
	}
	defer f.Close()

	schema, err := ioutil.ReadAll(f)
	if err != nil {
		log.Print(err)
		return
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		log.Print(err)
		return
	}
}

func newURL(c *gin.Context) {
	var u URL

	err := c.ShouldBindJSON(&u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	} else if u.Full == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Empty URL",
		})
		return
	}

	var last int
	row := db.QueryRow("SELECT id FROM urls ORDER BY id DESC LIMIT 1")
	err = row.Scan(&last)
	if err == sql.ErrNoRows {
		// Empty database, first entry.
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	u.ID = last + 1
	u.Hash = strconv.FormatInt(int64(u.ID), 36)
	stmt, err := db.Prepare("INSERT INTO urls(hash, full) VALUES(?, ?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	_, err = stmt.Exec(u.Hash, u.Full)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"hash": u.Hash,
		"url":  u.Full,
	})
}

func setupDB(dbFile string) {
	var err error
	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	migrate()
}

func setupRouter(r *gin.Engine) {
	r.GET("/", index)
	r.GET("/:hash", shortURL)
	r.GET("/:hash/details", details)
	r.POST("/", newURL)
}

func shortURL(c *gin.Context) {
	var u URL

	u.Hash = c.Param("hash")
	row := db.QueryRow("SELECT full FROM urls WHERE hash = ?", u.Hash)
	err := row.Scan(&u.Full)
	if err == sql.ErrNoRows {
		c.String(http.StatusNotFound, "Not Found")
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Redirect(http.StatusMovedPermanently, u.Full)
}

func main() {
	dbFile := os.Getenv("DB_FILE")
	if dbFile == "" {
		dbFile = "urls.db"
	}
	setupDB(dbFile)

	r := gin.Default()
	setupRouter(r)
	r.Run()
}
