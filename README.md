# ds 🔢
ds is a simple CLI tool to compute statistics of datasets

## Installation

```
brew install bwagner5/ds
```

## Usage

```
$ ds is a CLI tool to compute stats for data sets

Usage:
  ds [flags]

Flags:
  -f, --file string   Input file to compute statistics for
  -h, --help          help for ds
  -v, --version       the version
```

```
$ ds -f file-with-numbers ... or
$ cat file-with-numbers | ds ... or
$ ds < file-with-numbers
n:                   10
mean:                5.5
median:              5.5
std dev:             2.87228
min:                 1
max:                 10
P99.99:              9.5
P99:                 9.5
P95:                 9.5
P75:                 7.5
P25:                 2.5
P5:                  1.5
P1:                  1.5
P0.01:               1.5
Top Freq Num:        1
Top Freq:            1
> Top Freq:          9
< Top Freq:          0
```
