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
	"time"
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

type TorrentCreatorFunc func(t *torrent, w io.Writer) error
type FileCreatorFunc func(path string) (io.Writer, error)
type Runner func(args []string) error

type App struct {
	cli          Runner
	torrentMaker TorrentCreatorFunc
	fileCreator  FileCreatorFunc
}

func (a *App) Run(args []string) error {
	return a.cli(args)
}

func NewApp(fileCreator FileCreatorFunc, torrentCreator TorrentCreatorFunc) *App {
	cliApp := &cli.App{
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
				Name:     "name",
				Usage:    "Torrent name",
				Required: true,
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
			if ctx.NArg() != 1 {
				log.Printf("root directory is not specified")
				return cli.ShowAppHelp(ctx)
			}

			rootDir := ctx.Args().Get(0)

			t := &torrent{
				AnnounceUrls: ctx.StringSlice("announce"),
				Comment:      ctx.String("comment"),
				CreatedBy:    ctx.String("createdby"),
				Name:         ctx.String("name"),
				Private:      ctx.Bool("private"),
				Root:         rootDir,
				PieceLength:  ctx.Uint64("piecelength"),
			}

			outputPath := ctx.Path("output")
			w, err := fileCreator(outputPath)
			if err != nil {
				return err
			}

			return torrentCreator(t, w)
		},
	}
	return &App{
		cli: cliApp.Run,
	}
}

func main() {
	app := NewApp(
		createFile,
		makeTorrent,
	)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func createFile(path string) (io.Writer, error) {
	if path == "-" {
		return os.Stdout, nil
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, errors.Wrap(err, "create file")
	}
	return f, nil
}

func makeTorrent(t *torrent, w io.Writer) error {
	mi := metainfo.MetaInfo{
		AnnounceList: make([][]string, 0),
	}
	for _, a := range t.AnnounceUrls {
		mi.AnnounceList = append(mi.AnnounceList, []string{a})
	}

	mi.CreationDate = time.Now().Unix()
	if len(t.Comment) > 0 {
		mi.Comment = t.Comment
	}

	if len(t.CreatedBy) > 0 {
		mi.CreatedBy = t.CreatedBy
	} else {
		mi.CreatedBy = "torrentify"
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
