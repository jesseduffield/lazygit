# Profiling Lazygit

If you want to investigate what's contributing to CPU or memory usage, start
lazygit with the `-profile` command line flag. This tells it to start an
integrated web server that listens for profiling requests.

## Save profile data

### CPU

While lazygit is running with the `-profile` flag, perform a CPU profile and
save it to a file by running this command in another terminal window:

```sh
curl -o cpu.out http://127.0.0.1:6060/debug/pprof/profile
```

By default, it profiles for 30 seconds. To change the duration, use

```sh
curl -o cpu.out 'http://127.0.0.1:6060/debug/pprof/profile?seconds=60'
```

### Memory

To save a heap profile (containing information about all memory allocated so
far since startup), use

```sh
curl -o mem.out http://127.0.0.1:6060/debug/pprof/heap
```

Sometimes it can be useful to get a delta log, i.e. to see how memory usage
developed from one point in time to another. For that, use

```sh
curl -o mem.out 'http://127.0.0.1:6060/debug/pprof/heap?seconds=20'
```

This will log the memory usage difference between now and 20 seconds later, so
it gives you 20 seconds to perform the action in lazygit that you are interested
in measuring.

## View profile data

To display the profile data, you can either use speedscope.app, or the pprof
tool that comes with go. I prefer the former because it has a nicer UI and is a
little more powerful; however, I have seen cases where it wasn't able to load a
profile for some reason, in which case it's good to have the pprof tool as a
fallback.

### Speedscope.app

Go to https://www.speedscope.app/ in your browser, and drag the saved profile
onto the browser window. Refer to [the
documentation](https://github.com/jlfwong/speedscope?tab=readme-ov-file#usage)
for how to navigate the data.

### Pprof tool

To view a profile that you saved as `cpu.out`, use

```sh
go tool pprof -http=:8080 cpu.out
```

By default this shows the graph view, which I don't find very useful myself.
Choose "Flame Graph" from the View menu to show a much more useful
representation of the data.
