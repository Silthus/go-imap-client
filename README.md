# go-imap-client

A command line client implementing the [go-imap][go-imap] library to search IMAP mailboxes.

## Installation

Download the binary for your operating system from the [releases page][releases], and you are done :rocket:.

## Usage

You can download the binary from [releases][releases] or run it as the Docker container `silthus/go-imap-client`.

```bash
docker run -it silthus/go-imap-client -e IMAP_CLI_USERNAME=test -e IMAP_CLI_PASSWORD=test --server imap.google.com:993 --tls search test
```

You can always find the up-to-date and full help by running `go-imap-client -h`. Here is a list of possible flags and
their options.

| Flag            | Short | Description                                                                                         | Example                          | Default                                                                       |
|-----------------|-------|-----------------------------------------------------------------------------------------------------|----------------------------------|-------------------------------------------------------------------------------|
| `--config`      |       | Specify a config file to use. See [configuration](#Configuration) for more details.                 | `--config my-config.yaml"`       | Searches for a config in various places. See [configuration](#Configuration). |
| `--server`      | `-s`  | The IMAP server to connect against. Including its port.                                             | `--server "imap.google.com:993"` |                                                                               |
| `--username`    | `-u`  | The username used to authenticate.                                                                  | `--username "my.user@gmail.com"` |                                                                               |
| `--password`    | `-p`  | The password of the user. Use single quotes <kbd>'</kbd> to wrap passwords with special characters. | `--password 'password'`          |                                                                               |
| `--tls`         |       | Set to connect against servers using TLS. (required for most mail servers)                          | `--tls`                          | `false`                                                                       |
| `--skip-verify` |       | When using `--tls` can be used to disable the certificate check. Should only be used for testing!   | `--skip-verify`                  | `false`                                                                       |
| `--timeout`     |       | Set a timeout for connecting to the mail server.                                                    | `--timeout 10s`                  | `5s`                                                                          |

### Search

Search the inbox of a mailbox for mails matching a given subject. Here are the `search` specific flags.

| Flag                 | Short | Description                                                    | Example              | Default |
|----------------------|-------|----------------------------------------------------------------|----------------------|---------|
| `--mailbox`          | `-m`  | Set the mailbox/folder to search.                              | `--mailbox Archive`  | `INBOX` |
| `--no-results-error` | `-e`  | Returns an error and exit code if the search finds no results. | `--no-results-error` | `false` |

```bash
go-imap-client --server "imap.google.com:993" --tls --username "my-user@gmail.com" --password 'my_super_secret_PW!' search awesome search term
```

## Configuration

All parameters can also be configured using a YAML configuration or environment variables prefixed with `IMAP_CLI_`.  
The environment variable keys are all **UPPERCASE** and dashes are replaced with <kbd>_</kbd> (underscores).
For example `no-results-error` becomes `IMAP_CLI_NO_RESULTS_ERROR=true` and `username` is `IMAP_CLI_USERNAME=username`.

The configuration looks like this. All values can be overwritten by specifying the command line argument.
Config files named `.go-imap-client.yaml` located in `$HOME` or the executables' directory are sourced automatically.

```yaml
server: "imap.google.com:993"
tls: true
skip-verify: false
username: user@gmail.com
password: 'mycomplex!"#pas''sword' # The password is (note the escaped single quote): mycomplex!"#pas'sword
timeout: 10s

no-results-error: true
mailbox: INBOX
```

[go-imap]: https://github.com/emersion/go-imap

[releases]: https://github.com/Silthus/go-imap-client/releases/latest
