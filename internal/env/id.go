package env

import "github.com/rs/xid"

func GetID() string {
	guid := xid.New()
	return guid.String()
}
