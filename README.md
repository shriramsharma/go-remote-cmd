#### go-remote-cmd: Run any unix command on multiple remote servers. The output will be printed on the host STDOUT.

##### Usage:
* ssh-add /path/to/your/private/key
* go install go-remote-cmd
* go-remote-cmd /path/to/file/containing/ips "command"

Example: `go-remote-cmd ips "tail -f access_logs.log"
