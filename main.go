package main

import "fmt"

const (
	VERSION               = "0.01"
	ConfigurationFileName = "configuration.yaml"
	LogFileName           = "log.log"
)

var (
	configuration Configuration
	logger        *Logger
)

func main() {
	// Prepare logger
	logger = NewLogger(LogFileName)
	defer logger.Destroy()

	// Read configuration
	logger.Info("Read YAML configuration:")
	configuration = readConfigurationFile(ConfigurationFileName)
	logger.Info(" -- Account name: " + configuration.TelegraphAcountName)
	logger.Info(fmt.Sprintf(" -- Max image size (MB): %d", configuration.MaxImageSizeMb))
	logger.Info(fmt.Sprintf(" -- Debug mode: %t", configuration.Debug))

	// Get today folders
	currentDirList := readDirAndFilter(".", FilterOnlyFoldersWithContent())
	todayDirList := FilterFilesNames(
		FilterFilesList(currentDirList, FilterOnlyTodayFolders()))
	logger.Info("")
	logger.Info(fmt.Sprintf("Found %d today folders:", len(todayDirList)))
	for _, folder := range todayDirList {
		logger.Info(" -> " + folder)
	}

	// Generating tasks for each folder
	var tasks []Task
	for _, folder := range todayDirList {
		tasks = append(tasks, generateTaskFromFolder(folder))
	}

	// Connect to Telegraph
	client, _ := getTelegraphConnection(configuration.TelegraphAcountName)

	// Upload each folder (execute each task) and collecting urls
	for _, task := range tasks {
		executeTask(client, task)
	}

}
