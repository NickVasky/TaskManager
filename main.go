package main

import (
	"TaskManager/dbactions"
)

func main() {
	user := dbactions.UserEntity{Username: "nickvasky", Password: "qwerty123", Name: "Никита", Surname: "Никита"}
	dbactions.CreateUser(user)
	user = dbactions.UserEntity{Username: "zloilegin", Password: "qwerty123", Name: "Олег", Surname: "Ильин"}
	dbactions.CreateUser(user)
	dbactions.ChangePassword(7, "qwerty123", "HelloBoy13")
}
