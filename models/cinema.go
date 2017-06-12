package models

type Cinema struct {
	Id      int
	Name    string
	Address string
	Rooms   string
	Phone   string
	Web     string
	Movies  []Movie
}

func (self *Cinema) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["id"] = self.Id

	if len(self.Name) > 0 {
		m["name"] = self.Name
	}
	if len(self.Address) > 0 {
		m["address"] = self.Address
	}
	if len(self.Rooms) > 0 {
		m["rooms"] = self.Rooms
	}
	if len(self.Phone) > 0 {
		m["phone"] = self.Phone
	}
	if len(self.Web) > 0 {
		m["web"] = self.Web
	}
	movies := make([]map[string]interface{}, len(self.Movies))
	for i, v := range self.Movies {
		movies[i] = v.ToMap()
	}
	m["movies"] = movies
	return m

}
