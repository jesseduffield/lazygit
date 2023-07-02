# Searching/Filtering

## View searching/filtering

Depending on the currently focused view, hitting '/' will bring up a filter or search prompt. When filtering, the contents of the view will be filtered down to only those lines which match the query string. When searching, the contents of the view are not filtered, but matching lines are highlighted and you can iterate through matches with `n`/`N`.

We intend to support filtering for the files view soon, but at the moment it uses searching. We intend to continue using search for the commits view because you typically care about the commits that come before/after a matching commit.

If you would like both filtering and searching to be enabled on a given view, please raise an issue for this.

## Filtering files by status

You can filter the files view to only show staged/unstaged files by pressing `<c-b>` in the files view.

## Filtering commits by file path

You can filter the commits view to only show commits which contain changes to a given file path.

You can do this in a couple of ways:
1) Start lazygit with the -f flag e.g. `lazygit -f my/path`
2) From within lazygit, press `<c-s>` and then enter the path of the file you want to filter by
