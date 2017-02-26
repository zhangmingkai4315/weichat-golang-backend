package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

//status means the query running status.
type Status int32

//those status will act like a enum type
const (
	QUERY_INIT Status = iota
	QUERY_READY
	QUERY_RUNNING
	QUERY_STOPPED
	QUERY_ERROR
	QUERY_SKIP
)
const MAXCONCURRENT = 5

var re_date *regexp.Regexp

func init() {
	re_date = regexp.MustCompile("([0-9]+) (hours|hour|minute|minutes|days|day) ago")
}

func GetDateForThreadPostLayout(s string) (t time.Time, err error) {
	format := "January 2, 2006 , 3:04 pm"
	t, err = time.Parse(format, s)
	if err != nil {
		fmt.Println(err)
		return
	}
	return t, nil
}

func GetDateFromMMDDYYYY(s string) (t time.Time, err error) {
	const shortForm = "1/02/2006"
	t, err = time.Parse(shortForm, s)
	if err != nil {
		fmt.Println(err)
		return
	}
	return t, nil
}

func GetDateForInfoSec(s string) (t time.Time, err error) {
	format := "02 Jan 2006"
	t, err = time.Parse(format, s)
	if err != nil {
		fmt.Println(err)
		return
	}
	return t, nil
}

// GetDateFromString 转换相对时间到绝对时间
// Input :"1 hour ago" or "52 minutes ago" , "2 days ago"
func GetDateFromString(s string) (t time.Time, err error) {
	ret := re_date.FindStringSubmatch(s)
	if len(ret) == 0 || len(ret) != 3 {
		err = errors.New("Can't get anything")
		return
	}
	value := ret[1]
	unit := []byte(ret[2])[0]
	if unit == 'd' {
		i, aerr := strconv.Atoi(value)
		if aerr != nil {
			return
		}
		value = fmt.Sprintf("%dh", i*24)
	} else {
		value = fmt.Sprintf("%s%c", value, unit)
	}

	d, err := time.ParseDuration(value)
	if err != nil {
		return
	}
	return time.Now().Add(-d), nil
}

func GetMD5Hash(list ...string) string {
	var buffer bytes.Buffer
	for _, t := range list {
		buffer.WriteString(t)
	}
	hasher := md5.New()
	hasher.Write(buffer.Bytes())
	return hex.EncodeToString(hasher.Sum(nil))
}
