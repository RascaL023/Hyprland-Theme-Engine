package helper

func CheckErr(err error, msg string) {
	if err != nil {
		panic(msg);
	}
}
