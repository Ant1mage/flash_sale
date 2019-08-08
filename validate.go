package main

import (
	"flash-sale/common"
	"flash-sale/encrypt"
	"fmt"
	"net/http"
)

//执行正常业务逻辑
func Check(w http.ResponseWriter, r *http.Request) {
	fmt.Println("执行check！")
}

//统一验证拦截器
func Auth(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("执行验证！")
	//添加基于cookie的权限验证
	err := CheckUserInfo(r)
	if err != nil {
		return err
	}
	return nil
}

func CheckUserInfo(r *http.Request) error {
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		//return errors.New("用户UID Cookie 获取失败！")
	}

	signCookie, err := r.Cookie("sign")
	if err != nil {
		//return errors.New("用户加密串 Cookie 获取失败！")
	}

	signByte, err := encrypt.DePwdCode(signCookie.Value)
	if err != nil {
		//return errors.New("加密串已被篡改！")
	}
	if checkInfo(uidCookie.Value, string(signByte)) {
		//return nil
	}
	//return errors.New("身份校验失败！")
	return nil
}

func checkInfo(checkStr string, signStr string) bool {
	if checkStr == signStr {
		return true
	}
	return false
}

func main() {
	filter := common.NewFilter()
	filter.RegisterFilterUri("/check", Auth)
	http.HandleFunc("/check", filter.Handle(Check))
	_ = http.ListenAndServe(":8083" ,nil)
}
