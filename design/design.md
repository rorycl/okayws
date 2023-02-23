# Design document

Revision A  
23 February 2023

This is a very terse design document for a small webserver code named
`okayws`, intended to be an "okay web server" that makes it easy for
clubs, teams or those running a small newsletter to run a website.
(I like the name `mehws` too...).

![okayws icon](okayws.png)

I've had experience of running a company
[Django](https://www.djangoproject.com/) website without a database, and
of helping maintain a static document site built using
[Sphinx](https://www.sphinx-doc.org/en/master/). Based on those I think
a small web server that has web pages described simply with markdown
files would work well for most uses. The static website builder
[hugo](https://gohugo.io/) is great, but `okayws` can be simpler and
more opinionated, to make it easier to get a site up-and-running, but
also run the site in production. `okayws` will be a single go binary.

Apart from content, though, a user will also have to find, modify or
create a css style sheet.

Markdown is a very simple way of defining web pages using simple text
formatting instructions. For example text after a \# describes a level 1
header, while text between \*\* symbols will be shown in **bold**. See
the [github markdown](https://docs.github.com/en/get-started/writing-on-github/getting-started-with-writing-and-formatting-on-github/basic-writing-and-formatting-syntax)
page for more information about this type of markdown.

## Operation

Websites could be structured something like this:

```
    website
    ├── home.md
    ├── pelagic
    │   ├── herring.md
    │   └── mackerel.md
    └── scarabs
        ├── home.md
        ├── dung_beetles.md
        ├── Goliath Beetles.html
        └── june-beetles.md
```

This would generate the following urls and associated resources:

```
url               : resource
------------------------------------------------
/                 : home.md
/pelagic/herring  : /pelagic/herring.md
/pelagic/mackarel : /pelagic/mackarel.md
/scarabs          : /scarabs/home.md
/dung-beetles     : /scarabs/dung_beetles.md
/goliath-beetles  : /scarabs/Goliath Beetles.md
/june-beetles     : /scarabs/june-beetles.md
```

The special name `home.md` provides a home page for the directory. In
the example above the url `/scarabs` will show the content in
`/scarabs/home.md` while `/pelagic` has no `home.md' resource so will
result in a "page not found". Note how resources with spaces or underbar
characters have renamed urls, with the file extension removed.

The content would be supported by media in a `media` directory and two
templates in a `templates` directory, `home.html` and `inner.html` for
wrapping the website contents. As a result the top level directories for
the webserver will be

```
.
├── website
├── templates
└── media
```

The `okayws` server can be run in production or development mode. The
default will be development mode. If it is run in development mode and
no `website`, `templates` or `media` directories exist, these will be
made, together with some default content.

```
# run in development mode
./okayws

# run in production mode
./okayws --production <project_directory>
```

In development mode, changes on the filesystem, such as the creation of
a new file or saving a change, will automatically reload the content of
the server to help review the rendered changes.

In production mode content will be read on startup (only) and html
gzipped for high-performance.

## Tags

A very limited number of go html template tags can be used in the
files in the `templates` directory. These tags are very simple and
have html.

```
{{ .Content }}
```

The `.Content` tag points to the
data from the rendered markdown. For example for the url `/scarabs` it
would contain the html from rendering `/scarabs/home.md`.

```
{{ .Sections }}
```

The `.Sections` tag would output the list of different sections, in the
example `Scarabs` and `Pelagic` (the latter would point to
`/pelagic/herring` as the `pelagic` section has no home page). This is
output as an unordered list.
