package main

import (
	"GoAPi/controllers"

	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/urfave/negroni"
	"gopkg.in/mgo.v2"
)

var tpl *template.Template

/*

 */

/*

 */
func init() {
	tpl = template.Must(template.ParseGlob("views/*"))
}

func main() {
	//router := httprouter.New()
	uc := controllers.NewUserController(getSession())
	uc.InitKeys()
	// http.Handle("/files/", http.StripPrefix("/files",
	// 	http.FileServer(http.Dir("./static"))))
	//http.HandleFunc("/", index)
	//http.HandleFunc("/", helloo)
	//http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./public/pics/"))))
	// http.Handle("/", http.FileServer(http.Dir("./static/")))
	// http.Handle("/images", negroni.New(
	// 	negroni.HandlerFunc(uc.ValidateTokenMiddleware),
	// 	negroni.Wrap((http.FileServer(http.Dir("./public/pics/")))),
	// ))
	h := http.NewServeMux()
	// router.Handler("GET", "/images/", negroni.New(
	// 	//negroni.HandlerFunc(uc.ValidateTokenMiddleware),
	// 	negroni.Wrap(http.StripPrefix("/images/", http.FileServer(http.Dir("./public/pics/")))),
	// ))

	h.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, you hit foo!")
	})

	/*** for httprouter

	  ****
	//router.ServeFiles("/public/pics/*filepath", http.Dir("./public/pics/"))
	//
	router.Handler("POST", "/login", negroni.New(negroni.HandlerFunc(uc.LoginHandler)))
	//
	// router.Handler("GET", "/images", negroni.New(
	// 	negroni.HandlerFunc(uc.ValidateTokenMiddleware),
	// 	negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		http.FileServer(http.Dir("./public/pics/")).ServeHTTP(w, r)
	// 	})),
	// ))
	//
	router.Handler("GET", "/hello", negroni.New(
		negroni.HandlerFunc(uc.ValidateTokenMiddleware),
		negroni.HandlerFunc(uc.Hello),
	))

	//
	router.Handler("POST", "/uploadpost", negroni.New(
		negroni.HandlerFunc(uc.ValidateTokenMiddleware),
		negroni.HandlerFunc(uc.Uploadpost),
	))

	//
	router.Handler("POST", "/allposts", negroni.New(
		negroni.HandlerFunc(uc.ValidateTokenMiddleware),
		negroni.HandlerFunc(uc.Getallposts),
	))
	//
	router.POST("/user", uc.CreateUser)
	router.DELETE("/user/:username/:email", uc.DeleteUser)

	***
	  //////
	  ****/
	/*
	 *****/
	http.Handle("/images/", negroni.New(
		negroni.HandlerFunc(uc.ValidateTokenMiddleware),
		negroni.Wrap(http.StripPrefix("/images/", http.FileServer(http.Dir("./public/pics/")))),
	))

	http.Handle("/login", negroni.New(negroni.HandlerFunc(uc.LoginHandler)))
	http.Handle("/hello", negroni.New(
		//negroni.HandlerFunc(uc.ValidateTokenMiddleware),
		negroni.HandlerFunc(uc.Hello),
	))

	http.Handle("/uploadpost", negroni.New(
		negroni.HandlerFunc(uc.ValidateTokenMiddleware),
		negroni.HandlerFunc(uc.Uploadpost),
	))

	http.Handle("/allposts", negroni.New(
		negroni.HandlerFunc(uc.ValidateTokenMiddleware),
		negroni.HandlerFunc(uc.Getallposts),
	))

	/*
	*****/

	//http.Handle("/", http.HandlerFunc(helloo))
	//http.Handle("/", http.HandlerFunc(helloo))
	//n := negroni.New(negroni.HandlerFunc(uc.ValidateTokenMiddleware))
	//http.Handle("/login", negroni.New(negroni.HandlerFunc(uc.LoginHandler)))
	//router.Handler("POST", "/login", negroni.New(negroni.HandlerFunc(uc.LoginHandler)))
	//http.HandleFunc("/", helloo)
	//router.POST("/login", uc.LoginHandler)
	//router.ServeFiles("/views/*filepath", http.Dir("views"))
	// router.ServeFiles("/public/pics/*filepath", http.Dir("./public/pics/"))

	// router.GET("/images/", (func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// 	router.ServeFiles("/public/pics/*filepath", http.Dir("./public/pics/"))
	// 	http.StripPrefix("/images/", http.FileServer(http.Dir("./public/pics/")))
	// }))

	// router.Handler("GET", "/images", negroni.New(
	// 	negroni.HandlerFunc(uc.ValidateTokenMiddleware),
	// 	negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		http.FileServer(http.Dir("./public/pics/")).ServeHTTP(w, r)
	// 	})),
	// ))

	//router.Handler("GET", "/images", http.FileServer(http.Dir("/public/pics/")))
	// router.Handler("GET", "/images", negroni.New(
	// 	negroni.HandlerFunc(uc.ValidateTokenMiddleware),
	// 	negroni.Wrap((http.FileServer(http.Dir("./public/pics/")))),
	// ))
	// negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	http.FileServer(http.Dir("./public/pics/")).ServeHTTP(w, r)
	// }))

	// router.GET("/resource", negroni.New(
	// 	negroni.HandlerFunc(uc.ValidateTokenMiddleware),
	// 	negroni.New(negroni.HandlerFunc(uc.LoginHandler)),
	// ))
	// http.Handle("/hello", negroni.New(
	// 	negroni.HandlerFunc(uc.ValidateTokenMiddleware),
	// 	negroni.HandlerFunc(uc.Hello),
	// ))
	// router.Handler("GET", "/hello", negroni.New(
	// 	negroni.HandlerFunc(uc.ValidateTokenMiddleware),
	// 	negroni.HandlerFunc(uc.Hello),
	// ))
	// router.GET("/", uploadpost)
	// router.POST("/uploadpost", uploadpost)

	// router.Handler("POST", "/uploadpost", negroni.New(
	// 	negroni.HandlerFunc(uc.ValidateTokenMiddleware),
	// 	negroni.HandlerFunc(uc.Uploadpost),
	// ))

	// router.Handler("POST", "/allposts", negroni.New(
	// 	negroni.HandlerFunc(uc.ValidateTokenMiddleware),
	// 	negroni.HandlerFunc(uc.Getallposts),
	// ))
	//router.GET("/hello", uc.Hello)
	// router.GET("/user/:username/:email", uc.GetUser)

	//http.ListenAndServe("localhost:8080", nil)
	//n.UseHandler(router)
	log.Println("Now listening...")
	http.ListenAndServe("localhost:8000", nil)
	//http.ListenAndServe("localhost:8000", router)

}

func helloo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func getCookie(w http.ResponseWriter, req *http.Request) *http.Cookie {
	c, err := req.Cookie("session")
	if err != nil {
		sID := uuid.NewV4()
		c = &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
	}
	return c
}

func appendValue(w http.ResponseWriter, c *http.Cookie, fname string) *http.Cookie {
	s := c.Value
	if !strings.Contains(s, fname) {
		s += "|" + fname
	}
	c.Value = s
	http.SetCookie(w, c)
	return c
}

func getSession() *mgo.Session {
	// Connect to our local mongo
	//s, err := mgo.Dial("mongodb://mongo:mongo@ds119750.mlab.com:19750/user")

	// Check if connection error, is mongo running?
	// if err != nil {
	// 	fmt.Printf("failed")
	// 	panic(err)
	// }
	//return s

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{"ds119750.mlab.com:19750"},
		Username: "mongo",
		Password: "123",
		Database: "user",
		//ReplicaSetName: "ReplicaSetName",
	})

	if err != nil {
		fmt.Printf("failed")
		panic(err)
	}
	return session
}
