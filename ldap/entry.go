package ldap

import (
	"github.com/3th1nk/easygo/util/strUtil"
	"github.com/go-ldap/ldap"
	"github.com/toolkits/slice"
	"strings"
)

type IEntryParser interface {
	UUID(entry *ldap.Entry) string
	EntrySearchFilter() string
	ToUser(entry *ldap.Entry, orgUnit string) *User
	ToOrgUnit(entry *ldap.Entry, parentOrgUnit string) *OrgUnit
}

type DefaultEntryParser struct{}

func (e DefaultEntryParser) UUID(entry *ldap.Entry) string {
	// 唯一标识字段字段
	// WindowsAD: objectGUID
	// OpenLDAP:  entryUUID
	for _, attr := range []string{"objectGUID", "entryUUID"} {
		if guid := entry.GetAttributeValue(attr); guid != "" {
			return formatGUID(guid)
		}
	}

	// 调试OpenLDAP时没有获取到entryUUID属性值，但从phpLdapAdmin页面中可以看到, 使用DN作为唯一标识兼容处理一下
	// 格式：OU=R&D,DC=demo,DC=com
	values := strUtil.Split(entry.DN, ",", true, func(s string) string {
		kv := strings.Split(s, "=")
		if len(kv) > 1 {
			return kv[1]
		}
		return ""
	})
	return strings.Join(values, ".")
}

func (e DefaultEntryParser) EntrySearchFilter() string {
	return "(|(objectClass=organizationalPerson)(objectClass=inetOrgPerson)(objectClass=organizationalUnit))"
}

func (e DefaultEntryParser) isUser(entry *ldap.Entry) bool {
	// 对象类型字段，Windows AD和OpenLDAP一样
	for _, class := range entry.GetAttributeValues("objectClass") {
		if slice.ContainsString([]string{"organizationalPerson", "inetOrgPerson"}, class) {
			return true
		}
	}
	return false
}

func (e DefaultEntryParser) isOrgUnit(entry *ldap.Entry) bool {
	return slice.ContainsString(entry.GetAttributeValues("objectClass"), "organizationalUnit")
}

func (e DefaultEntryParser) ToOrgUnit(entry *ldap.Entry, parentOrgUnit string) *OrgUnit {
	if !e.isOrgUnit(entry) {
		return nil
	}

	var name string
	for _, attr := range []string{"ou", "cn", "name", "displayName"} {
		if name = entry.GetAttributeValue(attr); name != "" {
			break
		}
	}
	return &OrgUnit{
		UID:           e.UUID(entry),
		DN:            entry.DN,
		Name:          name,
		Manager:       entry.GetAttributeValue("managedBy"),
		ParentOrgUnit: parentOrgUnit,
	}
}

func (e DefaultEntryParser) ToUser(entry *ldap.Entry, orgUnit string) *User {
	if !e.isUser(entry) {
		return nil
	}

	// 由于人员信息存在同义的字段，按照常用的优先级获取
	var email string
	for _, attr := range []string{"mail", "Email"} {
		email = entry.GetAttributeValue(attr)
		if email != "" {
			break
		}
	}

	var account string
	for _, attr := range []string{"sAMAccountName", "uid"} {
		account = entry.GetAttributeValue(attr)
		if account != "" {
			break
		}
	}
	if account == "" {
		// 以上取不到的话，尝试从邮箱中截取
		account = strings.Split(email, "@")[0]
		if account == "" {
			// 账户作为唯一标识，不能为空
			logger.Warn("对象(%v)字段(sAMAccountName|uid)为空", entry.DN)
			return nil
		}
	}

	var name string
	for _, attr := range []string{"displayName", "cn", "name"} {
		name = entry.GetAttributeValue(attr)
		if name != "" {
			break
		}
	}

	var mobile string
	for _, attr := range []string{"mobile"} {
		mobile = entry.GetAttributeValue(attr)
		if mobile != "" {
			break
		}
	}

	var telephone string
	for _, attr := range []string{"telephoneNumber", "homePhone", "Telephone"} {
		telephone = entry.GetAttributeValue(attr)
		if telephone != "" {
			break
		}
	}

	return &User{
		UID:       e.UUID(entry),
		OrgUnit:   orgUnit,
		Account:   account,
		Name:      name,
		Mobile:    mobile,
		Telephone: telephone,
		Email:     email,
	}
}
