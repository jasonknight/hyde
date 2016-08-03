# Hyde
## Guaranteed to be written by Drunken Monkeys 

### For when you want a static website, not a Ph. D in Jekyll

A static website generator in Go, like Jekyll, without the nonsensical complications. Point it at an input directory and an output directory and it will make a website. 

No Ruby, No Gems, No Coffeescript. Just Go.

`I solemnly swear all commits to this repo were made while inebriated and dressed as a monkey.`

# Basic Idea

You create a series of files with markdown in them, the program recursively compiles them to HTML in the target directory in the same structure. Files named with an `_` are special and not compiled.

## Special files

### _layout.html

Will be picked up by the compiler and overrides the default compiled in. Only requirement is the special action `{{.Content}}`.

Each directory can have its own _layout.html

### xxxx--some-file-name.md

Naming a file like this (replacing xxxx with some alphanumeric id) allows you to refer to that file anywere with functions like `{{link_to "xxxx"}}` to get a full absolute URL.

### Partials (or any file really)

Any file has at least the id of its name. So if you create a file: _footer.html, then you can find it with the `partial` command, i.e. `{{partial "_footer.html"}}`

### Want permalinks?

Don't move things.

### Why the f***k didn't you use Hugo

Too big, too touchy feely. Have you seen their repo? It's scurry.
