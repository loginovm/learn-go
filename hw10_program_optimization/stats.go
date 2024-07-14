package hw10programoptimization

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := count(r, domain)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func count(r io.Reader, domain string) (result DomainStat, err error) {
	result = make(DomainStat)
	domLevel1 := "." + domain
	b := bufio.NewReader(r)
	var er error
	var line []byte
	for er == nil {
		line, er = b.ReadBytes('\n')
		if er != nil && er != io.EOF {
			return result, er
		}
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var user User
		if err = jsoniter.Unmarshal(line, &user); err != nil {
			return
		}
		if len(user.Email) > 0 {
			matched := domain == "" || strings.HasSuffix(user.Email, domLevel1)
			if matched {
				i := strings.IndexByte(user.Email, '@')
				domLevel2 := strings.ToLower(user.Email[i+1:])
				result[domLevel2]++
			}
		}
	}
	return
}
