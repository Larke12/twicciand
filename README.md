# twicciand

`twicciand` is the background process for [Twiccian](https://github.com/octotep/twiccian). It's job is to authenticate
with Twitch and handle all communication with the Twitch API and chat.

## Building

To build the application, first make sure [`GOPATH` is
set](https://golang.org/doc/code.html). Next, issue the command:

```go get "github.com/walle/cfg"```

to install config file library. Finally, run `go build` in
the project directory to build the project.

## Authentication

Currently, `twicciand` can only authenticate with Twitch on behalf of the user.
To do so, visit this
[URL](https://api.twitch.tv/kraken/oauth2/authorize?response_type=token&client_id=mya9g4l7ucpsbwe2sjlj749d4hqzvvj&redirect_uri=http://localhost:19210/&scope=user_read+user_follows_edit+channel_read+user_subscriptions+chat_login)
 while the server is running to generate a authentication token. The server
 will then echo this auth token. To avoid going through this process every time
 the server starts up, the token is saved in the twicciand configuration file.

## Configuration File

At startup, the server reads a startup configuration from a file called
`twicciand.conf` in the same directory as the server. Below is an example
configuration:

```
username=USERNAME
token=
```

Make sure you replace `USERNAME` with your twitch username. Build and run the project and follow the directions to generate a twitch auth token.
