# torrentify

A mini CLI application to create torrent files from a directory.

## Usage

```shell
> ./torrentify --help

NAME:
   torrentify - torrent creator

USAGE:
   torrentify [global options] <torrent root>

VERSION:
   v0.0.1-fa5d041

DESCRIPTION:
   torrentify creates torrent files from a directory.

GLOBAL OPTIONS:
   --announce value, -a value  tracker announce urls  (accepts multiple inputs) [$ANNOUNCE_URL]
   --comment value             torrent comment
   --createdby value           torrent creator name [$CREATED_BY]
   --help, -h                  show help (default: false)
   --output value, -o value    torrent file path
   --private                   set torrent as private (default: false) [$PRIVATE]
   --version, -v               print the version (default: false)
```
