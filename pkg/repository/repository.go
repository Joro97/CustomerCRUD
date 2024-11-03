package repository

import (
	"CustomerCRUD/utils"
	"database/sql"
)

func GetDB(isLocalDb bool, connStrEnvVar string) (*sql.DB, error) {
	if isLocalDb {
		return utils.GetLocalDB()
	} else {
		return nil, nil // TODO:
	}
}
