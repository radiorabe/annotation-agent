# RaBe Annotation Agent

The RaBe Annotation Agent is used to annotate audio files and other assets at RaBe. It is also an experiment in using linked data at [RaBe](https://rabe.ch).

* is triggered via several keys on a amqp topic
* downloads files from [our archive](https://archiv.rabe.ch)
* stores [speech/music segmentation](https://github.com/bbc/bbc-vamp-plugins) in annnotations
* stores [audiowaveform](https://github.com/bbc/audiowaveform) dat files in an object store and links them with files via an annotation
* receives events from [acrcloud](https://acrcloud.api.rabe.ch) and links them with a show (unfinished, doesn't work yet)

## Development

```bash
# clone the repo
git clone https://github.com/radiorabe/annotation-agent.git
cd annotation-agent

# run the command line locally
go run main.go --help

# build a binary
go build main.go -o annotation-agent
```

### pre-commit hook

#### pre-commit configuration

```bash
# setup hooks
pre-commit install

# run them all
pre-commit run -a
```

### Release Process

Create a git tag and push it to this repo or use the git web ui.

This is built on GitHub Actions and uses a `GH_PAT_TOKEN` secret to work. The access key must
have repo, read:packages, write:packages and delete:packages in it's scope.

## License
This software is free software: you can redistribute it and/or modify it under
the terms of the GNU Affero General Public License as published by the Free
Software Foundation, version 3 of the License.

## Copyright
Copyright (c) 2020 [Radio Bern RaBe](http://www.rabe.ch)
