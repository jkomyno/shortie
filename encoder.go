package main

import (
    "math"
    "strings"
)

const ALPHABET string = "23456789bcdfghjkmnpqrstvwxyzBCDFGHJKLMNPQRSTVWXYZ"
const BASE float64 = 49

func generateURL(n int64) string {
    s := ""
    var num float64 = float64(n)

    for num > 0 {
        s = string(ALPHABET[int(num) % int(BASE)])+s
        num = math.Floor(num / BASE)
	}

    return s
}

func decodeURL(s string) int{
    var num int = 0

    for _, element := range strings.Split(s, "") {
        num = num * int(BASE) + strings.Index(ALPHABET, element)
    }

    return num
}
