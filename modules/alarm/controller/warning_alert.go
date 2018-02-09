package controller

import (
	"cloudtropy.com/alert/libofm"
	"encoding/json"
	// "errors"
	"io/ioutil"
	"net/http"
	// "strings"
	// "time"
)

func PostAlert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		RenderMsgJson(w, err.Error())
		return
	}
	defer r.Body.Close()

	var b map[string]string
	err = json.Unmarshal(body, &b)
	if err != nil {
		RenderMsgJson(w, "InvalidRequestParams")
		return
	} else if b["user"] == "" || b["way"] == "" {
		RenderMsgJson(w, "InvalidRequestParams")
		return
	} else if b["way"] != "email" && b["way"] != "wechat" && b["way"] != "phone" {
		RenderMsgJson(w, "InvalidRequestParams")
		return
	}

	switch b["way"] {
	case "email":
		if b["email"] == "" {
			RenderMsgJson(w, "InvalidEmailAddr")
			return
		}
		email := libofm.NewEmail(b["email"], b["theme"], b["content"])
		err := libofm.SendEmail(email)
		if err != nil {
			RenderMsgJson(w, "SendEmailWarning Failed")
		}
		// SendEmail(b["email"], b["theme"], b["content"])
	case "wechat":
		if b["wechat"] == "" {
			RenderMsgJson(w, "InvalidUser")
			return
		}
		err := libofm.Wechat_warning(b["touser"], b["agentid"], b["content"])
		if err != nil {
			RenderMsgJson(w, "SendWechat Failed")
		}
		// SendWechat(b["wechat"], b["content"])
	case "phone":
		if b["phone"] == "" {
			RenderMsgJson(w, "InvalidUserPhoneNumber")
			return
		}
		err := libofm.NoticePersonByVoice(b["phone"], b["sex"], b["device"], b["info"])
		if err != nil {
			RenderMsgJson(w, "SendPhone Warning Failed")
		}
		// SendPhone(users[user]["phone"])
	}
	RenderMsgJson(w, "success")

}

func PostEmailWarning(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		RenderMsgJson(w, err.Error())
		return
	}

	var mss map[string]string
	err = json.Unmarshal(body, &mss)
	if err != nil {
		RenderMsgJson(w, "InvalidRequestBody: "+err.Error())
		return
	}

	if mss["emailAddr"] == "" {
		RenderMsgJson(w, "InvalidRequestBody(Empty Email Addr)")
		return
	}

	email := libofm.NewEmail(mss["emailAddr"], mss["title"], mss["message"])

	err = libofm.SendEmail(email)

	if err != nil {
		RenderMsgJson(w, "SendEmailWarning Failed")
	}

	RenderMsgJson(w, "success")

}

func PostWechatWarning(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		RenderMsgJson(w, err.Error())
		return
	}

	var mss map[string]string
	err = json.Unmarshal(body, &mss)
	if err != nil {
		RenderMsgJson(w, "InvalidRequestBody: "+err.Error())
		return
	}

	err = libofm.Wechat_warning(mss["touser"], mss["agentid"], mss["content"])

	if err != nil {
		RenderMsgJson(w, "SendWechatWarning Failed")
	}

	RenderMsgJson(w, "success")

}
