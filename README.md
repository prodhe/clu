# clu

## Introduction

Command line lookup utility that returns the Mozilla Developer Network article for a given keyword.

It fetches results from the HTML and CSS documentation and selects the one which answers first and non-empty. Then it transforms the relevant part from HTML to plain text and prints to standard output.

## Credits

As of now it uses [html2text](https://github.com/jaytaylor/html2text) for the plain text output.

## License

[MIT](LICENSE.txt).