package libofm

import (
	"cloudtropy.com/alert/g"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	// "strconv"
	"strings"
	"time"
)

var (
	phoneUrl = "http://gw.api.taobao.com/router/rest"
	mCharChn = map[string]string{".": "点", "0": "零", "1": "一", "2": "二", "3": "三",
		"4": "四", "5": "五", "6": "六", "7": "七", "8": "八", "9": "九"}
	signBorder       = "5d254790d2bc32ca22ac86eb20536946"
	phoneVoiceParams = map[string]string{
		"app_key":         "23863876",
		"called_show_num": "02131314050",
		"format":          "json",
		"method":          "alibaba.aliqin.fc.tts.num.singlecall",
		"sign_method":     "md5",
		"v":               "2.0",
	}
)

func NoticePersonByVoice(callNum, sex, ip, info string) error {
	phoneVoiceParams["called_num"] = callNum
	phoneVoiceParams["timestamp"] = time.Now().String()[0:19]
	if sex == "girl" {
		phoneVoiceParams["tts_code"] = "TTS_70155332"
	} else {
		phoneVoiceParams["tts_code"] = "TTS_69985305"
	}
	ttsParam := map[string]string{
		"device": IpToChn(ip),
		"info":   info,
	}
	ttsParamBytes, err := json.Marshal(ttsParam)
	if err != nil {
		log.Println("NoticePersonByVoice", err, ttsParam)
		return err
	}
	phoneVoiceParams["tts_param"] = string(ttsParamBytes)

	signStrs := make([]string, 0)
	query := url.Values{}
	for k, v := range phoneVoiceParams {
		signStrs = append(signStrs, k+v)
		query.Add(k, v)
	}

	sort.Strings(signStrs)
	sign := g.GetMd5(signBorder + strings.Join(signStrs, "") + signBorder)
	query.Add("sign", sign)

	queryUrl := phoneUrl + "?" + query.Encode()
	log.Println(queryUrl)
	resp, err := http.Get(queryUrl)
	if err != nil {
		log.Println("Get", queryUrl, err)
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Get", queryUrl, err)
		return err
	} else {
		log.Printf("%s\n", body)
	}
	return nil
}

func IpToChn(ip string) (res string) {
	for i := 0; i < len(ip); i++ {
		res += mCharChn[string(ip[i])]
	}
	return
}
