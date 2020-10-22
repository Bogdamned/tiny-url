package main

import (
	"errors"
	"log"
	"math"
	"net/url"
	"regexp"
	"strings"
)

const base = 62
const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Encode returns a base62 representation as
// string of the given integer number.
func Encode(x uint64) string {
	b := []byte{}

	if x == 0 {
		return charSet[:1]
	}

	for x > 0 {
		r := x % base
		x /= base
		b = append([]byte{charSet[r]}, b...)
	}

	return string(b)
}

// Decode returns uint64 representation of a string using base62 conversion
func Decode(s string) (uint64, error) {
	var result uint64 = 0
	pow := len(s) - 1

	for _, v := range s {
		pos := strings.Index(charSet, string(v))
		if pos == -1 {
			return 0, errors.New("Base62 decode: invalid string")
		}
		result += uint64(pos) * uint64(math.Pow(float64(base), float64(pow)))
		pow--
	}

	return result, nil
}

func urlIsValid(u string) bool {
	isValid := true

	if match, _ := regexp.MatchString("^https?://", u); !match {
		u = "http://" + u
	}

	_, err := url.ParseRequestURI(u)
	if err != nil {
		isValid = false
	}

	return isValid
}

func checkError(message string, err error) {
	if err != nil {
		log.Println(message, err)
	}
}
