package filter

// Filter is the common interface implemented by all record matchers
// (TimeRange, FieldFilter, MultiFilter, etc.).
// It is declared here to avoid duplicate-declaration conflicts when
// the interface is referenced across multiple files in the package.
// (If your build already defines Filter elsewhere, remove this file.)
