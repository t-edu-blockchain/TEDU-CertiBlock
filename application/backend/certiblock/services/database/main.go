package database

import (
	"certiblock/base"
	"database/sql"
	"fmt"
	"os"
	"slices"

	"github.com/go-sql-driver/mysql"
)

func InitConnection(context *base.ApplicationContext) error {
	dsn, err := getMysqlDSN(context)
	if err != nil {
		return err
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	context.DB = db
	err = migrate(context)
	if err != nil {
		return err
	}
	return nil
}

func migrate(context *base.ApplicationContext) error {
	preludeSql := "CREATE TABLE IF NOT EXISTS `migrations` (`id` INT NOT NULL AUTO_INCREMENT, `name` VARCHAR(255) NOT NULL, PRIMARY KEY (`id`));"

	_, err := context.DB.Exec(preludeSql)
	if err != nil {
		return err
	}

	fmt.Println("Migrator: Migrating database...")

	migrationFiles, err := os.ReadDir("services/database/migrations")
	if err != nil {
		return err
	}

	migrationFileNames := []string{}

	for _, migrationFile := range migrationFiles {
		// skip directories, symlinks and other non-regular files
		if migrationFile.Type().IsRegular() {
			migrationFileNames = append(migrationFileNames, migrationFile.Name())
		}
	}

	slices.Sort(migrationFileNames)

	for _, migrationFileName := range migrationFileNames {
		fmt.Printf("Migrator: Applying migration %s... ", migrationFileName)
		sql := "SELECT COUNT(*) FROM `migrations` WHERE `name` = ?;"
		row := context.DB.QueryRow(sql, migrationFileName)

		var count int
		err := row.Scan(&count)
		if err != nil {
			return err
		}

		if count == 0 {
			// APPLY MIGRATION
			// read sql file
			sqlFilePath := "services/database/migrations/" + migrationFileName
			sqlFile, err := os.ReadFile(sqlFilePath)
			if err != nil {
				return err
			}

			// execute sql
			_, err = context.DB.Exec(string(sqlFile))
			if err != nil {
				return err
			}

			// record migration
			sql = "INSERT INTO `migrations` (`name`) VALUES (?);"
			_, err = context.DB.Exec(sql, migrationFileName)
			if err != nil {
				return err
			}

			fmt.Println("SUCCESSFULLY APPLIED.")
		} else {
			fmt.Println("ALREADY APPLIED.")
		}
	}

	fmt.Println("Migrator: Done.")
	return nil
}

func getMysqlDSN(context *base.ApplicationContext) (string, error) {
	cfg := mysql.Config{
		User:                 context.Config.Database.Username,
		Passwd:               context.Config.Database.Password,
		Net:                  "tcp",
		Addr:                 context.Config.Database.Address,
		DBName:               context.Config.Database.DatabaseName,
		AllowNativePasswords: true,
		AllowOldPasswords:    true,
		MultiStatements:      true, // for migrations!
	}

	return cfg.FormatDSN(), nil
}
