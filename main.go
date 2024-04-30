package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"sync/atomic"

	"github.com/coder/serpent"
	"github.com/gabriel-vasile/mimetype"
	"github.com/pkoukk/tiktoken-go"
	ignore "github.com/sabhiram/go-gitignore"
)

func main() {
	var (
		model          string
		matchMime      string
		ignoreFilePath string
		verbose        bool
	)

	cmd := &serpent.Command{
		Use: "tokencounter [directory]",
		Handler: func(i *serpent.Invocation) error {
			logger := log.New(os.Stderr, "", 0)
			mimeMatcher, err := regexp.Compile(matchMime)
			if err != nil {
				return err
			}
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			if len(i.Args) > 0 {
				dir = i.Args[0]
			}
			ignoreFile, err := ignore.CompileIgnoreFile(filepath.Join(dir, ignoreFilePath))
			if err != nil {
				logger.Printf("not using ignore file: %s", err)
			} else {
				logger.Printf("using ignore file: %s", ignoreFilePath)
			}
			encoding, err := tiktoken.EncodingForModel(model)
			if err != nil {
				return err
			}
			var filesMu sync.Mutex
			skippedFiles := map[string]uint64{}
			handledFiles := map[string]uint64{}
			var tokenCount atomic.Int64
			var wg sync.WaitGroup
			err = filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}
				if ignoreFile != nil && ignoreFile.MatchesPath(path) {
					return nil
				}
				wg.Add(1)
				go func() {
					defer wg.Done()
					content, err := os.ReadFile(path)
					if err != nil {
						logger.Printf("error reading %s: %s", path, err)
						return
					}
					mime := mimetype.Detect(content)
					if !mimeMatcher.MatchString(mime.String()) {
						if verbose {
							logger.Printf("skipping %s (%s)", path, mime.String())
						}
						filesMu.Lock()
						skippedFiles[mime.String()]++
						filesMu.Unlock()
						return
					}
					filesMu.Lock()
					handled := handledFiles[mime.String()]
					handled++
					filesMu.Unlock()
					if handled%100 == 0 {
						logger.Printf("processed %q %d times", mime.String(), handled)
					}
					if verbose {
						logger.Printf("tokenizing %s (%s)", path, mime.String())
					}
					tokens := encoding.Encode(string(content), nil, nil)
					tokenCount.Add(int64(len(tokens)))
				}()
				return nil
			})
			if err != nil {
				return err
			}
			for mime, count := range skippedFiles {
				logger.Printf("skipped %q %d times", mime, count)
			}
			logger.Printf("waiting for processing to finish")
			wg.Wait()
			logger.Printf("total tokens: %d", tokenCount.Load())
			return nil
		},
		Options: serpent.OptionSet{{
			Flag:          "model",
			FlagShorthand: "m",
			Description:   "Model to use for tokenization.",
			Default:       "gpt-4",
			Value:         serpent.StringOf(&model),
		}, {
			Flag:          "ignore-file",
			FlagShorthand: "i",
			Description:   "File with patterns to ignore.",
			Default:       ".gitignore",
			Value:         serpent.StringOf(&ignoreFilePath),
		}, {
			Flag:          "verbose",
			FlagShorthand: "v",
			Description:   "Print verbose output.",
			Default:       "false",
			Value:         serpent.BoolOf(&verbose),
		}, {
			Flag:        "match-mime",
			Description: "Only process files with the given MIME type.",
			Value:       serpent.StringOf(&matchMime),
			Default:     "text/.*",
		}},
		Children: []*serpent.Command{{
			Use: "models",
			Handler: func(i *serpent.Invocation) error {
				for model := range tiktoken.MODEL_TO_ENCODING {
					fmt.Println(model)
				}
				return nil
			},
		}},
	}

	err := cmd.Invoke().WithOS().Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
