package payxgo_util

import (
	r "crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/url"
	"sort"

	"github.com/rs/xid"
)

func Xid() string {
	return xid.New().String()
}

func keySort(m map[string]interface{}) []string {
	if m == nil {
		return nil
	}
	var s []string
	for k := range m {
		s = append(s, k)
	}
	sort.Strings(s)
	return s
}

func dealParam(param map[string]interface{}) string {
	if len(param) == 0 {
		return ""
	}
	s := keySort(param)
	p := make(url.Values)

	for _, v := range s {
		if param[v] != nil {
			p.Add(v, fmt.Sprint(param[v]))
		}
	}
	return p.Encode()
}

func Sign(key string, param map[string]interface{}) string {
	return calcSha512(key + dealParam(param) + key)
}

func calcSha512(str string) string {
	h := sha512.New()
	_, err := io.WriteString(h, str)
	if err != nil {
		log.Println(NewError(1024, err.Error()).Error())
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}

func RsaEncrypt(data, key /*密钥*/ []byte) []byte {
	publicKey, err := x509.ParsePKCS1PublicKey(pBlock(key))
	if err != nil {
		log.Println(NewError(1021, err.Error()).Error())
		return nil
	}
	crypt, err := rsa.EncryptPKCS1v15(r.Reader, publicKey, data)
	if err != nil {
		log.Println(NewError(1022, err.Error()).Error())
		return nil
	}
	return crypt
}

func pBlock(key []byte) []byte {
	buf, err := base64.RawStdEncoding.DecodeString(string(key))
	if err != nil {
		log.Println(NewError(1023, err.Error()).Error())
		return nil
	}
	return buf
}
