package main

import (
	"database/sql"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	//	"gopkg.in/gorp.v1"
	// "fmt"
	"github.com/coopernurse/gorp"
	"github.com/pborman/uuid"
	"io/ioutil"
	"ucloud"
)

type Resp struct {
	Ret  int64       `form:"ret" json:"ret"`
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
	Id       int       `form:"id"` //  `form:"id"  db:"id,primarykey, autoincrement"`
	UserId   int       `form:"userId"  db:"user_id"`
	Phone    string    `form:"phone" binding:"required"  db:"phone"`
	Password string    `form:"password" binding:"required" db:"password"`
	Token    string    `db:"token"`
	Expires  time.Time `db:"expires"`
}

type OAuthUser struct {
	User  User
	OAuth OAuth
}

type LocalAuthUser struct {
	User      User
	LocalAuth LocalAuth
}

type Activity struct {
	Id             int       `form:"id"  db:"id"`
	Title          string    `form:"title"  binding:"required" db:"title"`
	Date           time.Time `db:"date"`
	ADate          int64     `form:"date" binding:"required" db:"-"`
	Desc           string    `form:"desc" binding:"required" db:"desc"`
	FontCover      string    `form:"fontCover" binding:"required" db:"front_cover"`
	Type           int       `form:"type" binding:"required" db:"type"`
	Price          int       `form:"price"  db:"price"`
	Password       string    `form:"password"  db:"pwd"`
	BelongUserId   int       `form:"belongUserId"  db:"belong_user_id"`
	VideoId        string    `form:"videoId"  db:"video_id"`
	VideoType      int       `form:"videoType" binding:"required" db:"video_type"`
	VideoPullPath  string    `form:"videoPullPath"  db:"video_pull_path"`
	VideoPushPath  string    `form:"videoPushPath"  db:"video_push_path"`
	VideoStorePath string    `form:"videoPushPath"  db:"video_store_path"`
}

func GetActivityList(r render.Render, dbmap *gorp.DbMap) {
	var activities []Activity
	_, err := dbmap.Select(&activities, "select * from t_activity")
	checkErr(err, "Select failed")
	if err != nil {
		r.JSON(200, Resp{2002, "查询活动失败", nil})
	} else {
		r.JSON(200, Resp{0, "查询活动成功", activities})
	}
}

func GetActivity(args martini.Params, r render.Render, dbmap *gorp.DbMap) {
	var activity Activity
	err := dbmap.SelectOne(&activity, "select * from t_activity where id =?", args["id"])
	checkErr(err, "Select failed")
	if err != nil {
		r.JSON(200, Resp{2003, "活动不存在", nil})
	} else {
		r.JSON(200, Resp{0, "查询活动成功", activity})
	}
}

func NewActivity(activity Activity, r render.Render, dbmap *gorp.DbMap) {
	activity.VideoId = uuid.New()
	activity.VideoPushPath = "xxxxxx"
	activity.BelongUserId = 123
	//奇葩
	// activity.Date = time.Unix(activity.ADate, 0).Format("2006-01-02 15:04:05")
	activity.Date = time.Unix(activity.ADate, 0)
	err := dbmap.Insert(&activity)
	checkErr(err, "Insert failed")
	if err != nil {
		r.JSON(200, Resp{2001, "添加活动失败", nil})
	} else {
		r.JSON(200, Resp{0, "添加活动成功", activity})
	}
}

func UpdateActivity(args martini.Params, activity Activity, r render.Render, dbmap *gorp.DbMap) {
	obj, err := dbmap.Get(Activity{}, args["id"])
	if err != nil {
		log.Println(err)
		r.JSON(200, Resp{2004, "更新活动失败", nil})
	} else {
		orgActivity := obj.(*Activity)
		orgActivity.Title = activity.Title
		orgActivity.Date = activity.Date
		orgActivity.Desc = activity.Desc
		orgActivity.Type = activity.Type
		orgActivity.VideoType = activity.VideoType
		orgActivity.FontCover = activity.FontCover
		_, err := dbmap.Update(orgActivity)
		checkErr(err, "Update failed")
		if err != nil {
			r.JSON(200, Resp{2004, "更新活动失败", nil})
		} else {
			r.JSON(200, Resp{0, "更新活动成功", activity})
		}
	}
}

