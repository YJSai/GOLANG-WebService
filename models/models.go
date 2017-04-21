package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Posts struct {
	Id          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Userid      string        `json:"userid" bson:"userid"`
	Pic         string        `json:"pic" bson:"pic"`
	Contenttext string        `json:"contenttext" bson:"contenttext"`
	Createdtime time.Time     `json:"createdtime" bson:"createdtime"`
	Updatedtime time.Time     `json:"updatedtime" bson:"updatedtime"`
}

type Users struct {
	//Id        bson.ObjectId `json:"id" bson:"_id"`
	Username  string `json:"username" bson:"username"`
	Firstname string `json:"firstname" bson:"firstname"`
	Lastname  string `json:"lastname" bson:"lastname"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
}

type UserCredentials struct {
	Id       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Username string        `json:"username"`
	Password string        `json:"password"`
}

type Response struct {
	Data string `json:"data"`
}

type Token struct {
	Token  string `json:"token"`
	UserId string `json:"userid"`
}

type Foodstalls struct {
	Name     string `json:"name"`
	Foodtype string `json:"foodtype"`
	Avgcost  int    `json:"avgcost"`
	Geometry struct {
		Type        string
		Coordinates [2]float32
	} `json:"geometry"`
	Reviews []struct {
		Ref  string `json:"ref"`
		Type bson.ObjectId
	} `json:"reviews"`
}

type Reviews struct {
	Title      string
	text       string
	Foodstalls struct {
		Type bson.ObjectId
		ref  string
	} `json:"foodstalls"`
}

func (self *Foodstalls) InitDefaults() {
	self.Geometry.Type = "Point"
}
