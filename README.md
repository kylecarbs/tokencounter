# tokencounter

A simple CLI to count OpenAI tokens for a set of files.

- Supports `.gitignore`
- Match mimetypes to include
- Specify any model

> [!NOTE]
> Uses [tiktoken-go](https://github.com/pkoukk/tiktoken-go) to tokenize. It may not be 100% accurate, but should be close.

```sh
# Install the CLI
go install github.com/kylecarbs/tokencounter@latest

# Compute tokens for the current directory recursively
tokencounter

# List supported models
tokencounter models
```

## Example

Counting the number of tokens in [coder/coder](https://github.com/coder/coder):

```sh
dev ~/p/c/coder (main *s =) ~/tokencount
using ignore file: .gitignore
adding files: 100
adding files: 200
adding files: 300
adding files: 400
adding files: 500
adding files: 600
adding files: 700
adding files: 800
adding files: 900
adding files: 1000
adding files: 1100
adding files: 1200
adding files: 1300
adding files: 1400
skipped "video/mp4" 7 times
skipped "image/jpeg" 3 times
skipped "image/bmp" 1 times
skipped "image/svg+xml" 243 times
skipped "application/gzip" 2 times
skipped "application/x-ndjson" 4 times
skipped "application/json" 87 times
skipped "application/x-tar" 2 times
skipped "application/zip" 1 times
skipped "image/png" 4097 times
skipped "image/gif" 2 times
skipped "video/quicktime" 2 times
skipped "application/octet-stream" 15 times
skipped "image/x-icon" 2 times
skipped "image/webp" 3 times
waiting for processing to finish
total tokens: 5631222
```
