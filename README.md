#Simple message sender

##Implementation requiremens
[./docs/challenge.pdf](./docs/challenge.pdf)

Simple message sender implementation with just standard Go library. Neither custom logs nor test framework are used in accordance to implementation requirements. To be consistent neither any task manager nor docker cluster manager are used

Dependencies are to be managed with dep [https://golang.github.io/dep/](https://golang.github.io/dep/)

To be honest there's no 3rd-party dependencies yet, but any project mush have an ability to import some vendor code, I suppose.

Tested with standard go test [https://golang.org/pkg/testing/](https://golang.org/pkg/testing/).

##Build

- manually: `dep ensure && go build -o build/messagesender`
- by script: `./build.sh [docker]` (use `docker` option to build with docker)

##Test

- manually: `go test ./schedule -race && go test ./transport -race` (only dependency packages are covered by unit tests; main package is not; it seems really unnecessary and even meaningless to store any complicated logic which would need unit tests coverage in main application package)
- by script: `./test.sh [<package_name>]` (use package_name to test specified package. If not specified, all packages will be tested)

##Run

- manually: 

`./build/commservice &` (see commservice PID in stdout)

`./build/messagesender [--url=<commservice_url>] [--schedule-delimiter=<schedule_delimiter>] [--csv-delimiter=<csv_delimiter>] [--file=<path_to_csv_file>]` (please wait until messagesender finishes)

Some options are acceptable; by default `messagesender` works with `http://localhost:9090/messages` url and `docs/customers.csv` file

`sudo kill -2 <commservice_PID>` (see commservice report in stdout)


- by script: `./run.sh [docker]` (use `docker` option to build with docker).

Being runned with local FS the script finishes and reports automatically. See commservice report in stdout.

If run with docker please stop commservice container manually, for example:

(use another terminal)

`docker exec commservice ps -aux` (see commservice PID)

`docker exec commservice kill -2 <commservice_PID>`
(see commservice report in main terminal)
