# wfs

wfs is a system for automating the execution of your dev tools based on file system events.

## Overview

* Supports file creation, change and deletion events
* Runs commands in any shell
* Provides useful environment variables to the shell
* Actions are defined on directory, file name pattern (glob) and event type.

## Getting Started

### Installation

```
go get github.com/JamesOwenHall/wfs
```

### Your First Config File

wfs uses a configuration file to define which actions to perform on various file system events.  The config file is written in [YAML](https://en.wikipedia.org/?title=YAML), which offers great readability.

Start off by creating a file named `wfs.yml` with the following contents:
```YAML
shell: bash
files: 
  - path: .
    name: "*.txt"
    create: echo Created a file $filename >> out.log
    change: echo Changed a file $filename >> out.log
    delete: echo Deleted a file $filename >> out.log
```

Let's break this down.  First you define `shell: bash` which tells wfs that it should use Bash to execute the commands.

Next you define a list called `files`.  Each entry in this list defines a group of files.  In our case, we only have one group.  This group monitors all files in `.` that have a file name ending in `.txt`.

Finally, we define the commands to execute whenever such a file is created, changed or deleted.

### Environment Variables

Whenever wfs needs to run a command, it first sets several environment variables which can be used in the command.  These variables are:

* `path` is the absolute path of the target file.  (Note: it is in lowercase so that it doesn't conflict with `PATH`)
* `dir` is the absolute path of the target file's parent directory.
* `dirname` is the parent directory's name without its full path.
* `filename` is the target file's name without its full path.
* `fileradical` is the target file's name without its file extension.
* `fileext` is the target file's extension (including the .)

## License

MIT Licensed.  See `LICENSE` file for details.