package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jeremib/backup-agent/mysql"
	"github.com/jeremib/backup-agent/remote/aws"
)

type Mysql struct {
	Host string `json:"host"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Db   string `json:"db"`
	Port int    `json:"port"`
}
type Source struct {
	Mysql []Mysql `json:"mysql"`
}

type S3 struct {
	AccessKey string
	Secret    string
	Region    string
	Bucket    string
}

type Destination struct {
	S3 []S3 `json:"s3"`
}

type Config struct {
	Source      Source      `json:"sources"`
	Destination Destination `json:"destinations"`
}

// type Config struct {
// 	Sources struct {
// 		Mysql []struct {
// 			Host string
// 			User string
// 			Pass string
// 			Db   string
// 			Port int
// 		}
// 	}
// 	Destinations struct {
// 		S3 []struct {
// 			AccessKey string
// 			Secret    string
// 			Region    string
// 			Bucket    string
// 		}
// 	}
// }

func main() {

	// resultFilename, _ := mysql.Dump()

	// f, _ := os.Open(resultFilename)
	// reader := bufio.NewReader(f)
	// content, _ := ioutil.ReadAll(reader)

	// name := resultFilename + ".gz"

	// f, _ = os.Create(name)
	// w := gzip.NewWriter(f)
	// w.Write(content)
	// w.Close()

	// aws.Upload(name, name)
	// _ = os.Remove(resultFilename)
	// _ = os.Remove(name)

	// backup.Compress("./", os.TempDir()+"/files.zip")
	// aws.Upload(os.TempDir()+"/files.zip", "test/files.zip")

	configFile, err := os.Open("./config.json")

	if err != nil {
		fmt.Println(err)
	}

	defer configFile.Close()

	bv, _ := ioutil.ReadAll(configFile)

	var config Config

	json.Unmarshal(bv, &config)
	for _, file := range backupSources(config.Source) {
		sendFileToDestinations(file, config.Destination)
		os.Remove(file)
	}

}

func backupSources(Source Source) []string {
	var ret []string

	for _, mysqlConfig := range Source.Mysql {
		fmt.Printf("Backing up MySQL: %s:%s\n", mysqlConfig.Host, mysqlConfig.Db)
		file, _ := mysql.Dump(
			mysqlConfig.Host,
			mysqlConfig.User,
			mysqlConfig.Pass,
			mysqlConfig.Db,
			mysqlConfig.Port,
		)
		ret = append(ret, file)
	}
	return ret
}

func sendFileToDestinations(filename string, destination Destination) {
	for _, s3Config := range destination.S3 {
		fmt.Printf("Sennding %s to S3:%s\n", filename, s3Config.Bucket)
		aws.Upload(
			s3Config.AccessKey,
			s3Config.Secret,
			s3Config.Region,
			s3Config.Bucket,
			filename,
			path.Base(filename),
		)
	}
}
