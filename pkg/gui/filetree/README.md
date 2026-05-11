## FileTree Package

This package handles the representation of file trees. There are two ways to render files: one is to render them flat, so something like this:

```
dir1/file1
dir1/file2
file3
```

And the other is to render them as a tree

```
dir1/
  file1
  file2
file3
```

Internally we represent each of the above as a tree, but with the flat approach there's just a single root node and every path is a direct child of that root. Viewing in 'tree' mode (as opposed to 'flat' mode) allows for collapsing and expanding directories, and lets you perform actions on directories e.g. staging a whole directory. But it takes up more vertical space and sometimes you just want to have a flat view where you can go flick through your files one by one to see the diff.

This package is not concerned about rendering the tree: only representing its internal state.
