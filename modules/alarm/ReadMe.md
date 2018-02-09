# **alert api**



               ____()()
              /      @@
        `~~~~~\_;m__m._>o 


## alert服务基本信息
* 域名：alert.cloudtropy.com
* 端口号：10010


## 报警平台 API 目录
1. alert上架相关api [(๑╹◡╹)ﾉ"""]

                  ___
         _  _  .-'   '-.
        (.)(.)/         \ 
         /@@             ;
        o_\\-mm-......-mm`~~~~~~~~~~~~~~~~`


---
---
## 接口目录
1. 邮件报警接口 [POST /email](#1)
2. 微信报警接口 [POST /wechat](#2)
3. 报警统一接口 [POST /alert](#3 阿里大于)

## 1. 邮件报警接口
* Path: Post /email
* 描述：通过邮件报警
* 返回Headers：
  * Content-Type: application/json
* 返回内容:

```javascript
// 异常时
{
    "msg": ""  // 代表异常信息
}
// 正常返回
{
  "msg": "" // "success"代表正常返回
}
```  

---
## 3. 报警统一接口
* Path: POST /alert
* 描述：报警统一接口
* 请求Headers: 
    * Content-Type: application/json
* 返回Headers:
    * Content-Type: application/json
* 请求body形式(3种操作情形)：

```json
{
    "user": String,             // 报警通知的用户
    "way": String,              // 报警的方式(三种可选): email || wechat || phone || all
    "content": String           // 报警内容
    "email": String,            // [可选] 邮件地址，way为email时使用
    "theme": String,            // [可选] 邮件主题，way为email时使用
    "touser": String            // [可选] 微信用户，way为wechat时使用
    "agentid": String           // [可选] 企业微信通知id，一般传1，way为wechat时使用
    "phone": String             // [可选] 电话号码，way为phone时使用
    "device": String            // [可选] 设备信息，way为phone时使用
    "sex": String               // [可选] 语音性别，way为phone时使用

}

```
* 返回body:

```javascript
{
    "msg": String,   //"success"代表正常返回，其他代表异常信息
}
```
    





