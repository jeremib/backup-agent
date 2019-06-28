package mysql

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/JamesStewy/go-mysqldump"
)

func Dump(host string, user string, pass string, db string, port int) (string, error) {
	// Open connection to database
	dumpDir := os.TempDir()                                     // you should create this directory
	dumpFilenameFormat := fmt.Sprintf("%s-20060102T150405", db) // accepts time layout string and add .sql at the end of file

	dba, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, pass, host, port, db))
	if err != nil {
		fmt.Println("Error opening database: ", err)
		return "Error opening database", err
	}

	// Register database with mysqldump
	dumper, err := mysqldump.Register(dba, dumpDir, dumpFilenameFormat)
	if err != nil {
		fmt.Println("Error registering databse:", err)
		return "Error registering database", err
	}

	// Dump database to file
	resultFilename, err := dumper.Dump()
	if err != nil {
		fmt.Println("Error dumping:", err)
		return "Error dumping database", err
	}
	defer dumper.Close()

	return resultFilename, nil
}
