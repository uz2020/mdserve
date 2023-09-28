# A simple service for the emacs users

Inspired by impatient-mode, an emacs package for lively reload of showing previews of markdown files, I wrote this simple service in Golang.

Not long ago, I had a need to read some Ethereum's specs in markdown format. Github's rendering of these files are awesome, but unfortunately, network connectivity is required for reading them. It not only slows down the speed I read the specs, but the bad news is sometimes I don't have stable network connectivity at hand. After a little struggling, I decided to write my own tools to settle this matter.

This HTTP service pushes the new content of a file to the browser whenever it's modified.

A huge thanks to the team of [bytemd](https://github.com/bytedance/bytemd). Remarkable work, really well done.

## Installation

The service will be listening at `127.0.0.1:8080`

```
make

mdserve
```

## Preview of a markdown file

```
google-chrome http://localhost:8080/?path=/home/l/bytemd/README.md&server=localhost:8080
```

## With emacs

You should make a little modification to the snippets below to use the correct host and post of the golang service.

```
(defun browsemd (name) (browse-url (concat "http://localhost:8080?path=" name "&server=127.0.0.1:8080")))

(defun browsemdbuffer () (interactive) (browsemd (buffer-file-name)))
```

Then by invoking `M-x browsemdbuffer`, this will open your brower with the url constructed that's conformed to the service's requirements.


