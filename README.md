# templ-calendar

Calendar components for [a-h/templ](https://github.com/a-h/templ). Styled with Tailwind CSS — so BYO tailwind build.

**[Examples and full documentation →](https://codypotter.github.io/templ-calendar/)**

## Install

```sh
go install github.com/codypotter/templ-calendar/cmd/templ-calendar@latest
```

## Getting started

```sh
templ-calendar add calendar ./components/calendar
templ-calendar add navigator ./components/calendar
templ-calendar add jumper ./components/calendar

templ generate
```

## Development

```sh
task dev
```

Requires [Task](https://taskfile.dev), [templ](https://templ.guide), and the [Tailwind CSS standalone CLI](https://tailwindcss.com/blog/standalone-cli).
