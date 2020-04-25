package build

type IndexPages struct {
	Posts    ByDate
	PageSize int
	IDFunc   func(i int) string
}

func (ip IndexPages) Len() int {
	quotient := len(ip.Posts) / ip.PageSize
	remainder := len(ip.Posts) % ip.PageSize
	if remainder != 0 {
		return quotient + 1
	}
	return quotient
}

func (ip IndexPages) At(i int) (string, interface{}) {
	start := i * ip.PageSize
	end := start + ip.PageSize
	if len(ip.Posts) < end {
		end = len(ip.Posts)
	}
	return ip.IDFunc(i), ip.Posts[start:end]
}
