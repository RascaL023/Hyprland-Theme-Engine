package helper

import "fmt"

func CheckErr(err error, msg string) {
	if err != nil {
		fmt.Println(msg);
		panic(err);
	}
}
