web-stream
==========
web-stream is an http-server that, when connected to, yields a
web-interface for navigating the filesystem below the configured
`root`.

When the user clicks on a video file, web-stream starts encoding
it to webm using ffmpeg and serves the resulting file, while it's
still being encoded. Because webms can be viewed with most mobile
browsers, this can be used to stream movies from your PC to your
phone.

Table of contents
-----------------
- [Installation](#installation)
- [Configuring the server](#configuring-the-server)
- [Configuring the web-interface](#configuring-the-web-interface)

## Requirements:
- As web-stream is written in go, it requires an up-to-date
[go environment](https://golang.org/doc/install) to work properly.
- To be able to encode the videos, web-stream requires
[ffmpeg](https://www.ffmpeg.org/).

## Installation
To install it, simply run the following commands:
```
go get github.com/froozen/web-stream
$GOPATH/src/github.com/froozen/web-stream/install.sh
```
Now you're good to go! You can now move on to the configuration.

## Configuring the server
web-stream is configured through a [json-encoded](http://www.json.org/)
config file, that can be found under `~/.web-stream/config.json`.

The default configuration file looks something like this:
```json
{
    "port": 2223,
    "args": [
        "-threads", "8",
        "-quality", "realtime"
    ],
    "hooks": [],
    "filetypes": [".mkv", ".mp4"]
}
```

#### Changing the `root`-value
The first thing you will want to do is changing the `root` value to
point to the directory you keep your movies in. The default is
`~/Videos`.

To do that, simply add a line like this:
```json
"root": "/path/to/the/dir"
```

#### The `args`-value
The default arguments for the ffmpeg-command look something like this:
```
ffmpeg -i file <args> file.webm
```
Without any `args` this won't do much more than encoding the video, which
will be pretty slow. If you want to speed up the encoding process, you'll
need to pass additional arguments to ffmpeg. The `args` in the default
configuration file are a pretty good start, but you might still want to configure
them to your liking.

**Note:** Each argument must be a single string in the list. For example, to
automaticly scale the video, you'd need to add:
```json
"-s", "1280x720"
```

#### Dynamicly generating additional `args` using `hooks`
There are cases, when you'd want to dynamicly generate arguments instead
of staticly predefining them. Luckily, web-stream has `hooks` for that.

`hooks` are commands that are called every time a file is encoded.
They receive the filename as the first argument and print a json-list
of strings, that will then be added to the `args` for the encoding.

**Example:** This is a `hook` that adds subtitles to the video
```bash
#!/bin/bash

# Generate the tempfile name
a=$(mktemp -u sub-XXXXX.ass -p /tmp)

# Extract the subtitles
ffmpeg -i "$1" -threads 8 -an -vn -loglevel quiet -c copy $a

# Print the arguments
echo "[\"-vf\", \"subtitles=$a\"]"
```

**Note:** A `hook` printing something invalid will result in web-stream
abborting, so be careful.

#### Other values
These are a few other configuration options, which are rather self explanatory
and therefore don't need a whole paragraph of text for each one:
- `port`: This is the port web-stream will listen and serve on.
- `filetypes`: This is a list of the filetypes that web-stream will encode
and display in its web-interface
- `delay`: This is an integer value representing the time in senconds that
is waited before serving the file for the first time. This can be viewed
as a sort of "buffering" mechanism. The default value is 3 seconds.

## Configuring the web-interface
Because of my lack of web-design skills, the default web-interface doesn't look
all that great. Luckily, the templating system is pretty simple, so it should
be no problem to edit the files yourself.

The templating system replaces every comment in the form of `<!-- Data:valuename -->`
with the data corresponding to `valuename`.

All of the following files can be found in `~/.web-stream/web`
- `style.css` is served over the /style url and is the css-stylesheet for the
pages
- `page.html` is the template for the whole page. The insertable values are:
    - `DIRNAME`, which is the name of the directory corresponding to the page
    - `CODE`, which is the code generated for the files and directories in the directory
- `file.html` and `dir.html` contain the templates for the code snippets used
to represent a file or directory contained withing the displayed directory. All
of them together form the `CODE`-value. The insertable values are:
    - `PATHNAME`, which is the path leading to the currently displayed directory
    - `ITEMNAME`, which is the file- or directory name
