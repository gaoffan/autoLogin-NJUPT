package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

var file,_ = os.OpenFile("autologin.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
var logger = log.New(file, "", log.LstdFlags)

var query,_ = url.Parse("http://p.njupt.edu.cn:801/eportal/?")
var xyw = "http://192.168.168.168/0.htm"

type config struct {
	WlanUserIp string
	WlanAcIp string
	WlanAcName string
}

type auth struct {
	username *string
	password *string
}

func main()  {
	auth1 := auth{
		username: flag.String("u","null","username, example: B20010101@cmcc"),
		password: flag.String("p","null","password"),
	}
	auth2 := auth{
		username: flag.String("v","null","secondary username, example: B20010101@cmcc"),
		password: flag.String("q","null","secondary password"),
	}
	flag.Parse()

	if *auth1.username == "null" || *auth1.password == "null"{
		panic("username or password is blank.")
	}
	for {
		cnf,err := getConfig()
		// 6.6.6.6 这个网站会在登录成功后一直超时，那么可知，没有超时就是未登录状态
		if err == nil{
			logger.Println("使用账号登录 ",*auth1.username)
			err1 := login(cnf,&auth1)
			if err1 != nil && *auth2.username != "null"{
				logger.Println(err1)
				logger.Println("使用账号登录 ",*auth2.username)
				err2 := login(cnf,&auth2)
				if err2 != nil{
					logger.Println(err2)
				}else{
					logger.Println("登录成功 ",*auth2.username)
				}
			}else{
				logger.Println("登录成功 ",*auth1.username)
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func getConfig() (*config,error){
	client := http.DefaultClient
	client.Timeout = 1 * time.Second
	resp, err := client.Get("http://6.6.6.6/")
	if err != nil{
		return nil, err
	}
	data, _ := ioutil.ReadAll(resp.Body)
	compile, _ := regexp.Compile("location.href=\\\"(http://.*)\\\"")
	submatch := compile.FindStringSubmatch(string(data))
	parse, _ := url.Parse(submatch[1])
	cnf := config{
		WlanUserIp: parse.Query().Get("wlanuserip"),
		WlanAcIp: parse.Query().Get("wlanacip"),
		WlanAcName: parse.Query().Get("wlanacname"),
	}
	return &cnf,nil
}

func login(cnf *config, auth *auth) error{
	form := url.Values{}
	form.Set("DDDDD", ",0," + *auth.username)
	form.Set("upass", *auth.password)
	form.Set("R1", "0")
	form.Set("R2", "0")
	form.Set("R3", "0")
	form.Set("R6", "0")
	form.Set("para", "00")
	form.Set("0MKKey", "123456")
	form.Set("buttonClicked", "")
	form.Set("redirect_url", "")
	form.Set("err_flag", "")
	form.Set("username", "")
	form.Set("password", "")
	form.Set("user", "")
	form.Set("cmd", "")
	form.Set("Login", "")
	form.Set("v6ip", "")

	params := url.Values{}
	params.Add("c","ACSetting")
	params.Add("a","Login")
	params.Add("protocol","http:")
	params.Add("hostname","p.njupt.edu.cn")
	params.Add("iTermType","1")
	params.Add("wlanuserip",cnf.WlanUserIp)
	params.Add("wlanacip",cnf.WlanAcIp)
	params.Add("wlanacname",cnf.WlanAcName)
	params.Add("mac","00-00-00-00-00-00")
	params.Add("ip",cnf.WlanUserIp)
	params.Add("enAdvert","0")
	params.Add("queryACIP","0")
	params.Add("loginMethod","1")

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	post, err := client.PostForm(query.String() + params.Encode(),form)
	if err != nil{
		return err
	}

	if strings.Contains(post.Header.Get("Location"),"3.htm"){
		return nil
	}
	return errors.New("登录失败")
}
