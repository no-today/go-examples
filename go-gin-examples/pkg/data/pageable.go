package data

import "strings"

const (
	// 防止一页请求太多数据拉垮数据库
	maxSize = 1000
)

// Pageable 调用各个数据库的API时使用, service 层、repository 层使用
type Pageable struct {
	Page int64
	Size int64

	Sorts  []Sort
	Fields []string

	// Response
	Total int64
}

// PageRequest 接口接收请求参数的结构, 只当作接口入参使用
type PageRequest struct {
	// Request
	Page int64 `form:"page"`
	Size int64 `form:"size"`

	// 用逗号分隔的字段列表。例如："foo,bar"。默认升序排列。应该（should）给字段添加后缀 " desc" 来表示降序。例如："foo desc,bar"。
	// 多余的空格可以忽略，"foo,bar desc"和 " foo , bar desc "是相等的。
	OrderBy string `form:"orderBy"`
}

// PageResponse 接口响应使用
type PageResponse struct {
	Page  int64 `json:"page,omitempty"`
	Size  int64 `json:"size,omitempty"`
	Total int64 `json:"total,omitempty"`

	// Data list
	List interface{} `json:"list,omitempty"`
}

type Sort struct {
	Field string
	Order Order
}

type Order int8

const (
	OrderASC  Order = 1
	OrderDESC       = -1
)

func (p *Pageable) Copy(total int64) *Pageable {
	return &Pageable{
		Page:   p.Page,
		Size:   p.Size,
		Sorts:  p.Sorts,
		Fields: p.Fields,
		Total:  total,
	}
}

func ToPageable(request *PageRequest) (pageable Pageable) {
	if request == nil {
		return
	}

	pageable.Page = request.Page
	pageable.Size = request.Size

	if pageable.Page == 0 {
		pageable.Page = 1
	}
	if pageable.Size == 0 || pageable.Size > maxSize {
		pageable.Size = 10
	}

	if request.OrderBy != "" {
		var sorts []Sort
		for _, sortStr := range strings.Split(request.OrderBy, ",") {
			sorts = append(sorts, parseSort(sortStr))
		}

		pageable.Sorts = sorts
	}

	return pageable
}

func parseSort(sortStr string) Sort {
	item := strings.Split(strings.TrimSpace(sortStr), " ")
	field := item[0]
	var order Order
	if len(item) == 1 {
		order = OrderASC
	} else {
		sort := item[len(item)-1]
		if sort == "desc" || sort == "DESC" {
			order = OrderDESC
		} else {
			order = OrderASC
		}
	}

	sort := Sort{
		Field: field,
		Order: order,
	}
	return sort
}

func ToPageResp(pageable *Pageable, list interface{}) *PageResponse {
	return &PageResponse{
		Page:  pageable.Page,
		Size:  pageable.Size,
		Total: pageable.Total,
		List:  list,
	}
}
