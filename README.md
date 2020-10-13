# staticserver
Localhost server for HTML, JS, CSS development.

# How to use it

Type into terminal:

```
go get -u github.com/luckysuperduper/staticserver

staticserver -ssl -gzip -cache
```

You can also type only `staticserver` and answer the questions for configuring the local server.

# Livereload

For livereload use `air` and configure it to run `staticserver`

Example:

```
go get -u github.com/cosmtrek/air
```

Inside the folder with HTML, CSS and JS create `.air.toml` file with the following content:

```toml
# Config file for [Air](https://github.com/cosmtrek/air) in TOML format

# Working directory
# . or absolute path, please note that the directories following must be under root.
root = "."
tmp_dir = "tmp"

[build]
# Just plain old shell command. You could use `make` as well.
cmd = "staticserver -ssl -gzip -cache"
# Watch these filename extensions.
include_ext = ["go", "html", "js", "css"]
# Ignore these filename extensions or directories.
exclude_dir = ["node_modules", "tmp"]
# Watch these directories if you specified.
include_dir = []
# Exclude files.
exclude_file = []
# This log file places in your tmp_dir.
log = "air.log"
# It's not necessary to trigger build each time file changes if it's too frequent.
delay = 1000 # ms
# Stop running old binary when build errors occur.
stop_on_error = true
# Send Interrupt signal before killing process (windows does not support this feature)
send_interrupt = false
# Delay after sending Interrupt signal
kill_delay = 500 # ms

[log]
# Show log time
time = false

[color]
# Customize each part's color. If no color found, use the raw app log.
main = "magenta"
watcher = "cyan"
build = "yellow"
runner = "green"

[misc]
# Delete tmp directory on exit
clean_on_exit = true
```

# Problems?

Make sure `go/bin` is into Path.
## Example Linux and Mac:

Inside .bash or .zshrc file write:

```
export GOBIN=$HOME/go/bin
PATH=$PATH:$HOME/go/bin
```