func DeleteActivity(args martini.Params, r render.Render, dbmap *gorp.DbMap) {
	_, err := dbmap.Exec("DELETE from t_activity WHERE id=?", args["id"])
	checkErr(err, "Delete failed")
	if err != nil {
		r.JSON(200, Resp{2005, "删除活动失败", nil})
	} else {
		r.JSON(200, Resp{0, "删除活动成功", nil})
	}
}

//登陆
func LoginOAuth(oauth OAuth, r render.Render, dbmap *gorp.DbMap) {
	err := dbmap.SelectOne(&oauth, "select * from t_oauth where openid=?", oauth.OpenId)
	checkErr(err, "SelectOne failed")
	if err != nil {
		r.JSON(200, Resp{1000, "登陆失败", nil})
	} else {
		dbmap.Update(&oauth)
		obj, err := dbmap.Get(User{}, oauth.UserId)
		if err != nil {
			log.Fatal(err)
		}
		if obj == nil {
			r.JSON(200, Resp{1004, "用户资料信息不存在", nil})
		}
		user := obj.(*User)
		r.JSON(200, Resp{0, "登陆成功", map[string]interface{}{"token": oauth.AccessToken, "user": user}})
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
		log.Println(register.OAuth)
		trans.Insert(&register.OAuth)
		err = trans.Commit()
		if err == nil {
			r.JSON(200, Resp{0, "注册成功", map[string]interface{}{"user": register.User}})
		} else {
			log.Println("trans", err)
			r.JSON(200, Resp{1003, "注册失败", nil})
		}
	} else {
		log.Println("registered", err)
		r.JSON(200, Resp{1003, "该账号已经注册", nil})
	}
}

func Login(localAuth LocalAuth, r render.Render, dbmap *gorp.DbMap) {
	var auth LocalAuth
	err := dbmap.SelectOne(&auth, "select * from t_local_auth where phone=? and password = ?", localAuth.Phone, localAuth.Password)
	checkErr(err, "SelectOne failed")
	if err != nil {
		r.JSON(200, "账户或密码错误")
	} else {
		auth.Token = generaToken()
		auth.Expires = time.Now().Add(time.Hour * 24 * 30)
		dbmap.Update(&auth)
		obj, err := dbmap.Get(User{}, auth.UserId)
		if err != nil {
			log.Fatal(err)
		}
		if obj == nil {
			r.JSON(200, Resp{1004, "用户资料信息不存在", nil})
		}
		user := obj.(*User)
		r.JSON(200, Resp{0, "登陆成功", map[string]interface{}{"token": auth.Token, "user": user}})
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
		authUser.LocalAuth.Expires = time.Now().Add(time.Hour * 24 * 30)
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
			r.JSON(200, Resp{0, "注册成功", map[string]interface{}{"token": auth.Token, "user": authUser.User}})
		} else {
			log.Fatal(err)
			log.Println("trans", err)
			r.JSON(200, Resp{1005, "注册失败", nil})
		}
	} else {
		log.Println("registered", err)
		r.JSON(200, Resp{1006, "该账号已经注册", nil})
	}

}

func UpdateUser(args martini.Params, user User, r render.Render, dbmap *gorp.DbMap) {
	obj, err := dbmap.Get(User{}, args["id"])
	if err != nil {
		log.Println(err)
		r.JSON(200, Resp{2004, "更新信息失败", nil})
	} else {
		orgUser := obj.(*User)
		if user.Name != "" {
			orgUser.Name = user.Name
		}
		if user.Avatar != "" {
			orgUser.Avatar = user.Avatar
		}
		_, err := dbmap.Update(orgUser)
		checkErr(err, "Update failed")
		if err != nil {
			r.JSON(200, Resp{2004, "更新信息失败", nil})
		} else {
			r.JSON(200, Resp{0, "更新信息成功", user})
		}
	}
}

