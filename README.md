# Hyde
## Guaranteed to be written by Drunken Monkeys 

### For when you want a static website, not a Ph. D in Jekyll

A static website generator in Go, like Jekyll, without the nonsensical complications. Point it at an input directory and an output directory and it will make a website. 

No Ruby, No Gems, No Coffeescript. Just Go.

`I solemnly swear all commits to this repo were made while inebriated and dressed as a monkey.`

![XKCD](http://imgs.xkcd.com/comics/ballmer_peak.png)

# Basic Idea

You create a series of files with markdown in them, the program recursively compiles them to HTML in the target directory in the same structure. Files named with an `_` are special and not compiled. ANY file not named with an `_` will end up in the target directory, and ALL files are run through the templating engine. Files named with `.md` are wrapped in the `layout`, then run through the templating engine, then the Markdown engine.

If you DON'T want to use layouts (or just for a directory), then create a `_layout.html` file with just the action tag `{{.Content}}` and then handle header and footer, or whatever on a per file basis.

You can setup variables on a per file basis by placing action code in a shebang line.

`#! {{$title := "This is the page title"}}`

Notice that it is '#!{SPACE}' not simply '#!'. You can set any variables you want, like meta_descriptions and so on, or select specific CSS. The system doesn't require ANY of this, it's all based on how you do the layout. You can echo out the default layout with the command:

`./hyde -layout`

# Special files

### _layout.html

Will be picked up by the compiler and overrides the default compiled in. Only requirement is the special action `{{.Content}}`.

Each directory can have its own _layout.html

### xxxx--some-file-name.md

Naming a file like this (replacing xxxx with some alphanumeric id) allows you to refer to that file anywere with functions like `{{link_to "xxxx"}}` to get a full absolute URL (minus the http). 

Example: `[Home](https://{{link_to "index.md"}})`

### Partials (or any file really)

Any file has at least the id of its name. So if you create a file: _footer.html, then you can find it with the `partial` command, i.e. `{{partial "_footer.html"}}`

The id is literally the path to that file (always absolute path), so to get to a tidbit of code in `_src/some/sub/dir/_partial.md` you would use `{{partial "some/sub/dir/_partial.md"}}`. NOTE the missing `_src`, the root is omitted in all ids, because the only thing that actually changes is the root, from _src to _dest.

All files are run through the template engine. If you name a file with `.md` suffix, it will be run through the template engine and then the Markdown engine.

### Want permalinks?

Don't move things.

### Why the f***k didn't you use Hugo

Too big, too touchy feely. Have you seen their repo? It's scurry. Also it's good practice as a programmer to build your own systems as often as you can. It's how you learn to program. Otherwise, you're just copying and pasting all the time from the work of others. That makes you a `coder`. 
