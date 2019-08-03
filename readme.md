# Twitter CLI for collect tweets
[![CircleCI](https://circleci.com/gh/mpppk/twitter.svg?style=svg)](https://circleci.com/gh/mpppk/twitter)
[![Build status](https://ci.appveyor.com/api/projects/status/39y6e4o6khig6mct?svg=true)](https://ci.appveyor.com/project/mpppk/twitter)
[![codecov](https://codecov.io/gh/mpppk/twitter/branch/master/graph/badge.svg)](https://codecov.io/gh/mpppk/twitter)
[![GoDoc](https://godoc.org/github.com/mpppk/twitter?status.svg)](https://godoc.org/github.com/mpppk/twitter)

## Installation

Download from [GitHub Releases](https://github.com/mpppk/twitter/releases).  
Or `go get github.com/mpppk/twitter` (go modules must be enabled)

## Usage

### Search
Search tweets by query and some options.  
Results are stored in local file DB. (You can specify the DB path by --db-path flag.)  
If you want to download images which contained in tweets, execute 'download images' command after search command.

Search command manage two state, 'minID' and 'maxID', which decide tweet ID range when searching. So search command can resume even if process is interrupted.

```bash
$ twitter search \
  ---db-path tweets.db \
  -query [some_words] \
  --exclude retweets \
  --filter images 
  
  100 tweets are saved. (1157495237557678080-1157489067203760128)
  maxID is updated => 1157489067203760127
  100 tweets are saved. (1157489006810206210-1157483911351377921)
  maxID is updated => 1157483911351377920
  ...
```

### Download images
Download images which contained tweets from DB file.  
You must execute 'search' command first for collect tweets to DB.

```bash
$ twitter download images -db-path tweets.db
```

### Configuration
Each DB file has two state, 'minID' and 'maxID', which decide tweet ID range when searching.  
These values are updated automatically by search command, but you can also update manually through 'config' command.

```bash
$ twitter config set [maxID|minID] [new tweet ID]
```

```bash
$ twitter config get [maxID|minID]
=> maxID / minID will be printed
```
