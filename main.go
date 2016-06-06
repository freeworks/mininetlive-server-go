package main

import (
	"database/sql"
	"log"
	"net/http"
	//	"strconv"
	"time"

	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	//	"gopkg.in/gorp.v1"
	"github.com/coopernurse/gorp"
)

type Resp struct {
	Ret  int64      `form:"ret" json:"ret"`
	Msg  string      `form:"msg" json:"msg"`
	Data interface{} `form:"data" json:"data"`
}

type User struct {
	Id     int    `form:"id" ` //db:"id,primarykey, autoincrement"
	Name   string `form:"name" binding:"required"  db:"name"`
	Avatar string `form:"avatar"  binding:"required"  db:"avatar"`
	Gender int    `form:"gender" binding:"required"  db:"gender"`
}

type OAuth struct {
	Id          int       `form:"id"` //  `form:"id"  db:"id,primarykey, autoincrement"`
	UserId      int       `form:"userId"  db:"user_id"`
	Plat        string    `form:"plat" binding:"required" db:"plat"`
	OpenId      string    `form:"openid" binding:"required" db:"openid"`
	AccessToken string    `form:"access_token" binding:"required" db:"access_token"`
	ExpiresIn   int       `form:"expires_in" binding:"required" db:"-"` //- 忽略的意思
	Expires     time.Time `db:"expires"`
}

type LocalAuth struct {
	Id          int       `form:"id"` //  `form:"id"  db:"id,primarykey, autoincrement"`
	UserId      int       `form:"userId"  db:"user_id"`
	Phone       string    `form:"phone" binding:"required"  db:"phone"`
	Password    string    `form:"password" binding:"required" db:"password"`
	Token 		string    `db:"token"`
	Expires     time.Time `db:"expires"`
}

type OAuthUser struct {
	User  User
	OAuth OAuth
}

type LocalAuthUser struct {
	User  User
	LocalAuth LocalAuth
}

type Activity struct {
	Id   int
	Title string
}

func GetActivityList(r render.Render, dbmap *gorp.DbMap) {
	var activities []Activity
	_, err := dbmap.Select(&activities, "select * from t_activity")
	checkErr(err, "Select failed")
	newmap := map[string]interface{}{"metatitle": "this is my custom title", "activities": activities}
	r.HTML(200, "posts", newmap)
}

func GetActivity(args martini.Params, r render.Render, dbmap *gorp.DbMap) {
	var activities []Activity
	_, err := dbmap.Select(&activities, "select * from t_activity")
	checkErr(err, "Select failed")
	newmap := map[string]interface{}{"metatitle": "this is my custom title", "activities": activities}
	r.HTML(200, "posts", newmap)

	var activity Activity
	err = dbmap.SelectOne(&activity, "select * from t_activty where id =?", args["id"])
	//simple error check
	if err != nil {
		newmap := map[string]interface{}{"metatitle": "404 Error", "message": "This is not found"}
		r.HTML(404, "error", newmap)
	} else {
		newmap := map[string]interface{}{"metatitle": activity.Title + " more custom", "activity": activity}
		r.HTML(200, "activity", newmap)
	}
}

func NewActivity(activity Activity, r render.Render, dbmap *gorp.DbMap) {
	log.Println(activity)
	err := dbmap.Insert(&activity)
	checkErr(err, "Insert failed")
	newmap := map[string]interface{}{"metatitle": "created activity", "activity": activity}
	r.HTML(200, "post", newmap)
}

//登陆
func LoginOAuth(oauth OAuth, r render.Render, dbmap *gorp.DbMap) {
	err := dbmap.SelectOne(&oauth, "select * from t_oauth where openid=?", oauth.OpenId)
	checkErr(err, "SelectOne failed")
	if err != nil {
		r.JSON(200, Resp{1000,"登陆失败",nil})
	}else {
		dbmap.Update(&oauth)
		obj, err := dbmap.Get(User{}, oauth.UserId)
		if err != nil {
			log.Fatal(err)
		}
		if obj == nil{
			r.JSON(200, Resp{1004,"用户资料信息不存在",nil})
		}
		user := obj.(*User)
		r.JSON(200, Resp{0,"登陆成功",map[string]interface{}{"token":oauth.AccessToken,"user":user}})
	}
}

func RegisterOAuth(register OAuthUser, r render.Render, dbmap *gorp.DbMap) {
	var oauth OAuth
	err := dbmap.SelectOne(&oauth, "select * from t_oauth where openid=?", register.OAuth.OpenId)
	checkErr(err, "SelectOne failed")
	if err != nil && oauth.OpenId != register.OAuth.OpenId {
		// err = dbmap.Insert(&register.User)
		// err = dbmap.Insert(&register.OAuth)
		trans, err := dbmap.Begin()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(register.User)
		trans.Insert(&register.User)
		register.OAuth.Expires = time.Now().Add(time.Second * time.Duration(register.OAuth.ExpiresIn))
		register.OAuth.UserId = register.User.Id
		log.Println( register.OAuth)
		trans.Insert(&register.OAuth)
		err = trans.Commit()
		if err == nil {
			r.JSON(200, Resp{0,"注册成功",map[string]interface{}{"user":register.User}})
		}else{
			log.Println("trans", err)
			r.JSON(200, Resp{1003,"注册失败",nil})
		}
	} else {
		log.Println("registered", err)
		r.JSON(200, Resp{1003,"该账号已经注册",nil})
	}
}


