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

func main() {

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
