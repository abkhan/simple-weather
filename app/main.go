package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"

	mysqlRepo "github.com/abkhan/simple-weather/internal/repository/mysql"

	"github.com/abkhan/simple-weather/article"
	"github.com/abkhan/simple-weather/internal/repository/weatherapi"
	"github.com/abkhan/simple-weather/internal/rest"
	"github.com/abkhan/simple-weather/internal/rest/middleware"
	"github.com/abkhan/simple-weather/internal/rest/wserver"
	"github.com/joho/godotenv"
)

const (
	defaultTimeout = 30
	defaultAddress = ":9090"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	//prepare database
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPass := os.Getenv("DATABASE_PASS")
	dbName := os.Getenv("DATABASE_NAME")
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)
	if err != nil {
		log.Fatal("failed to open connection to database", err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal("failed to ping database ", err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal("got error when closing the DB connection", err)
		}
	}()

	// prep for openWeather api access
	// Get apiKey from env
	owurl := os.Getenv("OPENWEATHER_URL") // example: https://api.openweathermap.org/data/2.5/weather
	if owurl == "" {
		owurl = "https://api.openweathermap.org/data/2.5/weather"
	}
	owApiKey := os.Getenv("OPENWEATHER_API_KEY")

	// first create weatherApi
	wapi := weatherapi.New(owurl, owApiKey)

	// then create and start server
	wserver := wserver.New(wapi)
	wserver.Start()

	// prepare echo
	e := echo.New()
	e.Use(middleware.CORS)
	timeoutStr := os.Getenv("CONTEXT_TIMEOUT")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		log.Println("failed to parse timeout, using default timeout")
		timeout = defaultTimeout
	}
	timeoutContext := time.Duration(timeout) * time.Second
	e.Use(middleware.SetRequestContextWithTimeout(timeoutContext))

	// Prepare Repository
	authorRepo := mysqlRepo.NewAuthorRepository(dbConn)
	articleRepo := mysqlRepo.NewArticleRepository(dbConn)

	// Build service Layer
	svc := article.NewService(articleRepo, authorRepo)
	rest.NewArticleHandler(e, svc)

	// Start Server
	addres := os.Getenv("SERVER_ADDRESS")
	if addres == "" {
		addres = defaultAddress
	}
	log.Fatal(e.Start(os.Getenv("SERVER_ADDRESS"))) //nolint
}
