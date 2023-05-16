package ldap

import (
	"testing"
)

func TestClient(t *testing.T) {
	cfg := Config{
		Host:        "192.168.1.119",
		Port:        389,
		Username:    "admin",
		Password:    "admin@123",
		BaseDN:      "DC=demo,DC=com",
		EnableTLS:   false,
		Timeout:     5,
		EntryParser: DefaultEntryParser{},
	}
	cli, err := NewClient(&cfg)
	if err != nil {
		t.Error(err)
		return
	}

	if _, err = cli.UserAuth("user1", "user1@123"); err != nil {
		t.Error("鉴权异常：", err)
	}

	orgUnits, users, err := cli.GetSubTreeOrganization(cfg.BaseDN)
	if err != nil {
		t.Error(err)
		return
	}

	var dnToUser = make(map[string]*User)
	for _, u := range users {
		dnToUser[u.DN] = u
	}

	t.Logf("====部门列表(%d)====", len(orgUnits))
	for _, u := range orgUnits {
		t.Log("uid:", u.UID)
		t.Log("\tparent org:", u.ParentOrgUnit)
		t.Log("\tDN:", u.DN)
		t.Log("\tname:", u.Name)
		if u.Manager != "" {
			if m, ok := dnToUser[u.Manager]; ok {
				t.Log("\tmanager:", m.Name)
			} else {
				manager, err := cli.Get(u.Manager)
				if err != nil {
					t.Log("\tmanager:", u.Manager)
				} else {
					t.Log("\tmanager:", cfg.EntryParser.UUID(manager))
				}
			}
		}
	}

	t.Logf("====人员列表(%d)====", len(users))
	for _, u := range users {
		t.Log("uid:", u.UID)
		t.Log("\torg:", u.OrgUnit)
		t.Log("\taccount:", u.Account)
		t.Log("\tname:", u.Name)
		t.Log("\tmobile:", u.Mobile)
		t.Log("\ttelephone:", u.Telephone)
		t.Log("\temail:", u.Email)
	}
}
