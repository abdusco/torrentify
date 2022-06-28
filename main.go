package main

import (
	_ "embed"
	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"os"
)

//go:embed version.txt
var version string

type torrent struct {
	AnnounceUrls []string
	Comment      string
	CreatedBy    string
	Name         string
	Private      bool
	Root         string
	PieceLength  uint64
}

func main() {
	app := &cli.App{
		Name:        "torrentify",
		Usage:       "Torrent creator",
		ArgsUsage:   "<torrent root>",
		Version:     version,
		Description: "torrentify creates torrent files from given root directory.",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "announce",
				Aliases:  []string{"a"},
				Usage:    "Tracker announce URLs. Separate multiple URLs with commas.",
				EnvVars:  []string{"ANNOUNCE_URL"},
				Required: true,
			},
			&cli.PathFlag{
				Name:      "output",
				Aliases:   []string{"o"},
				Usage:     "Output path. Defaults to stdout",
				Required:  true,
				TakesFile: true,
				Value:     "-",
			},
			&cli.StringFlag{
				Name:  "comment",
				Usage: "Torrent comment",
			},
			&cli.StringFlag{
				Name:    "createdby",
				Usage:   "Torrent creator name",
				EnvVars: []string{"CREATED_BY"},
			},
			&cli.BoolFlag{
				Name:    "private",
				Usage:   "Set torrent as private. Useful for private trackers",
				EnvVars: []string{"PRIVATE"},
			},
			&cli.Uint64Flag{
				Name:    "piecelength",
				Usage:   "Torrent piece length.",
				EnvVars: []string{"PIECE_LENGTH"},
				Value:   1024 * 1024,
			},
		},
		HideHelpCommand: true,
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() != 1 {
				return cli.ShowAppHelp(ctx)
			}

			t := &torrent{
				AnnounceUrls: ctx.StringSlice("announce"),
				Comment:      ctx.String("comment"),
				CreatedBy:    ctx.String("createdby"),
				Name:         ctx.String("name"),
				Private:      ctx.Bool("private"),
				Root:         ctx.Args().Get(0),
				PieceLength:  ctx.Uint64("piecelength"),
			}

			outputPath := ctx.Path("output")
			var w io.Writer
			if outputPath == "-" {
				w = os.Stdout
			} else {
				f, err := os.Create(outputPath)
				if err != nil {
					return errors.Wrap(err, "create file")
				}
				w = f
			}

			return makeTorrent(t, w)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func makeTorrent(t *torrent, w io.Writer) error {
	mi := metainfo.MetaInfo{
		AnnounceList: make([][]string, 0),
	}
	for _, a := range t.AnnounceUrls {
		mi.AnnounceList = append(mi.AnnounceList, []string{a})
	}
	mi.SetDefaults()
	if len(t.Comment) > 0 {
		mi.Comment = t.Comment
	}
	if len(t.CreatedBy) > 0 {
		mi.CreatedBy = t.CreatedBy
	}
	info := metainfo.Info{
		PieceLength: int64(t.PieceLength),
		Private:     &t.Private,
	}
	err := info.BuildFromFilePath(t.Root)
	if err != nil {
		return errors.Wrap(err, "hash files")
	}
	if t.Name != "" {
		info.Name = t.Name
	}
	mi.InfoBytes, err = bencode.Marshal(info)
	if err != nil {
		return errors.Wrap(err, "marshall torrent")
	}

	return mi.Write(w)
}
