package main

import (
	"flash-sale/common"
	"flash-sale/encrypt"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

// 集群地址
var hostArray = []string{"127.0.0.1", "127.0.0.1"}
var localHost = "127.0.0.1"
var port = "8081"
var hashConsistent *common.Consistent

// 控制信息
type AccessControl struct {
	// 用户想要存放的信息
	sourceArray map[int]interface{}
	sync.RWMutex
}
var accessControl = &AccessControl{
	sourceArray: make(map[int]interface{}),
}

func (a *AccessControl) GetNewRecord(uid int) interface{} {
	a.RWMutex.RLock()
	defer a.RWMutex.RUnlock()
	data := a.sourceArray[uid]
	return data
}

func (a *AccessControl) SetNewRecord(uid int) {
	a.RWMutex.Lock()
	defer a.RWMutex.Unlock()
	// FIXME: 修改
	a.sourceArray[uid] = "hello world"
}

func (a *AccessControl) GetDistributedRight(req *http.Request) bool {
	uid, err := req.Cookie("uid")
	if err != nil {
		return false
	}
	// 根据用户 id 获取具体机器
	hostRequest, err := hashConsistent.Get(uid.Value)
	if err != nil {
		return false
	}
	if hostRequest == localHost {
		// 本机数据校验
		return a.GetDataFromMap(uid.Value)
	}else {
		return GetDataFromOtherMap(hostRequest, req)
	}
}

func (a *AccessControl) GetDataFromMap(uid string)  bool {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}
	data := a.GetNewRecord(uidInt)
	if data != nil {
		return true
	}
	return false
}

func GetDataFromOtherMap(host string, req *http.Request) bool {
	uidPre, err := req.Cookie("uid")
	if err != nil {
		return false
	}
	uidSign, err := req.Cookie("sign")
	if err != nil {
		return false
	}
	client := &http.Client{}
	request, err := http.NewRequest("GET", "http://"+host+":"+port+"/access?uid=", nil)
	if err != nil {
		return false
	}
	// 排查多余 cookie
	cookieUid := &http.Cookie{Name:"uid", Value:uidPre.Value, Path:"/"}
	cookieSign := &http.Cookie{Name:"sign", Value:uidSign.Value, Path:"/"}
	request.AddCookie(cookieUid)
	request.AddCookie(cookieSign)

	res, err := client.Do(request)
	if err != nil {
		return false
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false
	}
	if res.StatusCode == 200 {
		if string(body) == "true" {
			return true
		}else {
			return false
		}
	}
	return false
}

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
	// 负载均衡 - 一致性 hash
	hashConsistent = common.NewConsistent()
	for _, v := range hostArray {
		hashConsistent.Add(v)
	}
	// 过滤器
	filter := common.NewFilter()
	filter.RegisterFilterUri("/check", Auth)
	http.HandleFunc("/check", filter.Handle(Check))
	_ = http.ListenAndServe(":8083" ,nil)
}
