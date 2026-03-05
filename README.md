# 💿 spin

**A command line last.fm scrobbler for techies.**

Interactively or programmatically scrobble tracks and albums to Last.fm from the
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
one profile, Spin will default to it, otherwise the active profile can be set.

### Scrobble Modes

TUI mode interactively prompts for track or album details, searches last.fm for
matches, and allows the user to select the correct release to scrobble. This
mode is intended to be human friendly.

CLI mode scrobbles tracks and albums directly using details provided as command
arguments. This mode allows spin to be easily integrated into automations.

## Usage

### Profiles

Adding a last.fm username will prompt for a password to authenticate.

```sh
spin profile add boldandbrad        # add a last.fm user profile
```

Other profile actions:

```sh
spin profile list                   # list added profiles
spin profile set boldandbrad        # set the active profile (automatically set to the first user added)
spin profile get                    # get the active profile
spin profile remove boldandbrad     # remove a last.fm user profile
```

### Scrobble

TUI Mode (interactive):

```sh
spin track      # interactively search for and scrobble a track
spin album      # interactively search for and scrobble an album
```

CLI Mode:

```sh
spin track --artist "Best Frenz" "Replay"                           # scrobble track
spin album --artist "Coldplay" "X&Y"                                # scrobble album
spin track -t 15:32 --artist "Joywave" "Nice House"                 # set specific time today
spin track -t 2026-02-27.15:32 --artist "Metric" "Gold Guns Girls"  # set scrobble date and time
spin album -p boldandbrad --artist "Electric Guest" "Mondo"         # specify a profile
```

### Recents

```sh
spin recent             # list active profile's recent scrobbles
spin recent -n 50       # set the number of results
```

## Why use spin?

- 🤖 Scriptable: use it to automatically scrobble locally playing music
- 🎮 Interactive: fun and easy to use
- 👥 Multi-user: profiles enable srobbling to multiple accounts

## Inspiration

> Coming soon.

## License

[MIT](./LICENSE)
