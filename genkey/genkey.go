package genkey

import (
	"crypto/sha256"
	"github.com/tsuka611/golang_sandbox/config"
	"github.com/tsuka611/golang_sandbox/log"
	"os"
	"time"
	"fmt"
)

func Gen() config.AppKey {
	host, err := os.Hostname()
	if err != nil {
		log.WARN.Println("Cannot get hostname.")
		host = "UNKNOWN HOST"
	}
	now := time.Now().String()
	key := fmt.Sprintf("%x", sha256.Sum256([]byte(host + now)))
	return config.AppKey(key)
}
