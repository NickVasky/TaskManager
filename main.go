package main

import (
	"TaskManager/api"
)

func main() {
	api.Serve()
	/*

		conn := repo.NewConnection()
		upr := conn.NewUserRepo()

		defer conn.Close()

		user := &repo.UserEntity{
			Username: "YaPerdolev",
			Password: "qwerty123",
			FirstName: sql.NullString{
				String: "Yaroslav",
				Valid:  true},
			SecondName: sql.NullString{
				String: "Perdolev",
				Valid:  true}}
		err0 := upr.Create(user)
		fmt.Printf("Err0: %v", err0)

		user1, err1 := upr.GetById(9)
		fmt.Printf("GetById - user1: %v err1: %v\n", user1, err1)

		user1.SecondName = sql.NullString{String: "Arsibekova", Valid: true}
		upr.Edit(user1)
		//user2, err2 := upr.GetByUsername("JOHNson")
		//fmt.Printf("GetByUserName - user2: %v err2: %v\n", user2, err2)

		//user2.Password = "SECURITY_PASSSSSSS"
		//upr.Edit(user2)

		t := &repo.TaskEntity{
			UserId:    user1.Id,
			Title:     "Get Beer",
			Goal:      "Get Craft beer in gose style",
			Measure:   "2 Bottles",
			Relevance: "I want to feel happines",
			Deadline: sql.NullTime{
				Time:  time.Now().Add(time.Hour * 24 * 4),
				Valid: true}}

		trp := conn.NewTaskRepo()

		trp.Create(t)

		t.Title = "Get Some Beer"
		trp.Edit(t)

		//upr.Delete(user2)
	*/
}
