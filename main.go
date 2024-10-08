package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/ncruces/zenity"
)

const (
	taskFile       = "/storage/external/work-personal/Personal/obsidian/Personal Notes/Work/My Tasks.md"
	taskFileBackup = "/storage/external/work-personal/Personal/obsidian/Personal Notes/Work/My Tasks.old"
)

func main() {
	fileInfo, err := os.Stat(taskFile)
	if err != nil {
		log.Fatalf("'%s' does not exist", taskFile)
	}

	if !fileInfo.Mode().IsRegular() {
		log.Fatalf("'%s' is not a regular file", taskFile)
	}

	task, err := zenity.Entry("Task:", zenity.Title("Add a Task"), zenity.Modal())
	if err != nil {
		log.Fatal("Error getting task")
	}

	// make a backup
	err = copy(taskFile, taskFileBackup)
	if err != nil {
		log.Fatalf("Failed to make backup of '%s'", taskFile)
	}

	f, err := os.Open(taskFile)
	if err != nil {
		log.Fatalf("Failed to open '%s' for reading", taskFile)
	}

	buf, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("Could not read '%s': %v", taskFile, err)
	}

	if err := f.Close(); err != nil {
		log.Fatalf("Failed to close '%s' from reading: %v", taskFile, err)
	}

	contents := string(buf)

	lines := strings.Split(contents, "\n")

	f, err = os.OpenFile(taskFile, os.O_RDWR|os.O_TRUNC, fileInfo.Mode().Perm())
	if err != nil {
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

	d, err := os.OpenFile(taskFileBackup, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0555)
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
