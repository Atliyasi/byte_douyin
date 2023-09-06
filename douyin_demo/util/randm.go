package util

import (
	"math/rand"
	"time"
)

func RandInt() int {
	rand.Seed(time.Now().UnixNano())

	// 生成1到9之间的随机整数
	randomNumber := rand.Intn(9) + 1

	return randomNumber
}
