# elixxir/user-discovery-bot

[![pipeline status](https://gitlab.com/elixxir/user-discovery-bot/badges/master/pipeline.svg)](https://gitlab.com/elixxir/user-discovery-bot/commits/master)
[![coverage report](https://gitlab.com/elixxir/user-discovery-bot/badges/master/coverage.svg)](https://gitlab.com/elixxir/user-discovery-bot/commits/master)

The user discovery bot helps users make first contact with other users. Users can search for other users using a string key (i.e. email address or phone number) and, if the user discovery bot finds a match for that user with the hash of the string, it will return a key ID. The user and the bot can then do a key exchange with the public key that the bot returns after the user queries that key ID to facilitate transfer of information that they need to talk to the user.

##Command-line options

|Long flag|Short flag|Effect|Example|
|---|---|---|---|
|--config| |Specify a different configuration file|--config udb2.yaml|
|--help|-h|Shows a help message|-h|
|--verbose|-v|Prints more log messages|-v|
|--version|-V|Prints generated version information for the UDB and its dependencies. To regenerate log messages, run `$ go generate cmd/version.go`.|-V|
|--ndf|-n|Path to the network definition file|-V|

## Example configuration

Note: Yaml prohibits the use of tabs. If you put tabs in your config file, the UDB will fail to parse it.

```yaml
# Path where UDB will store its logs
logPath: "udb.log"
# Path where UDB will store session file
sessionfile: "udb.session"

# Database connection information
dbUsername: "cmix"
dbPassword: ""
dbName: "cmix_server"
dbAddress: ""q
```

## Running

### Installing dependencies

`$ glide up` should automatically download or update all dependencies and place them in the vendor/ folder. If it's not working correctly, try removing `~/.glide/` and `./glide.lock` and trying again to clear the cache. If it still doesn't work, make sure you're pointing to the right versions and that you have access to all the repositories that are getting downloaded.

### Running tests

Simply run `$ go test ./...`
