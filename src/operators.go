package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

//CheckAuth via DB
func CheckAuth(key string) (int, AuthKey) {

	/*
		Auth
		0 --> Wrong AuthKey
		1 --> Correct AuthKey
		2 --> Out of date Authkey

		Type
		0 --> User
		1 --> Admin
	*/

	Localkey := GetAuthByKey(key)

	if Localkey.AUTHKEY == key {
		dbDate, err := time.Parse("2006-01-02 15:04:05", Localkey.TIMESTAMP)

		if err != nil {
			log.Panic(err)
		} else {
			days := time.Now().Sub(dbDate)
			if 22-days.Hours() <= 0 {
				return 2, Localkey
			} else {
				return 1, Localkey
			}
		}

	} else {
		return 0, Localkey
	}

	return 0, Localkey
}

func CheckLoginUser(email string, password string) (map[string]interface{}, bool, int) {

	resp := map[string]interface{}{}
	user, status := GetUserByMail(email)

	if status == true {
		//Check Password
		correct := bcrypt.CompareHashAndPassword([]byte(user.PASSWORD), []byte(password))

		if correct != nil {
			resp["status"] = "Wrong password!"
			return resp, false, -1

		}
		//Generate AuthKey
		auth := GetAuthByUser(user.ID)
		if auth.USERID == user.ID {
			//Check authdate
			dbDate, err := time.Parse("2006-01-02 15:04:05", auth.TIMESTAMP)

			if err != nil {
				log.Panic(err)
			} else {
				days := time.Now().Sub(dbDate)
				log.Println(fmtDuration(days))
				if 22-days.Hours() <= 0 {
					key := GenerateAuthKey(user)
					_, status := UpdateNewAuthkey(user, key)
					resp["key"] = key
					return resp, status, user.TYPE
				} else {
					//return key

					resp["key"] = auth.AUTHKEY
					resp["authdate"] = auth.TIMESTAMP
					resp["time_remaining"] = strconv.FormatFloat(22-days.Hours(), 'f', 0, 64) + " hours"
					return resp, true, user.TYPE
				}
			}

		} else {
			//create new one
			key := GenerateAuthKey(user)
			_, status := InsertNewAuthkey(user, key)
			resp["key"] = key
			return resp, status, user.TYPE
		}

	} else {
		//Not exists!
		resp["status"] = "User Not Exists!"
		return resp, false, -1
	}

	return resp, false, -1

}

func GenerateAuthKey(user User) string {
	localDate := time.Now()
	randomUID := RandStringBytesMaskImpr(50)

	var md1 = MD5(localDate.Format(time.RFC850))
	var md2 = MD5(randomUID)
	var md3, _ = json.Marshal([]byte(user.TIMESTAMP))
	var md4 = MD5(string(md3))
	var md5 = MD5(md1 + md2 + md4)
	return md5
}

func RandStringBytesMaskImpr(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
