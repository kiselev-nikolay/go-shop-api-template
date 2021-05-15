package detailres

type Detail struct {
	Value string `json:"detail"`
}

func New(detail string) *Detail {
	return &Detail{detail}
}
