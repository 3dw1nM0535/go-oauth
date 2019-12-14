package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/3dw1nM0535/go-auth/utils"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// User : retrieved and authenticated
type User struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Gender        string `json:"gender"`
}

var clientid, clientsecret string
var conf *oauth2.Config
var state string
var store = sessions.NewCookieStore([]byte("secret"))

// indexHandler : handle index page
func indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}

// getLogin : get state from authentication url
func getLoginURL(state string) string {
	return conf.AuthCodeURL(state)
}

// randToken : generate random token
func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// LoginHandler : store token in session
func loginHandler(c *gin.Context) {
	state = randToken()
	session := sessions.Default(c)
	session.Set("state", state)
	session.Save()
	c.Writer.Write([]byte("<html><title>Golang Google</title> <body><a href='" + getLoginURL(state) + "'><button>Login with Google</button></a></body></html>"))
}

func init() {
	err := godotenv.Load()
	clientid = utils.MustGet("ClientID")
	clientsecret = utils.MustGet("ClientSecret")
	file, err := ioutil.ReadFile("./cred.json")
	log.Println(file)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}

	conf = &oauth2.Config{
		ClientID:     clientid,
		ClientSecret: clientsecret,
		RedirectURL:  "http://localhost:9090/auth",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
}

func authHandler(c *gin.Context) {
	// Check state validity
	session := sessions.Default(c)
	retrievedState := session.Get("state")
	if retrievedState != c.Request.URL.Query().Get("state") {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state: %s", retrievedState))
		return
	}

	// Handle the exchange code to initiate transport
	token, err := conf.Exchange(oauth2.NoContext, c.Query("code"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		log.Println("Error during credential exchange: ", err.Error())
		return
	}
	// Construct the client
	client := conf.Client(oauth2.NoContext, token)
	res, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		log.Println("Error communicating with google API: ", err.Error())
		return
	}

	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	log.Println("Response Body: ", string(data))
}

func main() {
	app := gin.New()
	app.Use(sessions.Sessions("goquestsession", store))
	app.Static("/css", "./static/css")
	app.Static("/img", "./static/img")
	app.LoadHTMLGlob("templates/*")

	app.GET("/", indexHandler)
	app.GET("/login", loginHandler)
	app.GET("/auth", authHandler)
	app.Run("127.0.0.1:9090")
}
