# JS combiner CLI 🚀

A lightweight, zero-dependency Go utility to recursively walk a directory, concatenate JavaScript files, and provide basic minification. Perfect for projects using raw HTML5/CSS that want a simple build step without the overhead of Node.js or Webpack.

## Features
* **Recursive Walking:** Automatically finds all `.js` files in subdirectories.
* **Smart Concatenation:** Adds headers and safety newlines between files.
* **Basic Minification:** Strips comments and collapses whitespace for production.
* **No Dependencies:** Built entirely with the Go Standard Library.

## Installation

Ensure you have Go installed and your `GOBIN` is in your system path.

```bash
go install github.com/roborobs1023/js_combiner@latest
```
## Usage

```bash
js_combiner -src="./src/js" -out="./assets/js/bundle.js" -minify
```
