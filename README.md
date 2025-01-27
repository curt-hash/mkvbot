# mkvbot

## Overview

`mkvbot` streamlines ripping the "best" title from a Blu-ray or DVD using
MakeMKV. It runs a simple processing loop:

1. Wait for a disc
1. Search IMDb for movie metadata needed for Plex-friendly file names
1. Identify and rip the best title
1. Eject the disc

`mkvbot` has limited functionality and is not fully automated. That might work
for you. It works for me. If you are looking for something more complex, maybe
check out a project like
[automatic-ripping-machine](https://github.com/automatic-ripping-machine/automatic-ripping-machine).

`mkvbot` is designed to be lightweight and portable. It is a single executable
that can run on an old Raspberry Pi.

## Installation

`mkvbot` requires [makemkv](https://makemkv.com/).

### Release

Find an appropriate package for your operating system and architecture on the
[Releases](https://github.com/curt-hash/mkvbot/releases) page.

### Homebrew

```sh
brew tap curt-hash/homebrew-mkvbot
brew install mkvbot
```

### Snap

_This is not working yet (pending classic confinement approval)._

```
snap install mkvbot --classic
```

## Usage

`mkvbot` is a terminal program with a text-based user interface (TUI). Open a
terminal and run the executable (`mkvbot` or `mkvbot.exe` depending on
platform):

```sh
mkvbot.exe --output-dir Z:\\path\\to\\Movies
```

Run `mkvbot.exe -h` to see all of the command-line options.

Since it does not always pick the correct title or movie metadata, it currently
prompts for confirmation. It will also prompt you to choose the best title if
there is a tie. That may change as it gets smarter.

Audio track and subtitles selection is based on the value of
`app_DefaultSelectionString` in [profile.xml](profile.xml). For whatever reason,
`makemkvcon` (the CLI application) does not seem to honor the selection string
set in the GUI application preferences.
