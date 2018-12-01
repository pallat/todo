package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

var url = "mongodb://%s:27017"

const database = "tech_inno"
const collection = "test"

func main() {
	var err error

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "."))
	viper.SetDefault("mongodb", "localhost")
	port := viper.GetString("port")

	url = fmt.Sprintf(url, viper.GetString("mongodb"))
	fmt.Println(url)

	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	h := todoHandler{session: session}
	e.POST("/todos", h.NewTodoHandler)
	e.GET("/todos", ListTodoHandler, MiddlewareSession(session))
	e.PUT("/todos/:id", func(c echo.Context) error {
		err := DoneTodo(session, c.Param("id"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, nil)
	})
	e.Logger.Fatal(e.Start(":" + port))
}

type todo struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Topic string        `json:"topic"`
	Done  bool          `json:"done"`
}

type todoHandler struct {
	session *mgo.Session
}

func (h *todoHandler) NewTodoHandler(c echo.Context) error {
	var t todo
	err := c.Bind(&t)
	if err != nil {
		return err
	}
	NewTodo(h.session, t.Topic)
	return c.JSON(http.StatusOK, nil)
}

func NewTodo(session *mgo.Session, topic string) {
	s := session.Copy()
	defer s.Close()
	c := s.DB(database).C(collection)
	c.Insert(todo{ID: bson.NewObjectId(), Topic: topic})
}

func ListTodoHandler(c echo.Context) error {
	session := c.Get("mgoSession").(*mgo.Session)
	list := ListTodo(session)
	return c.JSON(http.StatusOK, list)
}

func ListTodo(session *mgo.Session) []todo {
	var all []todo
	s := session.Copy()
	defer s.Close()
	c := s.DB(database).C(collection)
	err := c.Find(nil).All(&all)
	if err != nil {
		return nil
	}
	return all
}

func DoneTodo(session *mgo.Session, id string) error {
	var otodo todo
	s := session.Copy()
	defer s.Close()
	c := s.DB(database).C(collection)
	err := c.FindId(bson.ObjectIdHex(id)).One(&otodo)
	if err != nil {
		return err
	}

	otodo.Done = true
	return c.UpdateId(bson.ObjectIdHex(id), &otodo)
}

func MiddlewareSession(session *mgo.Session) echo.MiddlewareFunc {
	return func(handler echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("mgoSession", session)
			return handler(c)
		}
	}
}
