package main

import (
	"fmt"
	"github.com/LucifaDva/go-mysql-struct/db"
	"github.com/LucifaDva/go-mysql-struct/models"
)

func main() {
	dbPool := new(db.DBPool)
	err := dbPool.Init("db.json", nil, 0, 0, false)
	if err != nil {
		fmt.Println("init dbpool err: ", err)
	}
	dbmodels := new(models.DBModel)
	dbmodels.Init(dbPool)

	persons := make([]models.Person, 0, 100)
	err = dbmodels.GetAllPersons(&persons)
	fmt.Println("persons:", persons)
	fmt.Println("err:", err)
}