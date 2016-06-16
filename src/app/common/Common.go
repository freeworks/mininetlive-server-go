package common

import ( 
	logger "app/logger"
)

func CheckErr(err error, msg string) {
	if err != nil {
		logger.Error(msg,err)
		// log.Fatalln(msg, err)
	}
}



func GeneraToken() string {
	return "genToken"
}