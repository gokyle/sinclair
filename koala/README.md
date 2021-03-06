## koala

`koala` is a small utility to crawl through a filesystem tree and
return the filenames contained. By default, it only returns files (not
the directories containing those files) and skips dot files.

### Usage

    koala [-ads] [-o style] [directories...]
            -a        show hidden (dot) files
            -d        show directories
            -o style  use style for output
            -s        don't strip the root directory from paths (i.e. return
                      absolute paths)


If no directories are specified, koala runs in the current directory.

### Styles

There are a few styles provided; the default is "space".

* lisp: outputs a list of pathnames
* space: outputs a space-separated list of file names
* comma: outputs a comma-separated list of file names
* null: outputs a NUL-separated list of file names

#### Examples:

    $ koala -o lisp
    (#P"README.md" #P"koala.go")
     
    $ koala -o space
    README.md koala.go
     
    $ koala -o comma
    README.md,koala.go
     
    $ koala -o null
    README.md koala.go


### License

> Copyright (c) 2014 Kyle Isom <kyle@tyrfingr.is>
> 
> Permission to use, copy, modify, and distribute this software for any
> purpose with or without fee is hereby granted, provided that the above 
> copyright notice and this permission notice appear in all copies.
> 
> THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
> WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
> MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
> ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
> WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
> ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
> OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE. 
