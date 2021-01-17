# Neon

Neon is a static site generator built by me and for me.

## Install

To install, you must have the Go toolchain installed at version 1.11 or later.

1. `go build` (or `GO111MODULE=on go build` if you've cloned this directory
   into your `$GOPATH`)
2. `mv ./neon $SOMEWHERE_IN_PATH`

## Usage

### Initializing a neon project

Neon projects have a file called `neon.yaml` in their root. Neon commands run
in any subdirectory of the project root will be run against the project. The
`neon.yaml` file for my blog currently looks like this (check out
[./config/config.go](./config/config.go) for marginally more information about
the configuration options):

```yaml
theme: default
live_reload_port: 8081
title: Craig Weber
description: "My blog; mostly tech"
author: Craig Weber
theme_config: # these params are only passed into the theme templates
  site_title: Craig Weber
  copyright: Craig Weber 2016
  contact:
    - url: http://twitter.com/weberc2
      text: Twitter
    - url: http://github.com/weberc2
      text: Github
  analytics_url: https://azd7mnx3pk.execute-api.us-east-2.amazonaws.com
index_page_size: 10
site_root: https://weberc2.github.io
code_highlight_theme: friendly
```

Ok, now we have our `neon.yaml` file--we should also make a `./posts` directory
with our posts represented as individual markdown files. These markdown files
should have some frontmatter, such as:

```md
---
Title: Deploying Go apps on Docker scratch images
Date: 2018-11-04
---
```

`Title` and `Date` are pretty important, but you can hang other metadata off of
this too.

Lastly, you should have a `./themes` directory with your themes. You can start
by checking out [`./themes/mono`](./themes/mono). Themes are just Go templates.
The `theme_config` from the `neon.yaml` file gets passed, so you have to make
sure the `theme_config` is populated as required by the theme.

### Building a neon project

Ok, so now that we've initialized a neon project, we should try to build it.
That much is easy, it's just `neon build` from anywhere in project. This will
create an `$NEON_PROJECT_ROOT/_output/` directory if none exists and dump the
files there. Note that if the output directory already exists, `neon build`
will overwrite files, but it will not delete any deleted or renamed posts. In
lieu of that, the best practice is to delete the whole output directory
before running `neon build`.

### Deploying the output directory

The output directory is ready-to-serve. Because I deploy to BitBucket pages (a
git repo), I just `rm -r $TARGET_REPO_ROOT/*` followed by
`cp -r $NEON_PROJECT_ROOT/* $TARGET_REPO_ROOT` and then I push the repo to
publish it.

### Iterating with `neon serve`

Running `neon serve` in your project directory will run a local development
server on port 8080. This is not a production server. It takes a `--port` flag
if you want to override the default port.
