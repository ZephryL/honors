package common

import (
	"errors"
	"flag"
	"fmt"
	"time"
	"github.com/gorilla/securecookie"
)

//--------------------------------------------------------
// System struct and methods
//--------------------------------------------------------
type System struct {
	User     string
	Password string
	Schema   string
	Port     int
	Cookie   *securecookie.SecureCookie
}

func (this *System) GetFlags() error {
	flag.StringVar(&this.User, "user", "", "The database User")
	flag.StringVar(&this.Password, "password", "password", "The database Password")
	flag.StringVar(&this.Schema, "schema", "zephry", "The database Schema")
	flag.IntVar(&this.Port, "port", 8080, "The server port")
	flag.Parse()
	// Verify there are flags
	if this.User == "" ||
		this.Password == "" ||
		this.Schema == "" ||
		this.Port < 1 {
		return errors.New("All flags are required. Some flags have either been omitted or have null values.")
	}
	return nil
}

func (this *System) SetCookie() {
	var hashKey = securecookie.GenerateRandomKey(32)
	var blockKey = securecookie.GenerateRandomKey(16)
	this.Cookie = securecookie.New(hashKey, blockKey)
}

func (this *System) SqlString() string {
	return fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v?parseTime=true", this.User, this.Password, this.Schema)
}

//--------------------------------------------------------
// Standard Response for handlers that have no return data
//--------------------------------------------------------
type StandardResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// or a passthru value from client
type Arbitrary struct {
	Passthru string `json:"passthru"`
}

// KeyValue
type KeyValue struct {
	Key       string  `json:"key"`
	Value     float32 `json:"value"`
	Precision int     `json:"precision"`
}

type KeyValueList struct {
	Count int        `json:"count"`
	Loval float32    `json:"loval"`
	Hival float32    `json:"hival"`
	List  []KeyValue `json:"list"`
}

// Number Cards
type NumberCards struct {
	List []KeyValue `json:"list"`
}

// Welcome
type Welcome struct {
	Date      time.Time  `json:"date"`
	Status    int        `json:"status"`
	Greeting  string     `json:"greeting"`
	Endpoints []Endpoint `json:"endpoints"`
}

// Endpoint
type Endpoint struct {
	Url         string `json:"url"`
	Description string `json:"desc"`
	ContentType string `json:"ctype"`
}

// Routes
type Route struct {
	Path string `json:"path"`
	Queries string `json:"queries"`
	Methods string `json:"methods"`
}
type RouteList struct {
	List []Route `json:"list"`
}

//--------------------------------------------------------
// Email Host and snippets
//--------------------------------------------------------
type EmailHost struct {
	Ehs_Key int `json:"ehskey"`
	Ehs_Default bool `json:"default"`
	Ehs_Name string `json:"name"`
	Ehs_Port int `json:"port"`
	Ehs_User  string `json:"user"`
	Ehs_Password string `json:"password"`
}

type EmailHostCollection struct {
	List []EmailHost `json:"list"`
}

//--------------------------------------------------------
// User and snippets
//--------------------------------------------------------

type Register struct {
	Reg_Key int `json:"regkey"`
	Reg_Firstnames string `json:"firstnames"`
	Reg_Surname string `json:"surname"`
	Reg_Email string `json:"email"`
	Reg_Password string `json:"password"`
	Reg_Token string `json:"token"`
	Reg_Date time.Time `json:"date"`
	Reg_Passthru string `json:"passthru"`
}

type Verify struct {
	Reg_Key int `json:"regkey"` //,string,omitempty"`
	Reg_Token string `json:"token"`
	Reg_Date string `json:"date"`
}

type Login struct {
	Email string `json:"email"`
	Token string `json:"token"`
	Passthru string `json:"passthru"`
}

type Account struct {
	Email string `json:"email"`
	Fullname string `json:"fullnames"`
}

type Token struct {
	Key int `json:"key"`
	Token string `json:"token"`
}

type User struct {
	Usr_Key int `json:"key"`
	Usr_Email string `json:"email"`
	Usr_Password string `json:"password"`
	Usr_Token  string `json:"token"`
	Usr_Firstnames string `json:"firstnames"`
	Usr_Surname string `json:"surname"`
	Usr_Identity string `json:"identity"`
}

//--------------------------------------------------------
// Application and snippets
//--------------------------------------------------------
type Application struct {
	App_Key int `json:"key"`
	App_Code string `json:"code"`
	App_Name string `json:"name"`
	App_Desc  string `json:"desc"`
	App_Lock bool `json:"lock"`
}

type ApplicationCollection struct {
	List []Application `json:"list"`
}

//--------------------------------------------------------
// Role and snippets
//--------------------------------------------------------
type Role struct {
	Key int `json:"key"`
	Code string `json:"code"`
	Name string `json:"name"`
	Desc  string `json:"desc"`
	Auth bool `json:"auth"`
	Lock bool `json:"lock"`
}

type RoleCollection struct {
	Account Account `json:"account"`
	List []Role `json:"list"`
}

