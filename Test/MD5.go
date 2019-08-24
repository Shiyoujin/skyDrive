package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func main() {
	h := md5.New()
	h.Write([]byte("123456")) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)
	fmt.Println(cipherStr)
	h.Write([]byte("987"))
	aa := h.Sum(nil)

	fmt.Println(aa)
	fmt.Println(hex.EncodeToString(aa))
	fmt.Printf("%s\n", hex.EncodeToString(cipherStr)) // 输出加密结果

}
