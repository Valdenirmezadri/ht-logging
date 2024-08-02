package main

import (
	htlogging "ht-logging"
	logging "ht-logging"
	"log"
	"os"
)

func main() {
	formatFile := logging.MustStringFormatter(
		`%{time:Jan 02 2006 15:04:05} %{shortfile} ▶ %{level:.4s} %{message}`,
	)

	formatConsole := htlogging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfile} ▶ %{level:.4s} %{color:reset}%{message}`,
	)
	console := htlogging.NewLogBackend(os.Stderr, "", 0)
	consoleFormatter := htlogging.NewBackendFormatter(console, formatConsole)

	hlog, err := htlogging.New("debug", consoleFormatter)
	if err != nil {
		log.Fatal(err)
	}

	hlog.SetLevel("DEBUG")

	hlog.Debugf("debug %s", "arg")
	hlog.Error("error")

	file, err := os.OpenFile("teste.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	writer := logging.NewLogBackend(file, "", 0)
	fileFromatter := logging.NewBackendFormatter(writer, formatFile)
	logfile, err := htlogging.New("debug", fileFromatter)
	if err != nil {
		log.Fatal(err)
	}

	logfile.SetLevel("DEBUG")

	logfile.Debugf("debug %s", "arg")
	logfile.Error("error")

}
