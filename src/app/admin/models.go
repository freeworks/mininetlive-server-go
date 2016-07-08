package admin

import (
	. "app/common"
	logger "app/logger"
	"app/sessionauth"
	"fmt"
	"time"

	"github.com/coopernurse/gorp"
)

type AdminModel struct {
	Id            int64     `form:"id" db:"id"`
	Uid           string    ` db:"uid"`
	Phone         string    `form:"phone" db:"phone"`
	NickName      string    `form:"nickName" db:"nickname"`
	Password      string    `form:"password" db:"password"`
	Avatar        string    `form:"avatar"  db:"avatar"`
	EasemobUUID   string    `form:"-"  db:"easemob_uuid"`
	Updated       time.Time `db:"update_time"`
	Created       time.Time `db:"create_time"`
	Authenticated bool      `form:"-" db:"-"`
}

func (admin *AdminModel) String() string {
	adminString := fmt.Sprintf("[%s, %s, %d]", admin.Id, admin.NickName, admin.Password)
	logger.Info(adminString)
	return adminString
}

func (admin *AdminModel) Login() {
	logger.Info("login ....")
	admin.Authenticated = true
}

func (admin *AdminModel) Logout() {
	admin.Authenticated = false
}

func (admin *AdminModel) IsAuthenticated() bool {
	return admin.Authenticated
}

func (admin *AdminModel) UniqueId() interface{} {
	return admin.Uid
}

func (admin *AdminModel) GetById(uid interface{}, dbmap *gorp.DbMap) error {
	logger.Info("GetById:", uid)
	err := dbmap.SelectOne(admin, "SELECT * FROM t_admin WHERE uid = ?", uid)
	CheckErr(err, "GetById select one")
	if err != nil {
		return err
	}
	return nil
}

func GenerateAnonymousUser() sessionauth.User {
	return &AdminModel{}
}
