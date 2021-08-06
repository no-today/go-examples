package main

const (
	ErrorNotFound = MapError("the key does not exist")
)

// MapError Map 错误, 实现 Error 接口
type MapError string

func (e MapError) Error() string {
	return string(e)
}

type Map struct {
	m map[string]string
}

// Get 根据 Key 获取值
func (m *Map) Get(key string) (string, error) {
	result, exists := m.m[key]
	if !exists {
		return "", ErrorNotFound
	}

	return result, nil
}

// Put 新增键值对
func (m *Map) Put(key, val string) {
	// 延迟初始化时机
	if m.m == nil {
		m.m = make(map[string]string)
	}

	m.m[key] = val
}

// Remove 根据 Key 删除条目
func (m *Map) Remove(key string) {
	delete(m.m, key)
}
