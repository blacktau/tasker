# Tasker: A Simple Task Capturing UI

Tasker is a lightweight Go application that provides a graphical interface for
capturing tasks and appending them to a markdown-based task list. It supports
backups to ensure your task file is secure.

## Features

- Graphical Task Input: Utilizes a user-friendly UI to capture tasks.
- Markdown Integration: Appends tasks directly to a specified markdown file.
- Configurable File Paths: Define the task file and backup file locations via a configuration file.
- Automatic Backup: Creates a backup of the task file before updating it.

## Prerequisites

Tasker requires a supported platform for [zenity](https://github.com/ncruces/zenity) (used for UI dialogs).

## Installation

## Configuration

Tasker requires a configuration file containing the location of the task list markdown file
