package main

import (
	"fmt"
	"github.com/mcuadros/go-defaults"
)

type UserInformation struct {
	userName           string `default:"root"`
	passWord           string `default:"123456"`
	sqlName            string `default:"mysql"`
	ipName             string `default:"127.0.0.1"`
	portName           string `default:"3306"`
	charsetInformation string `default:"utf8mb4"`
	parseTime          string `default:"True"`
	locInformation     string `default:"Local"`
}

func NewExampleBasic() *UserInformation {
	example := new(UserInformation)
	defaults.SetDefaults(example) //<-- This set the defaults values
	return example
}

//dsnInfor := "root:123456@tcp(127.0.0.1:3306)/mysql?charset=utf8mb4&parseTime=True&loc=Local"

func main() {
	targetInformation := NewExampleBasic()
	fmt.Println(targetInformation)
	fmt.Println(targetInformation.userName)
	dsnInformation := targetInformation.userName + ":" + targetInformation.passWord + "@tcp(" + targetInformation.ipName + ":" + targetInformation.portName + "0/" + targetInformation.sqlName
	fmt.Println(dsnInformation)
}
