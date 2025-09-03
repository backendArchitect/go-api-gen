package petstore




type Pets []Pet




type Pet struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	Tag string `json:"tag"`
}

