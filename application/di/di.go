package di

import (
	"darbelis.eu/stabas/dao"
	"darbelis.eu/stabas/db"
	"darbelis.eu/stabas/entities"
	"darbelis.eu/stabas/my_tests"
)

var liteDatabase *db.Database = nil

func GetLiteDatabase() *db.Database {
	if liteDatabase == nil {
		// TODO load path from env or a config param
		liteDatabase = db.NewDatabase("./missions/database.db")
	}

	return liteDatabase
}

func NewTaskRepository(environment string) dao.ITasksRepository {
	if environment == "dev" {
		return my_tests.NewTasksRepository()
	}

	if environment == "empty" {
		return dao.NewTasksRepository([]*entities.Task{}, 1)
	}

	if environment == "prod" {
		repo, err := dao.NewLiteTaskRepository(GetLiteDatabase())
		if err != nil {
			panic(err)
		}
		return repo
	}

	panic("wrong config for tasks repository creation " + environment)
}

func NewParticipantsRepository(environment string) dao.IParticipantsRepository {
	if environment == "dev" {
		return my_tests.NewParticipantsRepository()
	}
	if environment == "empty" {
		return dao.NewParticipantsRepository([]*entities.Participant{})
	}

	if environment == "prod" {
		repo, err := dao.NewLiteParticipantsRepository(GetLiteDatabase())
		if err != nil {
			panic(err)
		}
		return repo
	}

	panic("wrong config for participants repository creation " + environment)
}
