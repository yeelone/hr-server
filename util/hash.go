package util

import (
	"github.com/speps/go-hashids"
	"github.com/spf13/viper"
)

var salt = viper.GetString("hasid_salt")

//Encode 混淆
func HashIDEncode(data string) string {
	hd := hashids.NewData()
	hd.Salt = salt
	h, _ := hashids.NewWithData(hd)
	e, _ := h.EncodeHex("123456789")
	return e
}

//Decode 还原混淆
func HashIDDecode(data string) int {
	hd := hashids.NewData()
	hd.Salt = salt
	h, _ := hashids.NewWithData(hd)
	e, _ := h.DecodeWithError(data)
	return e[0]
}


