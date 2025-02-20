package expressions

import (
	"os"
	"strings"
)

// ColorEnabled controls whether we actually print ANSI escapes.
// Set via ENABLE_COLORS=1 or =true to enable color.
var ColorEnabled = initColorEnabled()

// ANSI reset code, applied after coloring a token.
var ColorReset = "\033[0m"

// The four palette names you can choose from:
const (
	PaletteMild      = "mild"
	PaletteVivid     = "vivid"
	PaletteDracula   = "dracula"
	PaletteSolarized = "solarized"
)

// PunctuationColor, StringColor, etc. are updated at init() time
// based on your chosen palette. They are used throughout the DSL .String() methods.
var (
	PunctuationColor string
	StringColor      string
	NumberColor      string
	OperatorColor    string
	BoolNullColor    string
	IdentifierColor  string
	LibraryColor     string
	FunctionColor    string
	ContextColor     string
)

// initColorEnabled checks if ENABLE_COLORS is "1" or "true" (case-insensitive).
func initColorEnabled() bool {
	val := strings.ToLower(os.Getenv("ENABLE_COLORS"))
	return val == "1" || val == "true"
}

func init() {
	// Determine which palette to load from env var COLOR_PALETTE.
	// If not set or unrecognized, we use "default".
	paletteName := strings.ToLower(os.Getenv("COLOR_PALETTE"))
	switch paletteName {
	case PaletteVivid:
		ApplyVividPalette()
	case PaletteDracula:
		ApplyDraculaPalette()
	case PaletteSolarized:
		ApplySolarizedPalette()
	case PaletteMild:
		ApplyMildPalette()
	default:
		ApplySolarizedPalette() // fallback
	}
}

// applyDefaultPalette sets a "mild" or "neutral" palette.
func ApplyMildPalette() {
	// One Dark–inspired
	PunctuationColor = "\033[38;2;92;99;112m"  // #5C6370 (comment-ish gray)
	StringColor = "\033[38;2;152;195;121m"     // #98C379 (green)
	NumberColor = "\033[38;2;209;154;102m"     // #D19A66 (orange-brown)
	OperatorColor = "\033[38;2;198;120;221m"   // #C678DD (purple)
	BoolNullColor = "\033[38;2;86;182;194m"    // #56B6C2 (cyan)
	IdentifierColor = "\033[38;2;229;192;123m" // #E5C07B (yellow-gold)
	LibraryColor = "\033[38;2;171;178;191m"    // #ABB2BF (soft foreground)
	FunctionColor = "\033[38;2;97;175;239m"    // #61AFEF (blue)
	ContextColor = "\033[38;2;224;108;117m"    // #E06C75 (reddish-pink)
}

// ApplyVividPalette sets a more saturated, bold color set.
func ApplyVividPalette() {
	// Extra bright / neon
	PunctuationColor = "\033[38;2;255;128;0m" // bright orange
	StringColor = "\033[38;2;255;85;85m"      // bright red/pink
	NumberColor = "\033[38;2;0;255;0m"        // lime green
	OperatorColor = "\033[38;2;255;0;255m"    // hot magenta
	BoolNullColor = "\033[38;2;0;170;255m"    // bright cyan
	IdentifierColor = "\033[38;2;255;215;0m"  // gold
	LibraryColor = "\033[38;2;255;160;0m"     // bright orange
	FunctionColor = "\033[38;2;85;85;255m"    // vivid blue
	ContextColor = "\033[38;2;255;20;147m"    // deep pink
}

// ApplyDraculaPalette sets colors inspired by the Dracula theme.
func ApplyDraculaPalette() {
	// Official Dracula color references:
	// https://draculatheme.com/contribute
	// background: #282a36  foreground: #f8f8f2
	// comment: #6272a4, cyan: #8be9fd, green: #50fa7b, orange: #ffb86c,
	// pink: #ff79c6, purple: #bd93f9, red: #ff5555, yellow: #f1fa8c

	PunctuationColor = "\033[38;2;98;114;164m" // #6272a4 (used often for comments/punctuation)
	StringColor = "\033[38;2;241;250;140m"     // #f1fa8c (yellow)
	NumberColor = "\033[38;2;189;147;249m"     // #bd93f9 (purple)
	OperatorColor = "\033[38;2;255;121;198m"   // #ff79c6 (pink)
	BoolNullColor = "\033[38;2;139;233;253m"   // #8be9fd (cyan)
	IdentifierColor = "\033[38;2;80;250;123m"  // #50fa7b (green)
	LibraryColor = "\033[38;2;255;184;108m"    // #ffb86c (orange)
	FunctionColor = "\033[38;2;255;85;85m"     // #ff5555 (red)
	ContextColor = "\033[38;2;248;248;242m"    // #f8f8f2 (foreground-ish)
}

// ApplySolarizedPalette sets colors inspired by the Solarized Dark theme.
func ApplySolarizedPalette() {
	// Official Solarized Dark references:
	// https://github.com/altercation/vim-colors-solarized
	//
	// base03: #002b36
	// base02: #073642
	// base01: #586e75
	// base00: #657b83
	// base0:  #839496
	// base1:  #93a1a1
	// base2:  #eee8d5
	// base3:  #fdf6e3
	// yellow: #b58900
	// orange: #cb4b16
	// red:    #dc322f
	// magenta:#d33682
	// violet: #6c71c4
	// blue:   #268bd2
	// cyan:   #2aa198
	// green:  #859900

	PunctuationColor = "\033[38;2;88;110;117m" // #586e75 (base01)
	StringColor = "\033[38;2;42;161;152m"      // #2aa198 (cyan)
	NumberColor = "\033[38;2;133;153;0m"       // #859900 (green)
	OperatorColor = "\033[38;2;108;113;196m"   // #6c71c4 (violet)
	BoolNullColor = "\033[38;2;38;139;210m"    // #268bd2 (blue)
	IdentifierColor = "\033[38;2;181;137;0m"   // #b58900 (yellow)
	LibraryColor = "\033[38;2;147;161;161m"    // #93a1a1 (base1)
	FunctionColor = "\033[38;2;211;54;130m"    // #d33682 (magenta)
	ContextColor = "\033[38;2;203;75;22m"      // #cb4b16 (orange)
}
