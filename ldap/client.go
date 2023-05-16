package ldap

import (
	"crypto/tls"
	"fmt"
	"github.com/3th1nk/easygo/util/logs"
	"github.com/go-ldap/ldap"
	"time"
)

var (
	logger = logs.Default
)

type Config struct {
	Host        string       `json:"host,omitempty"`
	Port        int          `json:"port,omitempty"`
	Username    string       `json:"username,omitempty"`
	Password    string       `json:"password,omitempty"`
	BaseDN      string       `json:"base_dn,omitempty"`
	Timeout     int          `json:"timeout,omitempty"` // 连接超时，单位：秒
	EnableTLS   bool         `json:"enable_tls,omitempty"`
	EntryParser IEntryParser `json:"-"`
}

type Client struct {
	conn *ldap.Conn
	cfg  *Config
}

func NewClient(cfg *Config) (*Client, error) {
	var conn *ldap.Conn
	var err error
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	if cfg.EnableTLS {
		conn, err = ldap.DialTLS("tcp", addr, &tls.Config{InsecureSkipVerify: true})
	} else {
		conn, err = ldap.Dial("tcp", addr)
	}
	if err != nil {
		return nil, &Error{Type: ErrRequestFail, Err: err}
	}

	if cfg.Timeout > 0 {
		conn.SetTimeout(time.Duration(cfg.Timeout) * time.Second)
	}

	if cfg.EntryParser == nil {
		cfg.EntryParser = DefaultEntryParser{}
	}

	if err = conn.Bind(cfg.Username, cfg.Password); err != nil {
		return nil, &Error{Type: ErrAuthFail}
	}

	return &Client{
		conn: conn,
		cfg:  cfg,
	}, nil
}

func (this *Client) Close() {
	this.conn.Close()
}

func (this *Client) Get(baseDN string, filter ...string) (*ldap.Entry, error) {
	if len(filter) == 0 {
		filter = append(filter, "(objectClass=*)")
	}
	req := ldap.NewSearchRequest(
		baseDN, ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false, filter[0], nil, nil)
	sr, err := this.conn.Search(req)
	if err != nil {
		return nil, &Error{Type: ErrRequestFail, Err: err}
	}
	if len(sr.Entries) > 0 {
		return sr.Entries[0], nil
	}
	return nil, &Error{Type: ErrObjectNotExist}
}

func (this *Client) FindSingleLevel(baseDN string, filter ...string) ([]*ldap.Entry, error) {
	if len(filter) == 0 {
		filter = append(filter, "(objectClass=*)")
	}
	req := ldap.NewSearchRequest(
		baseDN, ldap.ScopeSingleLevel, ldap.NeverDerefAliases, 0, 0, false, filter[0], nil, nil)
	sr, err := this.conn.SearchWithPaging(req, 1000)
	if err != nil {
		return nil, &Error{Type: ErrRequestFail, Err: err}
	}
	return sr.Entries, nil
}

func (this *Client) FindSubtree(baseDN string, filter ...string) ([]*ldap.Entry, error) {
	if len(filter) == 0 {
		filter = append(filter, "(objectClass=*)")
	}
	req := ldap.NewSearchRequest(
		baseDN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, filter[0], nil, nil)
	sr, err := this.conn.SearchWithPaging(req, 1000)
	if err != nil {
		return nil, &Error{Type: ErrRequestFail, Err: err}
	}
	return sr.Entries, nil
}

func (this *Client) Check(dn ...string) error {
	if len(dn) == 0 {
		dn = append(dn, this.cfg.BaseDN)
	}
	if _, err := this.Get(dn[0]); err != nil {
		return err
	}
	return nil
}

// UserAuth 使用指定的账户、密码鉴权，并返回对应的人员信息
func (this *Client) UserAuth(account, password string) (*User, error) {
	filter := fmt.Sprintf("(|(uid=%v)(sAMAccountName=%v))", account, account)
	entries, err := this.FindSubtree(this.cfg.BaseDN, filter)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, &Error{Type: ErrAccountNotExist}
	}
	user := this.cfg.EntryParser.ToUser(entries[0], "")
	if user == nil {
		return nil, &Error{Type: ErrAccountNotExist}
	}

	// 验证密码
	if err = this.conn.Bind(entries[0].DN, password); err != nil {
		return nil, &Error{Type: ErrAuthFail}
	}
	// 重新绑定回建立连接的账户
	if err = this.conn.Bind(this.cfg.Username, this.cfg.Password); err != nil {
		return nil, &Error{Type: ErrAuthFail}
	}

	return user, nil
}

func (this *Client) recursion(baseDN, baseUID string, orgMap map[string]*OrgUnit, userMap map[string]*User) error {
	entries, err := this.FindSingleLevel(baseDN, this.cfg.EntryParser.EntrySearchFilter())
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if org := this.cfg.EntryParser.ToOrgUnit(entry, baseUID); org != nil {
			orgMap[org.UID] = org
			if err = this.recursion(entry.DN, org.UID, orgMap, userMap); err != nil {
				return err
			}
			continue
		}

		if u := this.cfg.EntryParser.ToUser(entry, baseUID); u != nil {
			userMap[u.UID] = u
			continue
		}

		logger.Debug("对象(%v)未解析", entry.DN)
	}

	return nil
}

func (this *Client) GetSubTreeOrganization(baseDN string) (map[string]*OrgUnit, map[string]*User, error) {
	var orgMap = make(map[string]*OrgUnit)
	var userMap = make(map[string]*User)
	// 考虑到整体的数据量可能很大，递归遍历
	if err := this.recursion(baseDN, "", orgMap, userMap); err != nil {
		return nil, nil, err
	}
	return orgMap, userMap, nil
}

func (this *Client) GetSingleLevelOrganization(baseDN string) (map[string]*OrgUnit, map[string]*User, error) {
	entries, err := this.FindSingleLevel(baseDN, this.cfg.EntryParser.EntrySearchFilter())
	if err != nil {
		return nil, nil, err
	}

	var orgMap = make(map[string]*OrgUnit)
	var userMap = make(map[string]*User)
	for _, entry := range entries {
		if org := this.cfg.EntryParser.ToOrgUnit(entry, ""); org != nil {
			orgMap[org.UID] = org
			continue
		}

		if u := this.cfg.EntryParser.ToUser(entry, ""); u != nil {
			userMap[u.UID] = u
			continue
		}

		logger.Debug("对象(%v)未解析", entry.DN)
	}
	return orgMap, userMap, nil
}
