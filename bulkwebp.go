package main

import (
	"context"
	"flag"
	"fmt"
	_ "image/png"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/kolesa-team/go-webp/encoder"
	"golang.org/x/sync/errgroup"
)

func webpPath(from string) string {
	ext := filepath.Ext(from)
	without := strings.TrimSuffix(from, ext)
	return without + ".webp"
}

func isSupportedExt(path string) bool {
	switch filepath.Ext(path) {
	case ".png":
		return true
	default:
		return false
	}
}

func run(ctx context.Context) error {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [options] <files or directories...>\n\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()
	input := flag.Arg(0)
	if input == "" {
		flag.Usage()
		os.Exit(2)
	}

	options, err := encoder.NewLosslessEncoderOptions(encoder.PresetDefault, 6)
	if err != nil {
		return fmt.Errorf("create encoder options: %w", err)
	}

	conv := &converter{
		Options: options,
	}

	eg, ctx := errgroup.WithContext(ctx)

	for i := 0; i < flag.NArg(); i++ {
		base := flag.Arg(i)
		eg.Go(func() error {
			return filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
				if ctx.Err() != nil {
					return ctx.Err()
				}

				if d.IsDir() {
					return nil
				}
				if !isSupportedExt(path) {
					return nil
				}

				eg.Go(func() error {
					err := conv.Do(path)
					if err != nil {
						return fmt.Errorf("(%q) %w", path, err)
					}

					fmt.Printf("%q -> %q", path, webpPath(path))
					return nil
				})

				return nil
			})
		})
	}

	return eg.Wait()
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		<-ctx.Done()
		cancel()
	}()

	err := run(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
