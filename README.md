# Chess Server and engine

Just me trying to implement a chess server and chess engine.

## Engines built

All built engines are tagged and measured against the next version.

1. `randotron`: A chess engine only playing random moves.

Using the [`compile_engines.sh`](./compile_engines.sh) script you can compile all engine versions from any
git revision.

## Tests

Chess rules are tested against [stockfish](https://stockfishchess.org/) and can be verified using the engine perft test

```bash
go test ./internal/chess -v -run=^TestPerft
```

## Performance

Benchmark the perft command:

```bash
go test ./internal/chess -run=^$ -v -bench=BenchmarkPerft -cpuprofile=profile.prof
```

## Attributions

* chess piece svgs taken from https://commons.wikimedia.org/wiki/Category:SVG_chess_pieces
* Bootstrap Icons:
  * Trophy: https://icons.getbootstrap.com/icons/trophy-fill/
  * Flag: https://icons.getbootstrap.com/icons/flag/
