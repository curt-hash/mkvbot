# mkvbot

[![Check](https://github.com/curt-hash/mkvbot/actions/workflows/check.yml/badge.svg)](https://github.com/curt-hash/mkvbot/actions/workflows/check.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/curt-hash/mkvbot.svg)](https://pkg.go.dev/github.com/curt-hash/mkvbot)
[![Release](https://img.shields.io/github/v/release/curt-hash/mkvbot)](https://github.com/curt-hash/mkvbot/releases)

A simple CLI for ripping Blu-rays and DVDs with MakeMKV. Insert a disc, it picks the best title, looks it up on IMDb, and rips it to a folder named for Plex.

## What it does

1. Scans the disc and scores titles using configurable heuristics
2. Searches IMDb for metadata (name, year, ID)
3. Prompts you to confirm the metadata
4. Rips the best title to `<output>/<Name> (<Year>) {imdb-ttXXXXXXX}/`
5. Ejects the disc and beeps

Since auto-selection isn't perfect, you confirm metadata before it's committed. Use `--ask-title` to always pick the title manually.

It's single-threaded, single-drive, and single-movie — insert, wait, eject, repeat. No scheduling, no queue. If you need something more powerful, check out [automatic-ripping-machine](https://github.com/automatic-ripping-machine/automatic-ripping-machine).

One static binary, no runtime deps besides MakeMKV. Runs fine on a Pi over SSH.

## Requirements

- [MakeMKV](https://makemkv.com/) with a valid license key (needed for Blu-ray)
- `makemkvcon` on PATH (Linux/macOS) or in `Program Files (x86)\MakeMKV\` (Windows)

## Install

### Release

Find a package for your OS and architecture on the [Releases](https://github.com/curt-hash/mkvbot/releases) page (deb, rpm, apk, tarballs).

### Homebrew

```sh
brew tap curt-hash/homebrew-mkvbot
brew install mkvbot
```

### Go

```sh
go install github.com/curt-hash/mkvbot@latest  # or a specific tag, e.g. @v0.3.0
```

## Quick Start

```sh
mkvbot --create-config   # writes mkvbot.toml
mkvbot --create-profile  # writes makemkv.xml
```

> **Note:** Prior versions used `profile.xml`. Rename it to `makemkv.xml` or set `makemkv_profile_path = "profile.xml"` in your config file.

Edit `mkvbot.toml` and set `output_dir_path` to where you want rips to go. Everything else has sensible defaults.

Then:

```sh
mkvbot                    # reads mkvbot.toml from cwd
mkvbot -o /tmp/rips       # override output dir
mkvbot --ask-title        # always pick the title manually
mkvbot --debug            # verbose logging
```

Insert a disc and run `mkvbot`. The TUI shows drive/disc/movie/title info, a log, and input forms when decisions are needed. `Ctrl+C` exits cleanly.

## Beta Key

MakeMKV's beta key expires every few months. To fetch the current key from the [MakeMKV forum](https://forum.makemkv.com/forum/viewtopic.php?t=1053) and register it automatically:

```sh
mkvbot update-key
```

On Windows this writes to the registry; on Linux/macOS it updates `~/.MakeMKV/settings.conf`.

## Rip Filenames

Follows [Plex's naming convention](https://support.plex.tv/articles/naming-and-organizing-your-movie-media-files):

```
Movies/
  Blade Runner (1982) {imdb-tt0083658}/
    Blade Runner (1982) {imdb-tt0083658}.mkv
```

If the name is missing, uses `Unknown Title <timestamp>` to avoid collisions and keep Plex happy. If only the IMDb ID is missing, uses `Name (Year)` without the tag.

## Audio and Subtitles

Edit `app_DefaultSelectionString` in `makemkv.xml` to change which streams are selected. The default picks English audio and subtitles. `makemkvcon` doesn't read the GUI preferences, so the profile file is required.

## Config

`mkvbot.toml` is optional — all settings can be CLI flags. Flags override config file values.

| Setting | Default | |
|---------|---------|--|
| `output_dir_path` | `.` | where rips go |
| `makemkv_profile_path` | `makemkv.xml` | profile file |
| `cache_size` | `1024` | read cache (MiB) |
| `min_length` | `1800` | minimum title length (seconds) |
| `quiet` | `false` | suppress beep |
| `ask_for_title` | `false` | prompt for title |
| `best_title_heuristic_weights.*` | see below | scoring weights |

Title scoring weights (all matched titles accumulate their weight):

| Heuristic | Weight | |
|-----------|--------|--|
| `longest` | 1000 | longest title |
| `angle_one` | 300 | has angle 1 |
| `most_chapters` | 200 | most chapters |
| `most_streams` | 100 | most streams |

Tune them in `mkvbot.toml` or with flags like `--longest-title-weight 2000`.

## Known Issues

- Only the first optical drive is used if multiples exist
- TV shows aren't supported (no episode mapping)
- IMDb lookup is HTML scraping, so fragile to page structure changes. If it breaks, enter metadata manually at the prompt.
- Rips block until done; no queue or scheduling

## License

MIT
