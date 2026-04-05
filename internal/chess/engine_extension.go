package chess

// Render implements [uci.Extended].
func (e *Engine) Render() string {
	return e.pos.ASCIIArt()
}
