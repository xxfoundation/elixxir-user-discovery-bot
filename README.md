# elixxir/user-discovery-bot

[![pipeline status](https://gitlab.com/elixxir/user-discovery-bot/badges/master/pipeline.svg)](https://gitlab.com/elixxir/user-discovery-bot/commits/master)
[![coverage report](https://gitlab.com/elixxir/user-discovery-bot/badges/master/coverage.svg)](https://gitlab.com/elixxir/user-discovery-bot/commits/master)

The user discovery bot helps users make first contact with other users. Users can search for other users using a string key (i.e. email address or phone number) and, if the user discovery bot finds a match for that user with the hash of the string, it will return a key ID. The user and the bot can then do a key exchange with the public key that the bot returns after the user queries that key ID to facilitate transfer of information that they need to talk to the user.

##Command-line options

|Long flag|Short flag|Effect|Example|
|---|---|---|---|
|--config|-c|Specify a different configuration file|--config udb2.yaml|
|--port|-p|Port which UDB will listen on|--port 1234|
|--logLevel|-l|Sets the log message level to print. (0 = info, 1 = debug, >1 = trace)|-v 2|
|--log||Path where log file will be saved|--log ./udb-logs/udb.log|
|--certPath||Path to UDB TLS certificate|--certPath ./keys/udb.pem|
|--keyPath||Path to UDB TLS private key|--keyPath ./keys/udb.key|
|--permCertPath||Path to permissioning public certificate|--permCertPath ./keys/permissioning.pem|
|--sessionPass||Password for UDB session files|--sessionPass pass|
|--devMode||Activate developer mode|--devMode|
|--help|-h|Shows a help message|-h|

## Example configuration

Note: Yaml prohibits the use of tabs. If you put tabs in your config file, the UDB will fail to parse it.

```yaml
# Path where UDB will store its logs
log: "udb.log"
# Path to NDF
ndfPath: "path/to/ndf"
# Path where UDB will store session file
sessionPath: "/path/to/session"
sessionPass: "password for session"

# Port which UDB will listen on
port: "1234"

# Database connection information
dbUsername: "cmix"
dbPassword: ""
dbName: "cmix_server"
dbAddress: ""

# Certificates
certPath: "/path/udb.pem"
keyPath: "/path/udb.key"
permCertPath: "permissioning.pem"

# Twilio account information
# Note: running with --devMode bypasses twilio verification
twilioSid: "sid"
twilioToken: "token"
twilioVerification: "verification"
```

## Running

### Installing dependencies
Running `make release` will update all dependent repositories to most recent release version

Running `go mod vendor` causes updated repos to push to your vendor folder, so your project can use them

Note that repeatedly updating dependencies and running `go mod vendor` will add unneccesary lines to your go.mod file.  Run `go mod tidy` every so often to clean it up.  

### Building binaries

#### Linux

```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-w -s' -o udb main.go
```

#### Windows

```
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-w -s' -o udb main.go
```

or

```
GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -ldflags '-w -s' -o udb main.go
```

for a 32 bit version.

#### Mac OSX

```
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-w -s' -o udb main.go
```


## Running tests

Simply run `$ go test ./...`
