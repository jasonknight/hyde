# Hyde

### For when you want a static website, not a Ph. D in Jekyll

A static website generator in Go, like Jekyll, without the nonsensical complications. Point it at an input directory and an output directory and it will make a website. 

No Ruby, No Gems, No Coffeescript. Just Go.

# Basic Idea

You create a series of files with markdown in them, the program recursively compiles them to HTML in the target directory in the same structure. Files named with an _ are special and not compiled.

## Special files

### _layout.md

Will be picked up by the compiler and overrides the default compiled in. Only requirement is the special action {{.Content}}

### Want permalinks?

Don't move things.

 
