package ldap

// 不同类型的LDAP服务端schema可能不一样，这里统一定义下组织单元和人员的结构，方便上层处理

// OrgUnit 组织单元
type OrgUnit struct {
	UID           string `json:"uid,omitempty"`             // 唯一标识
	DN            string `json:"dn,omitempty"`              // 组织单元DN
	ParentOrgUnit string `json:"parent_org_unit,omitempty"` // 父级组织单元的UID
	Name          string `json:"name,omitempty"`            // 组织单元名称
	Manager       string `json:"manager,omitempty"`         // 组织单元主管DN，由调用方自行转换成对应的人员UID
}

// User 组织人员
type User struct {
	UID       string `json:"uid,omitempty"`       // 唯一标识
	DN        string `json:"dn,omitempty"`        // 人员DN
	OrgUnit   string `json:"org_unit,omitempty"`  // 所属组织单元的UID
	Account   string `json:"account,omitempty"`   // 账户
	Name      string `json:"name,omitempty"`      // 姓名
	Mobile    string `json:"mobile,omitempty"`    // 手机号码
	Telephone string `json:"telephone,omitempty"` // 电话号码
	Email     string `json:"email,omitempty"`     // 邮箱
}
