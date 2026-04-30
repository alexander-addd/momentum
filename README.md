# Momentum

Momentum is a small time tracking CLI app written in Go.

It is inspired by the useful parts of Toggl Track: start a timer, stop it when the work is done, and keep a clean local history of where time went.

The first version is intentionally simple:

```sh
momentum start "Write project README" --project momentum --tag planning
momentum status
momentum stop
momentum today
```

The project is local-first, personal-use friendly, and meant to be pleasant to use from the terminal.
