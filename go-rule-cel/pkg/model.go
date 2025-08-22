package pkg

// 模拟请求结构（字段名大写）
type Request struct {
	Auth AuthContext `json:"auth"`
}

type AuthContext struct {
	Claims map[string]string `json:"claims"`
}

//
