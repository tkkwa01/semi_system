package mysql

import (
	"fmt"
	"semi_systems/driver"
	"semi_systems/keijiban/domain"
)

func init() {
	err := driver.GetRDB().AutoMigrate(
		&domain.User{},
		&domain.Article{},
	)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
