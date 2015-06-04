package main

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "net/http"
)

func ComputeHmac256(message string, secret string) string {
    key := []byte(secret)
    h := hmac.New(sha256.New, key)
    h.Write([]byte(message))
    return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func authenticateRequest(r *http.Request, value string) bool {
    expectedHash := ComputeHmac256(value, config.AuthenticationSecret)
    recievedHash := r.Header.Get("X-Auth-Signature")
    return expectedHash == recievedHash
}
