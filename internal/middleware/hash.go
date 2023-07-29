package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
)

// hashWriter реализует интерфейс http.ResponseWriter
type hashWriter struct {
	http.ResponseWriter
	key []byte
}

func newHashWriter(w http.ResponseWriter, key []byte) *hashWriter {
	return &hashWriter{
		ResponseWriter: w,
		key:            key,
	}
}

func (h *hashWriter) Write(p []byte) (int, error) {
	if !bytes.Equal(h.key, nil) {
		sha := hmac.New(sha256.New, h.key)
		sha.Write(p)
		dst := hex.EncodeToString(sha.Sum(nil))
		h.Header().Add("HashSHA256", dst)
	}
	return h.ResponseWriter.Write(p)
}

func WithHash(key []byte, h http.HandlerFunc) http.HandlerFunc {
	hashFn := func(w http.ResponseWriter, r *http.Request) {
		ow := w
		if !bytes.Equal(key, nil) {
			sha, err := hex.DecodeString(r.Header.Get("HashSHA256"))
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if !bytes.Equal(sha, nil) {
				h := hmac.New(sha256.New, key)
				body, err := io.ReadAll(r.Body)

				if err == nil {
					h.Write(body)
					dst := h.Sum(nil)
					if !bytes.Equal(sha, dst) {
						fmt.Println(key)
						fmt.Println(sha)
						fmt.Println(dst)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					r.Body = io.NopCloser(bytes.NewBuffer(body))
				}
			}
			hw := newHashWriter(w, key)
			ow = hw
		}
		h.ServeHTTP(ow, r)
	}
	return http.HandlerFunc(hashFn)
}
