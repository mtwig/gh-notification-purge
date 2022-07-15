package app

type Config struct {
	PrintSubjects bool
	DryRun        bool
}

type application struct {
	printSubjects bool
	dryRun        bool
}

func CreateNew(config Config) (appRef *application, err error) {
	var app application
	app.printSubjects = config.PrintSubjects
	app.dryRun = config.DryRun
	appRef = &app
	return
}
