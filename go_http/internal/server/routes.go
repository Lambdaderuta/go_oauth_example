package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.healthHandler)
	r.GET("/auth/:provider/callback", s.GetAuthCallback)
	r.GET("/auth/:provider/logout", s.OauthLogout)
	r.GET("/auth/:provider", s.OauthProvider)

	return r
}

func (s *Server) GetAuthCallback(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Add("provider", c.Param("provider"))
	c.Request.URL.RawQuery = q.Encode()

	fmt.Println("here cb")

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		fmt.Fprintln(c.Writer, err)
		return
	}

	u, _ := json.Marshal(user)

	fmt.Println(u, 'u')

	if err := gothic.StoreInSession("user", string(u), c.Request, c.Writer); err != nil {
		fmt.Println(err)
	}

	http.Redirect(c.Writer, c.Request, "http://localhost:5173", http.StatusFound)
}

func (s *Server) OauthLogout(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Add("provider", c.Param("provider"))
	c.Request.URL.RawQuery = q.Encode()

	gothic.Logout(c.Writer, c.Request)
	c.Writer.Header().Set("Location", "/")
	c.Writer.WriteHeader(http.StatusTemporaryRedirect)
}

func (s *Server) OauthProvider(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Add("provider", c.Param("provider"))
	c.Request.URL.RawQuery = q.Encode()

	if user, err := gothic.CompleteUserAuth(c.Writer, c.Request); err == nil {
		fmt.Println(user)
	} else {
		gothic.BeginAuthHandler(c.Writer, c.Request)
	}
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	// u, _ := gothic.GetFromSession("user", c.Request)
	resp["message"] = "asdsd"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
