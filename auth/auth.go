package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Result struct {
	Sta  int
	Code int
	Msg  string
	Data User
}
type User struct {
	ID          int       `gorm:"AUTO_INCREMENT"`
	OpenID      string    `gorm:"unique;size:255"` // string默认长度为255, 使用这种tag重设。
	NickName    string    `gorm:"size:255"`
	HeadImgURL  string    `gorm:"size:255"`
	Password    string    `gorm:"size:255"`       // string默认长度为255, 使用这种tag重设。
	Token       string    `gorm:"index;size:255"` // string默认长度为255, 使用这种tag重设。
	LastLogin   time.Time `gorm:"default:null"`
	Status      int       `gorm:"default:1"` // 1未激活2已激活3封禁
	ExpiredAt   time.Time `gorm:"default:null"`
	ExpiredAtTs int64     `gorm:"-"` //ignore this field
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var m = sync.Map{}

func Auth(apiAddress, token string) bool {
	defer recoverName()
	_, ok := m.Load(token)
	if ok {
		log.Println("auth cache")
		return true
	}
	resp, err := http.PostForm(apiAddress,
		url.Values{"token": {token}})

	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		log.Println(err)
	}
	result := Result{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("err on http:%#v", err)
		return false
	}
	if result.Sta == 0 {
		return false
	}
	if result.Data.ExpiredAt.Unix() < time.Now().Unix() {
		return false
	}
	log.Println("auth http")
	m.Store(token, true)
	return true
}

func recoverName() {
	if r := recover(); r != nil {
		fmt.Println("recovered from Login:", r)
	}
}
