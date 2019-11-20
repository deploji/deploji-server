package dto

type Version struct {
	Name    string
	Value   string
	SortKey string `json:"-"`
}

type ByName []Version

func (s ByName) Len() int {
	return len(s)
}
func (s ByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByName) Less(i, j int) bool {
	return s[i].SortKey > s[j].SortKey
}