func Login(localAuth LocalAuth, r render.Render, dbmap *gorp.DbMap) {
	var auth LocalAuth
	err := dbmap.SelectOne(&auth, "select * from t_local_auth where phone=? and password = ?", localAuth.Phone,localAuth.Password)
	checkErr(err, "SelectOne failed")
	if err != nil {
		r.JSON(200, "账户或密码错误")
	}else {
		auth.Token = generaToken()
		auth.Expires = time.Now().Add(time.Hour * 24*30)
		dbmap.Update(&auth)
		obj, err := dbmap.Get(User{}, auth.UserId)
		if err != nil {
			log.Fatal(err)
		}
		if obj == nil{
			r.JSON(200, Resp{1004,"用户资料信息不存在",nil})
		}
		user := obj.(*User)
		r.JSON(200, Resp{0,"登陆成功",map[string]interface{}{"token":auth.Token,"user":user}})
	}
}

func Register(authUser LocalAuthUser, r render.Render, dbmap *gorp.DbMap) {
	var auth LocalAuth
	err := dbmap.SelectOne(&auth, "select * from t_local_auth where phone=?", authUser.LocalAuth.Phone)
	checkErr(err, "SelectOne failed")
	if err != nil && auth.Phone == "" {
		// err = dbmap.Insert(&authUser.LocalAuth)
		trans, err := dbmap.Begin()
		if err != nil {
			log.Fatal(err)
		}
		err = trans.Insert(&authUser.User)
		if err != nil {
			log.Fatal(err)
		}
		authUser.LocalAuth.UserId = authUser.User.Id
		log.Println(authUser.LocalAuth)
		err = trans.Insert(&authUser.LocalAuth)
		if err != nil {
			log.Fatal(err)
		}
		err = trans.Commit()
		if err != nil {
			log.Fatal(err)
		}
		if err == nil {
			r.JSON(200, Resp{0,"注册成功",map[string]interface{}{"token":auth.Token,"user":authUser.User}})
		}else{
			log.Fatal(err)
			log.Println("trans", err)
			r.JSON(200, Resp{1005,"注册失败",nil})
		}
	} else {
		log.Println("registered", err)
		r.JSON(200, Resp{1006,"该账号已经注册",nil})
	}

}

func UpdateUser(user User, dbmap *gorp.DbMap) {
	dbmap.Update(&user)
}

func DeleteUser(args martini.Params, dbmap *gorp.DbMap) {
	_, err := dbmap.Exec("delete from t_user where id=?", args["id"])
	checkErr(err, "Exec failed")
}

func GetUser(args martini.Params, r render.Render, dbmap *gorp.DbMap) {
	var user User
	err := dbmap.SelectOne(&user, "select * from t_user where id=?", args["id"])
	//simple error check
	if err != nil {
		r.JSON(200, "not found")
	} else {
		newmap := map[string]interface{}{"metatitle": user.Name + " more custom", "user": user}
		r.JSON(200, newmap)
	}
}

func main() {
	//setup db
	dbmap := initDb()
	defer dbmap.Db.Close()

	m := martini.Classic()

	m.Map(dbmap)

	m.NotFound(func() {
		// 处理 404
	})
	m.Use(func(res http.ResponseWriter, req *http.Request) {
		//		if req.Header.Get("X-API-KEY") != "secret123" {
		//			res.WriteHeader(http.StatusUnauthorized)
		//		}
	})
	m.Use(render.Renderer())

	m.Group("/activity", func(r martini.Router) {
		r.Get("/", GetActivityList)
		r.Get("/:id", GetActivity)
		r.Post("/", NewActivity)
		//		r.Put("/update/:id", UpdateActivity)
		//		r.Delete("/delete/:id", DeleteActivity)
	})

	m.Group("/user", func(r martini.Router) {
		r.Post("/oauth/login",binding.Bind(OAuth{}), LoginOAuth)
		r.Post("/oauth/register",binding.Bind(OAuthUser{}), RegisterOAuth)
		r.Post("/login", binding.Bind(LocalAuth{}), Login)
		r.Post("/register", binding.Bind(LocalAuthUser{}), Register)
		r.Get("/:id", GetUser)
		r.Put("/:id", binding.Bind(User{}), UpdateUser)
		r.Delete("/:id", DeleteUser)
	})

	m.Run()
}

func initDb() *gorp.DbMap {
	db, err := sql.Open("mysql", "root:weiwanglive@tcp(106.75.19.205:3306)/minnetlive?parseTime=true")
	checkErr(err, "sql.Open failed")
	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	dbmap.AddTableWithName(User{}, "t_user").SetKeys(true, "Id")
	dbmap.AddTableWithName(OAuth{}, "t_oauth").SetKeys(true, "Id")
	dbmap.AddTableWithName(LocalAuth{}, "t_local_auth").SetKeys(true, "Id")

	//TODO cdreate table
	//	// add a table, setting the table name to 'posts' and
	//	// specifying that the Id property is an auto incrementing PK
	//	dbmap.AddTableWithName(Post{}, "posts").SetKeys(true, "Id")
	//	// create the table. in a production system you'd generally
	//	// use a migration tool, or create the tables via scripts
	//	err = dbmap.CreateTablesIfNotExists()
	//	checkErr(err, "Create tables failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		// log.Fatalln(msg, err)
		log.Println(msg,err)
	}
}

func generaToken() string {

	return "genToken"
}
