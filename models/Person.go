package models

// struct to table
type Person struct {
	Id 		int64
	Name 	string
	Sex 	string
	Age 	int
}

func (m *DBModel) GetAllPersons(result interface{}) (err error) {
	dbServer, err := m.DbPool.Get("center")
	if err != nil {
		//fmt.Println("pool err :", err)
		return
	}
	query := "select * from test"
	err = dbServer.QueryHelper(result, query)
	return
}


