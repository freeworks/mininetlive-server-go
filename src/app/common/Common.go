package common

import ( "log" )

func CheckErr(err error, msg string) {
	if err != nil {
		log.Println(msg, err)
		// log.Fatalln(msg, err)
	}
}



func GeneraToken() string {
	return "genToken"
}