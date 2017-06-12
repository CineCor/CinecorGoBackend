package models

type UploadFormat struct {
	Cinemas    []Cinema
	LastUpdate string
}

func (self *UploadFormat) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["last_update"] = self.LastUpdate
	cinemas := make([]map[string]interface{}, len(self.Cinemas))
	for i, v := range self.Cinemas {
		cinemas[i] = v.ToMap()
	}
	m["cinemas"] = cinemas
	return m
}
