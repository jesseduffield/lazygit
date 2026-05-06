# Searching/Filtering

## View searching/filtering

Depending on the currently focused view, hitting '/' will bring up a filter or search prompt. When filtering, the contents of the view will be filtered down to only those lines which match the query string. When searching, the contents of the view are not filtered, but matching lines are highlighted and you can iterate through matches with `n`/`N`.

### Regular expression filters

You can filter with a [Go regular expression](https://go.dev/s/re2syntax) in two ways:

- Set `gui.filterMode` to `regexp` in your config so every filter string is treated as a regexp.
- Or keep `substring` / `fuzzy` and prefix a single filter with `gui.regexpFilterPrefix` (default `re:`, for example `re:^main` to match branch names that start with `main`). You can change the prefix in config if `re:` collides with how you name branches or paths.

If the pattern has no uppercase letters, matching is case-insensitive (the same rule as plain substring filters). Invalid regexps match nothing. In regexp mode, `.` and other metacharacters are active, so file paths like `foo.go` need `foo\.go` to match a literal dot. Unlike substring mode, regexp mode uses one pattern for the whole line (whitespace inside the pattern is not split into multiple AND terms).

We intend to continue using search for the commits view because you typically care about the commits that come before/after a matching commit.

If you would like both filtering and searching to be enabled on a given view, please raise an issue for this.

## Filtering files by status

You can filter the files view to only show staged/unstaged files by pressing `<c-b>` in the files view.

## Filtering commits by file path

You can filter the commits view to only show commits which contain changes to a given file path.

You can do this in a couple of ways:
1) Start lazygit with the -f flag e.g. `lazygit -f my/path`
2) From within lazygit, press `<c-s>` and then enter the path of the file you want to filter by
