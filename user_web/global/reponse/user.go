package reponse

import (
	"fmt"
	"time"
)

type JsonTime time.Time

func (j JsonTime) MarshalJson() ([]byte, error) {
	var stmp = fmt.Sprintf("\"%s\"", time.Time(j).Format("2006-01-02 01:02:30"))
	return []byte(stmp), nil
}

type UserResponse struct {
	Id       int32  `json:"id"`
	NickName string `json:"name"`
	//Birthday time.Time `json:"birthday"`
	//Birthday JsonTime `json:"birthday"`
	//Birthday uint64 `json:"birthday"`
	Birthday string `json:"birthday"`
	Gender   string `json:"gender"`
	Mobile   string `json:"mobile"`
}
