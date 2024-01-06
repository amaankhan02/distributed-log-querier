# Distributed Log Querier

This program allows you to query distributed log files on multiple machines, from any one of those machines that are running the program. From any running machine, you can run a grep query that runs on all the log files across all machines and prints output to your terminal (with the appropriate line counts, i.e., number of matching lines, and file names to designate where each log entry line came from). The distributed log querier also caches outputs of grep queries so that repeated grep queries are executed much faster.

## Build Instruction
* Install Go with at least version `1.19`
* Once installed, from the root directory of this project, run `go build cmd/main.go`
which will create an executable `./main` on Linux/Unix systems, or `main.exe` on Windows

## Instructions to run `./main`
* `./main` takes the following arguments
  * `-n` (number of total machines in network to use)
    * **type**: int
    * **default value**: 10
    * **usage**: Number of VMs you are going to run the program with. In the range of [2, 10]
    * **_NOTE_**: if you set n = 6, for instance, then you must boot up VMs 1...6. You cannot 
    use any VM larger than n to create your n VMs as this will not work since it automatically
    figures out your hostname and the other peer machine host names based off n
  * `-f` (log file name: _REQUIRED_)
    * **type**: string
    * **default** value: ""
    * **usage**: The filename of the log file
  * `-c` (cache size: _OPTIONAL_)
    * **type**: int
    * **default value**: 10
    * **usage**: Size of the in-memory cache on this VM
  * `-v` (verbose: _OPTIONAL_)
    * **type**: bool
    * **default value**: `false`
    * **usage**: Indicates if you want extra messages to be printed out aside
    from the outputs
  * `-t` (test directory: _OPTIONAL_)
    * **type**: string
    * **default value**: "" (NOT required field)
    * **usage**: Input a directory if you wish to boot up in TEST mode. 
    If the value is not "", then it boots up program in testing mode, and
    will store the Grep Outputs into JSON files under the directory you provide.
    Currently, the directory MUST already exist - it won't create one for you. 
    In future improvement we will add support for creating a new directory.
