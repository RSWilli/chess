package player

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rswilli/chess/internal/chess"
	"github.com/rswilli/chess/internal/uci"
)

func NewEnginesRegistry(folder string) (EnginesRegistry, error) {
	fi, err := os.Stat(folder)

	if err != nil {
		return EnginesRegistry{}, err
	}

	if !fi.IsDir() {
		return EnginesRegistry{}, err
	}

	return EnginesRegistry{
		enginesFolder: folder,
	}, nil
}

type EnginesRegistry struct {
	enginesFolder string
}

func (r EnginesRegistry) AvailableEngines() ([]string, error) {
	cwd, err := os.Getwd()

	if err != nil {
		return nil, fmt.Errorf("could not get cwd: %w", err)
	}

	entries, err := os.ReadDir(filepath.Join(cwd, r.enginesFolder))

	if err != nil {
		return nil, fmt.Errorf("could not read engines folder: %w", err)
	}

	engines := make([]string, 0, len(entries)+3)

	engines = append(engines, EngineHuman, EngineStockfish, EngineLocal)

	for _, f := range entries {
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}

		engines = append(engines, f.Name())
	}

	return engines, nil
}

const EngineStockfish = "stockfish"
const EngineHuman = "human"
const EngineLocal = "local"

func (r EnginesRegistry) NewEngine(name string) (uci.Engine, error) {
	switch name {
	case EngineLocal:
		return chess.NewEngine(), nil
	case EngineStockfish:
		return uci.NewStockfish()
	case EngineHuman:
		return NewHuman(), nil
	}

	cwd, err := os.Getwd()

	if err != nil {
		return nil, fmt.Errorf("could not get cwd: %w", err)
	}

	return uci.NewClient(filepath.Join(cwd, r.enginesFolder, name))
}
