## dibbler

`dibbler` is the node scanner for Sinclair. It performs one of two tasks:

1. It scans the files passed in as an argument, returning their mtime
   and the path in a JSON structure.

2. It scans and loads the files passed in as an argument, returning
   their metadata and markdown-rendered content as a JSON file. This
   is used to populate a new node structure.

The default behaviour is to load a node; the alternate behaviour may
be specified with the `-mod` flag.

### Usage

    dibbler [-mod] [files...]

### Output

`dibbler` will attempt to spit out JSON containing the relevant
information.

#### Node-loading

The returned JSON is in the form:

    {
        "message": "",
        "nodes": [
            {
                "body": "<p>And now we speak of things to come&hellip;</p>\n",
                "mtime": 1393914468,
                "path": "testpost.md",
                "slug": "bespoke-blogging-engines",
                "static": false,
                "tags": [
                    "bespoke",
                    "blogging",
                    "lisp",
                    "hacking"
                ]
            }
        ],
        "success": true
    }

Notice that if success is true, there will be no message; the message
field is used only to convery error messages. If success is false, the
error message will be found in message.

A node is an array of objects containing:

* title: the node's title as used for HTML rendering
* static: true if the node is a :page, false if the node is a post
* tags: a list of metadata tags for the node
* mtime: last modified time
* path: the path to the node (which is the same path that was passed
  to `dibbler`)
* slug: the url slug to use; for posts, this will end up as
  /yyyy/mm/dd/slug; for pages, it will end up as /slug. If not
  present, the basename of the path relative to the source directory
  path will be used.
* body: the rendered HTML fragment of the file.

#### Node mtime scanning

    {
        "message": "",
        "results": [
            {
                "mtime": 1393914468,
                "path": "testpost.md"
            }
        ],
        "success": true
    }

The same success and message fields are used in this container;
however, the results field contains a list of objects storing the path
to the node and the last modified time. This can be used to build a
list of nodes that should be updated.

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

