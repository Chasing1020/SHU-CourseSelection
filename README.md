# SHU-CourseSelection

## 1. Overview

An auto-course-selection program written in Golang.

## 2. Quick Start

```sh
git clone https://github.com/Chasing1020/SHU-CourseSelection.git
cd SHU-CourseSelection
go mod tidy

# *** Then modify the configuration file: `config.yaml` ***

go run .
```

And if you want to use a cronjob, you should use `crontab -e` command and add
```text
# take Fri Jan 7 20:30:00 CST 2022 as an example
30 20 7 1 * go run /[PATH_TO_DIR]/main.go 2>&1 >> ~/selection.log
```

## 3. Bugs

Have found any bugs or suggestions? 
Please visit the [issue tracker](https://github.com/Chasing1020/SHU-CourseSelection/issues).

I'm glad if you have any feedback or give pull requests to this project.

## 4. Disclaimer

This project is only available for free academic discussions.

## 5. License

Licensed under the [Apache License](https://www.apache.org/licenses/LICENSE-2.0), Version 2.0 (the "License");
