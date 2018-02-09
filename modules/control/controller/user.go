package controller

import (
  "net/http"
  "io/ioutil"
  "encoding/json"
  "time"
  // "strings"

  "github.com/cloudtropy/open-operation/modules/control/ctx"
  "github.com/cloudtropy/open-operation/modules/control/db/mysql"
  "github.com/cloudtropy/open-operation/modules/control/db/redis"
  "github.com/cloudtropy/open-operation/utils/fun"
  "github.com/tidwall/gjson"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.NotFound(w, r)
    return
  }

  defer r.Body.Close()
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    HttpResMsg(w, "InternalError", err.Error())
    return
  }

  user := mysql.User{}
  err = json.Unmarshal(body, &user)
  if err != nil {
    HttpResMsg(w, "InvalidRequestParams", err.Error())
    return
  } else if (user.Sex != "male" && user.Sex != "female") ||
    user.User == "" || user.Passwd == "" {
    HttpResMsg(w, "InvalidRequestParams", "")
    return
  }

  user.Salt = fun.GetUUIDV4()[:24]
  user.Passwd = fun.GetMd5(user.Passwd + user.Salt)

  err = mysql.CreateUser(&user)
  if err != nil {
    log.Println("mysql.CreateUser", user, err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }
  HttpResMsg(w, "Success", "")
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.NotFound(w, r)
    return
  }

  defer r.Body.Close()
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    HttpResMsg(w, "InternalError", err.Error())
    return
  }

  userId := gjson.GetBytes(body, "id").Int()
  if userId <= 0 || !fun.IsSqlId(gjson.GetBytes(body, "id").String()) {
    HttpResMsg(w, "InvalidRequestParams", "")
    return
  }

  err = mysql.DeleteUser(userId)
  if err != nil {
    log.Println("mysql.DeleteUser", userId, err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }
  HttpResMsg(w, "Success", "")
}

