/*
 * NSSF Logger
 */

package flog

// Escape character
const Escape string = "\033"

// Text style
const (
    NoEffect int = iota
    Bold
    Faint
    Highlight
    Underline
)

// Text color
const (
    Black int = iota + 30
    Red
    Green
    Yellow
    Blue
    Magenta
    Cyan
    White
)

// Hi-intensity text color
const (
    HiBlack int = iota + 90
    HiRed
    HiGreen
    HiYellow
    HiBlue
    HiMagenta
    HiCyan
    HiWhite
)
