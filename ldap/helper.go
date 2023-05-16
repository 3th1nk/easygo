package ldap

import "fmt"

func formatGUID(guid string) string {
	if len(guid) < 10 {
		return guid
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", guid[0:4], guid[4:6], guid[6:8], guid[8:10], guid[10:])
}
