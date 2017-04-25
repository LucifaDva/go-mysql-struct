package main
import(
"db"
"fmt"
"models"
)
func main() {
	dbPool := new(db.DBPool)
	err := dbPool.Init("db.json", nil, 0, 0, false)
	if err != nil {
		fmt.Println("init dbpool err: ", err)
	}
	dbmodels := new(models.DBModel)
	dbmodels.Init(dbPool)

	persons := make([]models.Person, 0 , 100)
	err = dbmodels.GetAllPersons(&persons)
	fmt.Println("persons:", persons)
	fmt.Println("err:", err)
/*	dbserver, err := dbPool.Get("center")

	type Shop struct {
		Id 	int64
		Name string
		Sex string
		Age int
	}
	result := make([]Shop, 0, 1)
	err = q(dbserver, &result)
	fmt.Println("query db err: ", err)
	fmt.Println(result)*/
}
/*
func q(dbserver db.DBServer,result interface{}) (err error) {
	
	if err != nil {
		fmt.Println("db pool err: ", err)
		return
	}
	query := "select * from test"
	err = dbserver.QueryHelper(result, query)
	return
}*/