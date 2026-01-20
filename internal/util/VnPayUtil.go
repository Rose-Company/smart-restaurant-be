package util

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

func HmacSHA512(key, data string) string {
	if key == "" || data == "" {
		return ""
	}
	h := hmac.New(sha512.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func GetIPAddress(r *http.Request) string {
	if r == nil {
		return "127.0.0.1" // Default IP for server-side requests
	}
	ip := r.Header.Get("X-FORWARDED-FOR")
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}

func GetRandomNumber(length int) string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	var sb strings.Builder

	for i := 0; i < length; i++ {
		sb.WriteByte(digits[rand.Intn(len(digits))])
	}
	return sb.String()
}

func GetPaymentURL(params map[string]string, encodeKey bool) string {
	keys := make([]string, 0)

	for k, v := range params {
		if v != "" {
			keys = append(keys, k)
		}
	}

	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		key := k
		if encodeKey {
			key = url.QueryEscape(k)
		}
		value := url.QueryEscape(params[k])
		parts = append(parts, key+"="+value)
	}

	return strings.Join(parts, "&")
}
