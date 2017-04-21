package controllers

import (
	"crypto/rsa"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"goapi/models"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// added session to our userController
type UserController struct {
	session *mgo.Session
}

// added session to our userController
func NewUserController(s *mgo.Session) *UserController {
	return &UserController{s}
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

const (
	// For simplicity these files are in the same folder as the app binary.
	// You shouldn't do this in production.
	privKeyPath = "controllers/demo.rsa"
	pubKeyPath  = "controllers/demo.rsa.pub"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func (uc UserController) InitKeys() {
	signBytes, err := ioutil.ReadFile(privKeyPath)
	fatal(err)

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	fatal(err)

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)
}

func (uc UserController) uploadFile(w http.ResponseWriter, r *http.Request) {

}

func (uc UserController) Hello(w http.ResponseWriter, r *http.Request, _ http.HandlerFunc) {
	w.Write([]byte("hello"))
	//fmt.Fprint(w, "hello")
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	//id := p.ByName("id")
	username := p.ByName("username")
	email := p.ByName("email")

	// if !bson.IsObjectIdHex(id) {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }
	// oid := bson.ObjectIdHex(id)

	u := models.Users{}

	if err := uc.session.DB("user").C("users").Find(bson.M{"username": username, "email": email}).One(&u); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("no this user"))
		return
	}

	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := models.Users{}

	json.NewDecoder(r.Body).Decode(&u)

	// u.Id = bson.NewObjectId()

	uc.session.DB("user").C("users").Insert(u)

	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// TODO: only write code to delete user
	id := p.ByName("id")
	//u := models.Users{}
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	oid := bson.ObjectIdHex(id)
	if err := uc.session.DB("user").C("users").RemoveId(oid); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//Remove(bson.M{"isbn": isbn})
	w.WriteHeader(http.StatusOK) // 200
	fmt.Fprint(w, "Write code to delete user\n")
}

func (uc UserController) Getallposts(w http.ResponseWriter, req *http.Request, _ http.HandlerFunc) {

	type Userid struct {
		Id string `json:"id" bson:"id"`
	}
	fmt.Println("now in getallposts")
	//var user models.UserCredentials

	type UserPosts struct {
		Posts []models.Posts `json:"posts"`
	}

	var posts []models.Posts
	var uid Userid
	err := json.NewDecoder(req.Body).Decode(&uid)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	fmt.Println("parse " + uid.Id + " json successfully")
	if err := uc.session.DB("user").C("posts").Find(bson.M{"userid": uid.Id}).All(&posts); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	userPosts := UserPosts{posts}
	uj, err := json.Marshal(userPosts)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc UserController) Uploadpost(w http.ResponseWriter, req *http.Request, _ http.HandlerFunc) {
	//tpl.ExecuteTemplate(w, "index.html", nil)
	post := models.Posts{}

	//c := getCookie(w, req)

	if req.Method == http.MethodPost {
		//fmt.Println("post")
		fmt.Println(req.FormValue("userid"))
		fmt.Println(req.FormValue("contenttext"))

		post.Userid = req.FormValue("userid")
		post.Contenttext = req.FormValue("contenttext")
		mf, fh, err := req.FormFile("avatarFile")
		if err != nil {
			fmt.Println(err)
		}
		defer mf.Close()
		// create sha for file name
		ext := strings.Split(fh.Filename, ".")[1]
		h := sha1.New()
		io.Copy(h, mf)
		fname := fmt.Sprintf("%x", h.Sum(nil)) + "." + ext
		// create new file
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		path := filepath.Join(wd, "public", "pics", fname)
		// d1 := []byte("hello\ngo\n")
		// error := ioutil.WriteFile("./public/pics", d1, 0644)
		// if error != nil {
		// 	fmt.Println(error)
		// }
		//fmt.Println(path)
		nf, err := os.Create(path)
		if err != nil {
			fmt.Println(err)
		}
		defer nf.Close()
		// copy
		mf.Seek(0, 0)
		io.Copy(nf, mf)
		post.Pic = "http://localhost:8000/images/" + fname
		post.Createdtime = time.Now()
		post.Updatedtime = time.Now()
		uc.session.DB("user").C("posts").Insert(post)
		//c = appendValue(w, c, fname)
		p, err := json.Marshal(&post)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(p))
		//uc.session.DB("user").C("posts").Insert(p)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // 201
		fmt.Fprintf(w, "%s\n", p)

	}
}

func (uc UserController) LoginHandler(w http.ResponseWriter, r *http.Request, _ http.HandlerFunc) {

	if r.Method == http.MethodPost {
		var user models.UserCredentials
		// var foodStalls models.Foodstalls
		// foodStalls.InitDefaults
		//_, params, _ :=  router.Lookup("GET", req.URL.Path)

		err := json.NewDecoder(r.Body).Decode(&user)
		//json, err := json.Marshal(r.Body)
		// fmt.Println([]byte(json))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "can't parse json send from client")
			return
		}

		//userinfo = uc.session.DB("user").C("users").Find(bson.M{"username": user.Username, "password": user.Password}).One(&user)
		if err := uc.session.DB("user").C("users").Find(bson.M{"username": user.Username, "password": user.Password}).One(&user); err != nil {
			w.WriteHeader(http.StatusForbidden)
			fmt.Println("Error logging in")
			fmt.Fprint(w, "Invalid credentials")
			return
		}

		// if strings.ToLower(user.Username) != "someone" {
		// 	if user.Password != "p@ssword" {
		// 		w.WriteHeader(http.StatusForbidden)
		// 		fmt.Println("Error logging in")
		// 		fmt.Fprint(w, "Invalid credentials")
		// 		return
		// 	}
		// }

		token := jwt.New(jwt.SigningMethodRS256)
		claims := make(jwt.MapClaims)
		claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
		claims["iat"] = time.Now().Unix()
		token.Claims = claims

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Error extracting the key")
			fatal(err)
		}

		tokenString, err := token.SignedString(signKey)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Error while signing the token")
			fatal(err)
		}

		response := models.Token{tokenString, string(user.Id.Hex())}
		JsonResponse(response, w)
	}

}

func JsonResponse(response interface{}, w http.ResponseWriter) {

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (uc UserController) ProtectedHandler(w http.ResponseWriter, r *http.Request) {

	response := models.Response{"Gained access to protected resource"}
	JsonResponse(response, w)

}

func (uc UserController) ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		})

	if err == nil {
		if token.Valid {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized access to this resource")
	}

}