func DeleteUser(args martini.Params, r render.Render, dbmap *gorp.DbMap) {
	_, err := dbmap.Exec("DELETE from t_user WHERE id=?", args["id"])
	checkErr(err, "Delete failed")
	if err != nil {
		r.JSON(200, Resp{2005, "删除用户失败", nil})
	} else {
		r.JSON(200, Resp{0, "删除用户成功", nil})
	}
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
		r.Post("", binding.Bind(Activity{}), NewActivity)
		r.Get("", GetActivityList)
		r.Get("/:id", GetActivity)
		r.Put("/:id", binding.Bind(Activity{}), UpdateActivity)
		r.Delete("/:id", DeleteActivity)
	})

	m.Group("/user", func(r martini.Router) {
		r.Post("/oauth/login", binding.Bind(OAuth{}), LoginOAuth)
		r.Post("/oauth/register", binding.Bind(OAuthUser{}), RegisterOAuth)
		r.Post("/login", binding.Bind(LocalAuth{}), Login)
		r.Post("/register", binding.Bind(LocalAuthUser{}), Register)
		r.Get("/:id", GetUser)
		r.Put("/:id", binding.Bind(User{}), UpdateUser)
		r.Delete("/:id", DeleteUser)
	})
	m.Post("/upload", Upload)
	m.Group("/amdin", func(r martini.Router) {
		r.Get("/", AdminMain)
	})

	m.Run()
}

func AdminMain(r render.Render) {

	r.HTML(200, "hello", "amdin")
}

func initDb() *gorp.DbMap {
	db, err := sql.Open("mysql", "root:weiwanglive@tcp(106.75.19.205:3306)/minnetlive?parseTime=true")
	checkErr(err, "sql.Open failed")
	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	dbmap.AddTableWithName(User{}, "t_user").SetKeys(true, "Id")
	dbmap.AddTableWithName(OAuth{}, "t_oauth").SetKeys(true, "Id")
	dbmap.AddTableWithName(LocalAuth{}, "t_local_auth").SetKeys(true, "Id")
	dbmap.AddTableWithName(Activity{}, "t_activity").SetKeys(true, "Id")

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

func Upload(r *http.Request, render render.Render) {
	log.Println("parsing form")
	err := r.ParseMultipartForm(100000)
	// checkErr("upload ParseMultipartForm",err)
	if err != nil {
		render.JSON(500, "server err")
	}
	file, head, err := r.FormFile("file")
	checkErr(err, "upload Fromfile")
	log.Println(head.Filename)
	defer file.Close()
	tempDir := "/Users/cainli/dev/mininetlvie/temp/"
	filepath := tempDir + head.Filename
	fW, err := os.Create(tempDir + head.Filename)
	checkErr(err, "create file error")
	defer fW.Close()
	_, err = io.Copy(fW, file)
	checkErr(err, "create file error")
	fileUUID := uuid.New()
	err = UploadToUCloudCND(filepath,fileUUID, render)
	if err == nil {
		render.JSON(200, Resp{0, "上传成功", "token :"+ fileUUID })
	} else {
		render.JSON(200, Resp{3001, "上传失败",nil})
	}
}

func UploadToUCloudCND(path string, fileName string, render render.Render) error {
	publicKey := "enqyjAgoDAQm0mx6A/xk8eyxEuEJWK+LQ6n258NtsT6lARMyF+YFgA=="
	privateKey := "2e3da80f079d3362f504a5db3776a9cd41feeea2"
	u := ucloud.NewUcloudApiClient(
		publicKey,
		privateKey,
	)
	contentType := "image/jpeg"
	bucketName := "mininetlive"
	data, err := ioutil.ReadFile(path)
	checkErr(err, "ReadFile")
	resp, err := u.PutFile(fileName, bucketName, contentType, data)
	checkErr(err,"upload ucloud")
	if err == nil {
		log.Println(resp.StatusCode)
		log.Println(string(resp.Content))
	} 
	return err
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Println(msg, err)
		// log.Fatalln(msg, err)
	}
}

func generaToken() string {
	return "genToken"
}
