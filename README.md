# 💿 spin

**A command line last.fm scrobbler for techies.**

Interactively or programmatically scrobble tracks and albums to last.fm from the
terminal. That's it.

## Install

```sh
brew tap boldandbrad/tap
brew install spin
```

Or download a release and add it to your path.

> Coming soon.

Or build from source.

> Coming soon.

## How it works

### User Management

Spin stores session details for last.fm users as profiles. Multiple profiles can
be created to enable scrobbling to different last.fm accounts. If there is only
one profile, Spin will default to it, otherwise the active profile can be
manually set.

User session data is stored at [TODO: finalize data storage location].

### Scrobble Modes

- 👤 **TUI mode**: A human friendly mode that interactively prompts for track or
album details, searches last.fm for matches, and allows the user to select the
correct release to scrobble. TUI mode closes automatically when a scrobble is
submitted, or can be closed with `Ctrl+C`.

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
spin -v/--version       # print version
spin -h/--help          # print help message
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

Tracks and albums are scrobbled with their own dedicated commands.

#### TUI mode

TUI mode is launched when no arguments are provided:

```sh
spin track      # interactively search for and scrobble a track
spin album      # interactively search for and scrobble an album
```

[TODO: add tui mode gif]

#### CLI mode

CLI mode is enabled when arguments are present. In this mode both the track and
album commands require two positional arguments: the `artist`, and then `track`
or `album` respectively.

```sh
spin track <artist> <track>                                 # scrobble track
spin album <artist> <album>                                 # scrobble album
```

Available CLI mode options:
- `-d|--date`: date of listen (default: current day)
- `-t|--timestamp`: time of listen (default: current time)
- `-p|--profile`: profile to scrobble with (default: active profile)

CLI mode examples:

```sh
spin track "Best Frenz" "Replay"                            # scrobble track now
spin track "Joywave" "Nice House" -t 12:46 -p boldandbrad   # specific time today and profile
spin album "Coldplay" "X&Y" -t 15:32                        # specific time today
spin album "Electric Guest" "Mondo" -d 2026-01-31 -t 01:14  # specific date and time
```

### History

```sh
spin history                # list active profile's recent scrobbles
spin history -n/--limit 50  # set the number of results
```

## Why use spin?

- 🤖 Scriptable: use it to automatically scrobble locally playing music
- 🎮 Interactive: fun and easy to use on the fly
- 👥 Multi-user: profiles enable scrobbling to multiple accounts
- 🔧 No config: just works out of the box

## License

[MIT](./LICENSE)
