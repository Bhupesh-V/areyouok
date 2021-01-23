<h1 align="center">AreYouOK?</h1>
<blockquote align="center">A minimal, fast & easy to use URL health checker</blockquote>
<p align="center">
  <img align="center" alt="cat areyouok logo" height="100px" src="https://user-images.githubusercontent.com/34342551/103980534-0da0e000-51a6-11eb-8e67-f4c41599ce1e.png" />
  <br><br>
  <a href="https://goreportcard.com/report/github.com/Bhupesh-V/areyouok">
    <img alt="Go report card" src="https://goreportcard.com/badge/github.com/Bhupesh-V/areyouok" target="_blank">
  </a>
  <a href="https://github.com/Bhupesh-V/areyouok/blob/master/LICENSE">
    <img alt="License" src="https://img.shields.io/github/license/Bhupesh-V/areyouok?color=blue" target="_blank">
  </a>
  <a href="https://github.com/Bhupesh-V/areyouok/issues">
    <img alt="Issues" src="https://img.shields.io/github/issues/Bhupesh-V/areyouok?color=blueviolet" />
  </a>
  <a href="https://twitter.com/bhupeshimself">
    <img alt="Twitter: bhupeshimself" src="https://img.shields.io/twitter/follow/bhupeshimself.svg?style=social" target="_blank" />
  </a>
  <img alt="areyouok-v1.0.0-demo" title="areyouok-v1.0.0-demo" src="https://user-images.githubusercontent.com/34342551/105573121-5a2d1380-5d81-11eb-8ae9-baf6d3d50e3a.gif">
</p>

## Who is `AreYouOk` made for ?

- **OSS Package Maintainers** 📦️:<br>
  With packages comes documentation which needs to be constantly updated & checked for dead/non-functioning URLs.
- **Tech Bloggers** ✍️ :<br>
  If you are someone who writes countless tutorials & owns the source for your website, Use areyouok to make sure your blogs don't contain any non-functioning URLs.
- Literally anyone who wants to check a bunch of URLs for being dead ☠️  or not, SEO experts? 

