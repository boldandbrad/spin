# 💿 spin

**A command line last.fm scrobbler for techies.**

Interactively or programmatically scrobble tracks and albums to last.fm from the
terminal. That's it.

## Install

```sh
brew install boldandbrad/tap/spin
```

Or download a [release](https://github.com/boldandbrad/spin/releases) and add
the `spin` binary to your path.

Or build from source.

```sh
git clone https://github.com/boldandbrad/spin.git
cd spin
go build -o spin .
./spin --version  # verify installation
```

## How it works

### User Management

Spin stores session details for last.fm users as profiles. Multiple profiles can
be created to enable scrobbling to different last.fm accounts. If there is only
one profile, Spin will default to it, otherwise the active profile can be
manually set.

Session keys are stored securely in the system keychain (macOS Keychain or
Linux Secret Service). Profile metadata is stored at `~/.config/spin/`.

### Scrobble Modes

- 👤 **TUI mode**: An interactive mode that prompts for artist and release details,
  searches last.fm for the best match, and scrobbles automatically. TUI mode closes
  automatically when a scrobble is submitted, or can be closed with `Ctrl+C`.

- 🤖 **CLI mode**: An automation friendly mode that scrobbles tracks and albums
  directly using details provided as command arguments. Scrobbles are submitted as
  soon as the command is run.

By default, Spin uses the current time for scrobbles. However, both modes
provide ways to set custom timestamps.

## Quick Start

Add a profile for your last.fm account. You will be prompted for your last.fm
password.

```sh
spin profile add <lastfm-username>
```

Then, scrobble some music!

```sh
spin track
OR
spin album
```

## Usage

### Global flags

```sh
spin --version       # print version
spin --help          # print help message
spin --debug         # enable debug logging
```

### Profiles

Adding a profile will prompt for the given username's last.fm password to
authenticate.

```sh
spin profile add <lastfm-username>      # add a last.fm user
```

> At least one profile must exist in order to scrobble.

Other profile actions:

```sh
spin profile list                       # list added profiles
spin profile set <lastfm-username>      # set the active profile
spin profile get                        # get the active profile
spin profile delete <lastfm-username>   # remove a profile
```

> The active profile can also be set using the `-p/--profile` flag in CLI mode.
> See below.

### Scrobble

Spin provides dedicated commands for scrobbling individual tracks and full
albums. Both commands can be used in either TUI or CLI mode.

#### TUI mode

TUI mode interactively prompts for scrobble deetails and is launched when no
arguments are provided:

```sh
spin track      # interactively search for and scrobble a track
spin album      # interactively search for and scrobble an album
```

TUI mode example:

[TODO: add tui mode gif]

Available TUI mode options:
- `--dryrun`: show what would be scrobbled without submitting

#### CLI mode

CLI mode is enabled when arguments are present. Both commands require two
positional arguments: `artist`, and then `track` or `album` respectively.

```sh
spin track <artist> <track>                                 # scrobble track
spin album <artist> <album>                                 # scrobble album
```

Available CLI mode options:
- `-d|--date`: date of listen (YYYY-MM-DD)
- `-t|--timestamp`: time of listen (HH:MM)
- `-p|--profile`: profile to scrobble with (default: active profile)
- `--dryrun`: show what would be scrobbled without submitting

CLI mode examples:

```sh
spin track "Best Frenz" "Replay"                            # scrobble track now
spin track "Joywave" "Nice House" -t 12:46 -p boldandbrad   # specific time today and profile
spin track "Joywave" "Nice House" --dryrun                  # preview without scrobbling
spin album "Coldplay" "X&Y" -t 15:32                        # specific time today
spin album "Electric Guest" "Mondo" -d 2026-01-31 -t 01:14  # specific date and time
```

### History

Review recent scrobbles from the active profile. Useful for validating scrobbles
were successful without launching last.fm in the browser.

```sh
spin history                # list active profile's recent scrobbles
spin history -n 50          # set the number of results
```

## Why use spin?

- 🤖 Scriptable: use it to automatically scrobble locally playing music
- 🎮 Interactive: fun and easy to use on the fly
- 👥 Multi-user: profiles enable scrobbling to multiple accounts
- 🔧 No config: just works out of the box
- 🔒 Secure: session keys stored in system keychain

## License

[MIT](./LICENSE)
