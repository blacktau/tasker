package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/ncruces/zenity"
	"github.com/spf13/viper"
)

func main() {
	taskFile, taskFileBackup := loadConfig()

	fileInfo, err := os.Stat(taskFile)
	if err != nil {
		showError(fmt.Sprintf("'%s' does not exist", taskFile))
		log.Fatalf("'%s' does not exist", taskFile)
	}

	if !fileInfo.Mode().IsRegular() {
		showError(fmt.Sprintf("'%s' is not a regular file", taskFile))
		log.Fatalf("'%s' is not a regular file", taskFile)
	}

	task, err := zenity.Entry("Task:", zenity.Title("Add a Task"), zenity.Modal())
	if err != nil {
		showError("Error getting task")
		log.Fatal("Error getting task")
	}

	// make a backup
	err = copy(taskFile, taskFileBackup)
	if err != nil {
		showError(fmt.Sprintf("Failed to make backup of '%s': %v", taskFile, err))
		log.Fatalf("Failed to make backup of '%s': %v", taskFile, err)
	}

	f, err := os.Open(taskFile)
	if err != nil {
		showError(fmt.Sprintf("Failed to open '%s' for reading", taskFile))
		log.Fatalf("Failed to open '%s' for reading", taskFile)
	}

	buf, err := io.ReadAll(f)
	if err != nil {
		showError(fmt.Sprintf("Could not read '%s': %v", taskFile, err))
		log.Fatalf("Could not read '%s': %v", taskFile, err)
	}

	if err := f.Close(); err != nil {
		showError(fmt.Sprintf("Failed to close '%s' from reading: %v", taskFile, err))
		log.Fatalf("Failed to close '%s' from reading: %v", taskFile, err)
	}

	contents := string(buf)

	lines := strings.Split(contents, "\n")

	f, err = os.OpenFile(taskFile, os.O_RDWR|os.O_TRUNC, fileInfo.Mode().Perm())
	if err != nil {
		showError(fmt.Sprintf("Failed to open '%s' for writing: %v", taskFile, err))
		log.Fatalf("Failed to open '%s' for writing: %v", taskFile, err)
	}
	defer f.Close()

	inTodo := false
	inTodoTasks := false

	for _, line := range lines {

		if strings.Contains(line, "## To do") {
			inTodo = true
		}

		if inTodo && strings.Contains(line, "- [ ]") {
			inTodoTasks = true
		}

		if inTodo && inTodoTasks && strings.TrimSpace(line) == "" {
			fmt.Fprintf(f, "- [ ] %s\n", task)
			inTodo = false
			inTodoTasks = false
		}

		fmt.Fprintln(f, line)
	}

	f.Sync()
}

func copy(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to copy '%s' to '%s': %w", src, dst, err)
	}

	defer s.Close()

	d, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to copy '%s' to '%s': %w", src, dst, err)
	}

	defer d.Close()

	_, err = io.Copy(d, s)
	if err != nil {
		return fmt.Errorf("failed to copy '%s' to '%s': %w", src, dst, err)
	}

	err = d.Sync()
	if err != nil {
		return fmt.Errorf("failed to copy '%s' to '%s': %w", src, dst, err)
	}

	return nil
}

func showError(message string) {
	zenity.Error(message)
}

func loadConfig() (string, string) {
	viper.SetConfigName("tasker")
	viper.AddConfigPath("$XDG_CONFIG_HOME/tasker/")
	viper.AddConfigPath("$LOCALAPPDATA/tasker/")
	viper.AddConfigPath("$HOME/.tasker")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			createConfig()
		} else {
			showError(fmt.Sprint("Error reading config file: %v", err))
			log.Fatalf("Error reading config file: %v", err)
		}
	}

	taskFile := viper.GetString("taskfile")

	if len(strings.TrimSpace(taskFile)) == 0 {
		showError("Config file is missing taskfile location")
		log.Fatal("Config file is missing taskfile location")
	}

	taskFileBackup := viper.GetString("backupfile")

	if len(strings.TrimSpace(taskFileBackup)) == 0 {
		showError("Config file is missing taskfileBackup location")
		log.Fatal("Config file is missing taskfileBackup location")
	}

	return taskFile, taskFileBackup
}

func createConfig() {
	zenity.Info("tasker could not find its configuration file. A new configuration file will be created", zenity.Title("tasker configuration"), zenity.WarningIcon)
	taskFile, err := zenity.SelectFile(zenity.Title("Choose your task file"), zenity.FileFilters{
		{Name: "Markdown files", Patterns: []string{"*.md"}, CaseFold: false},
	})
	if err != nil {
		showError(fmt.Sprintf("Failed to select task file: %v", err))
		log.Fatalf("failed to choose task file: %v", err)
	}

	viper.Set("taskfile", taskFile)

	err = zenity.Question("Do you want to create a backup file when a task is added?")
	createBackup := true

	if err != nil {
		createBackup = false
	}

	viper.SetConfigType("yaml")
	err = viper.SafeWriteConfig()
	if err != nil {
		showError(fmt.Sprintf("Error saving configuration file: %v", err))
		log.Fatalf("failed to create configuration file: %v", err)
	}
}