func ResetUserPassword(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.NotFound(w, r)
    return
  }

  defer r.Body.Close()
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    HttpResMsg(w, "InternalError", err.Error())
    return
  }

  var mss map[string]string
  err = json.Unmarshal(body, &mss)
  if err != nil {
    HttpResMsg(w, "InvalidRequestParams", err.Error())
    return
  }

  if mss["user"] == "" || mss["passwd"] == "" {
    HttpResMsg(w, "InvalidRequestParams", "user or passwd")
    return
  }

  user, err := mysql.GetUser(mss["user"])
  if err != nil {
    log.Println("mysql.GetUser", mss["user"], err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  } else if user == nil {
    HttpResMsg(w, "InvalidUser", "")
    return
  }

  passwd := fun.GetMd5(mss["passwd"] + user.Salt)
  err = mysql.UpdatePassword(mss["user"], passwd)
  if err != nil {
    log.Println("mysql.UpdatePassword", mss["user"], passwd, err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  HttpResMsg(w, "Success", "")
}

func UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.NotFound(w, r)
    return
  }

  defer r.Body.Close()
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    HttpResMsg(w, "InternalError", err.Error())
    return
  }

  username := GetUsername(r)
  if username == "" {
    HttpResMsg(w, "NotLogin", "")
    return
  }

  var mss map[string]string
  err = json.Unmarshal(body, &mss)
  if err != nil {
    HttpResMsg(w, "InvalidRequestParams", err.Error())
    return
  }
  if mss["passwd"] == "" || mss["newPasswd"] == "" {
    HttpResMsg(w, "InvalidRequestParams", "passwd or newPasswd")
    return
  }

  user, err := mysql.GetUser(username)
  if err != nil {
    log.Println("mysql.GetUser", username, err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  } else if user == nil {
    HttpResMsg(w, "InvalidUser", "")
    return
  }

  oldPasswd := fun.GetMd5(mss["passwd"] + user.Salt)
  if oldPasswd != user.Passwd {
    HttpResMsg(w, "InvalidPassword", "")
    return
  }
  
  newPasswd := fun.GetMd5(mss["newPasswd"] + user.Salt)
  err = mysql.UpdatePassword(username, newPasswd)
  if err != nil {
    log.Println("mysql.UpdatePassword", username, newPasswd, err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  HttpResMsg(w, "Success", "")
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.NotFound(w, r)
    return
  }

  defer r.Body.Close()
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    HttpResMsg(w, "InternalError", err.Error())
    return
  }

  var mss map[string]string
  err = json.Unmarshal(body, &mss)
  if err != nil {
    HttpResMsg(w, "InvalidRequestParams", err.Error())
    return
  }

  user, err := mysql.GetUser(mss["user"])
  if err != nil {
    log.Println("mysql.GetUser", mss["user"], err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  } else if user == nil {
    HttpResMsg(w, "InvalidUser", "")
    return
  }

  // todo: check
  user.Alias = mss["alias"]
  user.Email = mss["email"]
  user.Phone = mss["phone"]
  user.Wechat = mss["wechat"]
  user.Sex = mss["sex"]

  err = mysql.UpdateUser(user)
  if err != nil {
    log.Println("mysql.UpdateUser", user, err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }
  HttpResMsg(w, "Success", "")
}

func UpdateUserSelf(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.NotFound(w, r)
    return
  }

  username := GetUsername(r)
  if username == "" {
    HttpResMsg(w, "NotLogin", "")
    return
  }

  defer r.Body.Close()
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    HttpResMsg(w, "InternalError", err.Error())
    return
  }

  var mss map[string]string
  err = json.Unmarshal(body, &mss)
  if err != nil {
    HttpResMsg(w, "InvalidRequestParams", err.Error())
    return
  }

  user, err := mysql.GetUser(username)
  if err != nil {
    log.Println("mysql.GetUser", username, err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  } else if user == nil {
    HttpResMsg(w, "InvalidUser", "")
    return
  }

  // todo: check
  user.Alias = mss["alias"]
  user.Email = mss["email"]
  user.Phone = mss["phone"]
  user.Wechat = mss["wechat"]
  user.Sex = mss["sex"]

  err = mysql.UpdateUser(user)
  if err != nil {
    log.Println("mysql.UpdateUser", user, err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }
  HttpResMsg(w, "Success", "")
}

func GetUserList(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.NotFound(w, r)
    return
  }

  users, err := mysql.GetUsers()
  if err != nil {
    log.Println("mysql.GetUsers", err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  }

  res := make([]map[string]interface{}, 0)
  for _, user := range users {
    userInfo := map[string]interface{}{
      "id": user.Id,
      "user": user.User,
      "alias": user.Alias,
      "email": user.Email,
      "phone": user.Phone,
      "wechat": user.Wechat,
      "sex": user.Sex,
      "createTime": user.CreateTime,
    }
    res = append(res, userInfo)
  }

  HttpResData(w, res)
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.NotFound(w, r)
    return
  }

  r.ParseForm()
  username := r.FormValue("user")
  if username == "" {
    HttpResMsg(w, "InvalidRequestParams", "user")
    return
  }

  user, err := mysql.GetUser(username)
  if err != nil {
    log.Println("mysql.GetUser", username, err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  } else if user == nil {
    HttpResMsg(w, "InvalidUser", "")
    return
  }

  userInfo := map[string]interface{}{
    "id": user.Id,
    "user": user.User,
    "alias": user.Alias,
    "email": user.Email,
    "phone": user.Phone,
    "wechat": user.Wechat,
    "sex": user.Sex,
    "createTime": user.CreateTime,
  }
  HttpResData(w, userInfo)
}

func GetUserInfoSelf(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.NotFound(w, r)
    return
  }

  username := GetUsername(r)
  if username == "" {
    HttpResMsg(w, "NotLogin", "")
    return
  }

  user, err := mysql.GetUser(username)
  if err != nil {
    log.Println("mysql.GetUser", username, err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  } else if user == nil {
    HttpResMsg(w, "InvalidUser", "")
    return
  }

  userInfo := map[string]interface{}{
    "id": user.Id,
    "user": user.User,
    "alias": user.Alias,
    "email": user.Email,
    "phone": user.Phone,
    "wechat": user.Wechat,
    "sex": user.Sex,
    "createTime": user.CreateTime,
  }
  HttpResData(w, userInfo)
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.NotFound(w, r)
    return
  }

  defer r.Body.Close()
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    HttpResMsg(w, "InternalError", err.Error())
    return
  }

  var mss map[string]string
  err = json.Unmarshal(body, &mss)
  if err != nil {
    HttpResMsg(w, "InvalidRequestParams", err.Error())
    return
  }

  if mss["username"] == "" || mss["password"] == "" {
    HttpResMsg(w, "InvalidRequestParams", "username or password")
    return
  }
  var remember = false
  if mss["remember"] == "true" {
    remember = true
  }

  user, err := mysql.GetUser(mss["username"])
  if err != nil {
    log.Println("mysql.GetUser", mss["username"], err)
    HttpResMsg(w, "MysqlError", err.Error())
    return
  } else if user == nil {
    HttpResMsg(w, "InvalidUser", "")
    return
  }

  saltPw := fun.GetMd5(mss["password"] + user.Salt)
  if saltPw != user.Passwd {
    HttpResMsg(w, "InvalidPassword", "")
    return
  }

  // arrAccesses := make([]string, 0)
  // if mss["username"] == "root" {
  //   arrAccesses = append(arrAccesses, "AdministratorAccess")
  // } else {
  //   accesses, err := mysql.GetAccessesOfUser(mss["username"])
  //   if err != nil {
  //     log.Println("mysql.GetAccessesOfUser", mss["username"], err)
  //     HttpResMsg(w, "MysqlError", err.Error())
  //     return
  //   }
    
  //   for _, access := range accesses {
  //     arrAccesses = append(arrAccesses, access["access"])
  //   }
  // }

  // todo: model auth
  sid := fun.GetUUIDV4()
  err = redis.AddUserSession(sid, remember, map[string]string{
    "username": mss["username"],
    // "accesses": strings.Join(arrAccesses, ","),
  })
  if err != nil {
    log.Println("redis.AddUserSession", err)
    HttpResMsg(w, "RedisError", err.Error())
    return
  }

  var cookieMaxAge = int(redis.UserCookieExpire/time.Second) - 2
  if remember {
    cookieMaxAge = int(redis.UCRememberExpire/time.Second) - 2
  }

  SetResCookieKV(w, "sid", sid, cookieMaxAge)
  // SetResCookieKV(w, "username", mss["username"], 0)
  // SetResCookieKV(w, "accesses", strings.Join(arrAccesses, ","), 0)

  HttpResMsg(w, "Success", "")
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.NotFound(w, r)
    return
  }

  c, err := r.Cookie("sid") // c.Name is sid
  if err != nil {
    HttpResMsg(w, "NotLogin", "")
    return
  }

  err = redis.DelUserSession(c.Value)
  if err != nil {
    HttpResMsg(w, "RedisError", err.Error())
    return
  }

  SetResCookieKV(w, "sid", "logout", -1)
  // SetResCookieKV(w, "username", "logout", -1)
  // SetResCookieKV(w, "accesses", "logout", -1)

  HttpResMsg(w, "Success", "")
}

func GetUsername(r *http.Request) string {
  c, err := r.Cookie("sid")
  if err != nil {
    log.Println("GetUsername after check cookie sid error:", err)
    return ""
  }
  mssUser, err := redis.GetUserSession(c.Value)
  if err != nil {
    log.Println("redis.GetUserSession", c.Value, err)
    return ""
  }
  return mssUser["username"]
}

func SetResCookieKV(w http.ResponseWriter, name, value string, maxAge int) {
  // todo
  if maxAge == 0 {
    maxAge = int(time.Hour * 24 * time.Duration(ctx.Cfg().User.RememberTimeout) / time.Second)
  }
  http.SetCookie(w, &http.Cookie{
    Name:   name,
    Value:  value,
    // Domain: ctx.Cfg().CookieDomain,
    Path:   "/",
    MaxAge: maxAge,
  })
}
