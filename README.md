# torrentify

A mini CLI application to create torrent files from a directory.

## Usage

**torrentify** requires the root directory of the files to be specified as the only argument.
The remaining parameters are specified as flags,
some of which can be set using environment variables.

```

```shell
> ./torrentify --help

NAME:
   torrentify - Torrent creator

USAGE:
   torrentify [global options] <torrent root>

DESCRIPTION:
   torrentify creates torrent files from given root directory.

GLOBAL OPTIONS:
   --announce value, -a value  Tracker announce URLs. Separate multiple URLs with commas.  (accepts multiple inputs) [$ANNOUNCE_URL]
   --comment value             Torrent comment
   --createdby value           Torrent creator name [$CREATED_BY]
   --help, -h                  show help (default: false)
   --output value, -o value    Output path. Defaults to stdout (default: "-")
   --piecelength value         Torrent piece length. (default: 1048576) [$PIECE_LENGTH]
   --private                   Set torrent as private. Useful for private trackers (default: false) [$PRIVATE]
   --version, -v               print the version (default: false)
```