With time _AreYouOk_ can evolve to analyze URLs over a remote resource as well, send your ideas ✨️ through [**Discussions**](https://github.com/Bhupesh-V/areyouok/discussions)

## Installation

- Linux
  ```bash
  curl -LJO https://github.com/Bhupesh-V/areyouok/releases/latest/download/areyouok-linux-amd64
  # or
  wget -q https://github.com/Bhupesh-V/areyouok/releases/latest/download/areyouok-linux-amd64
  ```

- MacOS
  ```bash
  curl -LJO https://github.com/Bhupesh-V/areyouok/releases/latest/download/areyouok-darwin-amd64
  # or
  wget -q https://github.com/Bhupesh-V/areyouok/releases/latest/download/areyouok-darwin-amd64
  ```

- Windows
  ```bash
  curl -LJO https://github.com/Bhupesh-V/areyouok/releases/latest/download/areyouok-windows-amd64
  # or
  wget -q https://github.com/Bhupesh-V/areyouok/releases/latest/download/areyouok-windows-amd64
  ```

Check installation by running `areyouok -v`
```bash
$ mv areyouok-darwin-amd64 areyouok
$ areyouok -v
v1.0.0
```

Download builds for other architectures from [**releases**](https://github.com/Bhupesh-V/areyouok/releases/latest)


## Usage

AreYouOk provides 3 optional arguments followed by a directory path (default: current directory)
1. `-t` type of files to scan for links
2. `-i` list of files or directories to ignore for links (node_modules, .git)
3. `-r` type of report to generate

Some example usages:

- Analyze all HTML files for hyperlinks. The default directory path is set to current directory.
  ```bash
  areyouok -t=html Documents/my-directory
  ```
  There is not limitation on file type, to analyze json files use `-t=json`. The default type is set to `md` (markdown)

- To ignore certain directories or file-types use the ignore `-i` flag.
  ```bash
  areyouok -i=_layouts,.git,_site,README.md,build.py,USAGE.md
  ```

- By default AreYouOk outputs analyzed data directly into console. To generate a report, use the `-r` flag
  ```bash
  areyouok -i=_layouts,.git,_site,README.md,build.py,USAGE.md -r=html ~/Documents/til/
  ```
  Currently supported report formats are: `json`, `txt`, `html` & `github`.

  Different report types provide different levels of information which is briefly summarized below:
  1. **JSON** (report.json)<br>
     The JSON report can be used for other computational tasks as required (E.g emailing dead urls to yourself)
     ```json
     {
        "http://127.0.0.1:8000/": {
                "code": "",
                "message": "Get \"http://127.0.0.1:8000/\": dial tcp 127.0.0.1:8000: connect: connection refused",
                "response_time": ""
        },
        "http://freecodecamp.org": {
                "code": "200",
                "message": "OK",
                "response_time": "5.44s"
        },
        "http://ogp.me/": {
                "code": "200",
                "message": "OK",
                "response_time": "3.60s"
        },
        "http://prnbs.github.io/projects/regular-expression-parser/": {
                "code": "200",
                "message": "OK",
                "response_time": "0.25s"
        },
        "https://bhupeshv.me/30-Seconds-of-C++/": {
                "code": "404",
                "message": "Not Found",
                "response_time": "3.84s"
        },
        ...
     }
     ```

  2. **Text** (report.txt)<br>
     The text format just lists the URLs which were not successfully fetched. Useful if you just want dead urls.
     Text report also puts the no.of hyperlinks analyzed along with total files & total reponse time.

     ```
     74 URLs were analyzed across 31 files in 21.69s

     Following URLs are not OK:

     http://freecodecamp.org`
     http://127.0.0.1:8000/
     https://drive.google.com/uc?export=view&id=<INSERT-ID>`
     https://drive.google.com/file/d/
     https://drive.google.com/uc?export=view&id=$get_last
     https://github.com/codeclassroom/PlagCheck/blob/master/docs/docs.md
     https://bhupeshv.me/30-Seconds-of-C++/
     ```
     Note that the total time would vary according to your internet speed & website latency.

  3. **HTML** (report.html)<br>
     The html report is the most superior formats of all & can be used to have a visual representaion of analyzed links.<br>
     Below is demo of how this HTML report looks like, [**you can see it live**]()
     ![report-latest](https://user-images.githubusercontent.com/34342551/105046278-e80db380-5a8e-11eb-8371-124fae8b3d7f.png)

  4. **GitHub** (report.github)<br>
     The github report format is well suited if you are utilizing Github Actions. The format generated is largely HTML, compatible with github's commonmark markdown renderer.<br>
     Below is a demo of a Github Action which reports the analyzed URLs through github issues. [Here is a demo link](https://github.com/Bhupesh-V/til/issues/2)
     
     ![demo-action](https://user-images.githubusercontent.com/34342551/105579706-169cce80-5dae-11eb-8dd6-b51bf23e63ee.png)
     


## Development

#### Prerequisites

- [Go 1.16](https://golang.org/dl/#unstable)

1. Clone the repository.
   ```bash
   git https://github.com/Bhupesh-V/areyouok.git
   ```
2. Run tests.
   ```bash
   go test -v
   ```
3. Format & Lint the project.
   ```bash
   gofmt -w areyouok.go && golint areyouok.go
   ```

## 📝 Changelog

See the [CHANGELOG.md](CHANGELOG.md) file for details.

## ☺️ Show your support

Support me by giving a ⭐️ if this project helped you! or just [![Twitter URL](https://img.shields.io/twitter/url?style=social&url=https%3A%2F%2Fgithub.com%2FBhupesh-V%2Fareyouok%2F)](https://twitter.com/intent/tweet?url=https://github.com/Bhupesh-V/areyouok&text=areyouok%20via%20@bhupeshimself)

<a href="https://liberapay.com/bhupesh/donate">
  <img alt="Donate using Liberapay" src="https://liberapay.com/assets/widgets/donate.svg" width="100">
</a>
<a href="https://ko-fi.com/bhupesh">
  <img title="ko-fi/bhupesh" alt="Support on ko-fi" src="https://user-images.githubusercontent.com/34342551/88784787-12507980-d1ae-11ea-82fe-f55753340168.png" width="185">
</a>


## 📝 License

Copyright © 2020 [Bhupesh Varshney](https://github.com/Bhupesh-V).<br />
This project is [MIT](https://github.com/Bhupesh-V/areyouok/blob/master/LICENSE) licensed.

## 👋 Contributing

Please read the [CONTRIBUTING](CONTRIBUTING.md) file for the process of submitting pull requests to us.
