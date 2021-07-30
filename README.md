# Promptkit

[![Go Reference](https://pkg.go.dev/badge/github.com/erikgeiser/promptkit.svg)](https://pkg.go.dev/github.com/erikgeiser/promptkit)

Promptkit is a collection of common command line prompts for interactive
programs. Each prompts comes with sensible defaults, re-mappable key bindings
and many opportunities for heavy customization.

## Selection

Selection with filter and pagination support:

[![asciicast](https://asciinema.org/a/4ZK5HZ2uJm8NtC0rs8rnqxUwS.svg)](https://asciinema.org/a/4ZK5HZ2uJm8NtC0rs8rnqxUwS)

https://github.com/erikgeiser/promptkit/blob/ea17c82a1ba5299a2eb2b00bc1b1b5baf4e52a5e/examples/selection/main.go#L11-L15

The selection prompt is highly customizable and also works well with custom
types which for example enables filtering only on custom fields:

[![asciicast](https://asciinema.org/a/T9SG8WwP683dZxRdh1cAD6Deu.svg)](https://asciinema.org/a/T9SG8WwP683dZxRdh1cAD6Deu)

https://github.com/erikgeiser/promptkit/blob/ea17c82a1ba5299a2eb2b00bc1b1b5baf4e52a5e/examples/selection_custom/main.go#L55-L70

## Text Input

A text input that supports editable default values:

[![asciicast](https://asciinema.org/a/tJCUnnKxoXivvSf0gSkZfAjdn.svg)](https://asciinema.org/a/tJCUnnKxoXivvSf0gSkZfAjdn)

https://github.com/erikgeiser/promptkit/blob/f29e12dd8eb290771e9652a0eda6cac0e3895976/examples/textinput/main.go#L11-L16

Custom validation is also supported:

[![asciicast](https://asciinema.org/a/LNsZi7yrk7SvrcYCLROnUk7Of.svg)](https://asciinema.org/a/LNsZi7yrk7SvrcYCLROnUk7Of)

https://github.com/erikgeiser/promptkit/blob/f29e12dd8eb290771e9652a0eda6cac0e3895976/examples/textinput_custom/main.go#L18-L24

The text input can also be used in hidden mode for password prompts:

[![asciicast](https://asciinema.org/a/HcqfFKCIPSBClYYjqJdDqJ35z.svg)](https://asciinema.org/a/HcqfFKCIPSBClYYjqJdDqJ35z)

https://github.com/erikgeiser/promptkit/blob/f29e12dd8eb290771e9652a0eda6cac0e3895976/examples/textinput_hidden/main.go#L11-L17
