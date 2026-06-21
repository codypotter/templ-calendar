# templ-calendar

Calendar components for [a-h/templ](https://github.com/a-h/templ). Styled with Tailwind CSS — so BYO tailwind build.

**[Examples and full documentation →](https://codypotter.github.io/templ-calendar/)**

## Getting started

Add all components at once:

```sh
go run github.com/codypotter/templ-calendar/cmd/templ-calendar@latest add all ./components/calendar
```

Or pick individual ones:

```sh
go run github.com/codypotter/templ-calendar/cmd/templ-calendar@latest add calendar ./components/calendar
go run github.com/codypotter/templ-calendar/cmd/templ-calendar@latest add navigator ./components/calendar
go run github.com/codypotter/templ-calendar/cmd/templ-calendar@latest add jumper ./components/calendar
```

Then run:

```sh
templ generate
```

## Development

```sh
task dev
```

Requires [Task](https://taskfile.dev), [templ](https://templ.guide), and the [Tailwind CSS standalone CLI](https://tailwindcss.com/blog/standalone-cli).
