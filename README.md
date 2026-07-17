<p align="center">
  <h1 align="center"><b>promptkit</b></h1>
  <p align="center"><i>Interactive command line prompts with style!</i></p>
  <p align="center">
    <a href="https://github.com/erikgeiser/promptkit/releases/latest"><img alt="Release" src="https://img.shields.io/github/release/erikgeiser/promptkit.svg?style=for-the-badge"></a>
    <a href="https://pkg.go.dev/github.com/erikgeiser/promptkit"><img alt="Go Doc" src="https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge"></a>
    <a href="https://github.com/erikgeiser/promptkit/actions?workflow=Build"><img alt="GitHub Action: Build" src="https://img.shields.io/github/actions/workflow/status/erikgeiser/promptkit/build.yml?branch=main&style=for-the-badge"></a>
    </br>
    <a href="https://github.com/erikgeiser/promptkit/actions?workflow=Check"><img alt="GitHub Action: Check" src="https://img.shields.io/github/actions/workflow/status/erikgeiser/promptkit/check.yml?branch=main&style=for-the-badge"></a>
    <a href="/LICENSE.md"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge"></a>
  </p>
</p>

Promptkit is a collection of common command line prompts for interactive
programs. Each prompts comes with sensible defaults, re-mappable key bindings
and many opportunities for heavy customization.

---

**Disclaimers:**

- The API of library is not yet stable. Expect significant changes in minor
  versions before `v1.0.0`.

---

## Selection

Selection with filter and pagination support: [Example Code](https://github.com/erikgeiser/promptkit/blob/main/examples/selection/main.go)

![Selection Prompt](.assets/selection.gif)

---

The selection prompt is highly customizable and also works well with custom
types which for example enables filtering only on custom fields: [Example Code](https://github.com/erikgeiser/promptkit/blob/main/examples/selection_custom/main.go)

![Custom Selection Prompt](.assets/selection_custom.gif)

The selection module also contains a multi-selection prompt (this is still a preview, use `go get go get github.com/erikgeiser/promptkit@main`): [Example Code](https://github.com/erikgeiser/promptkit/blob/main/examples/multi_selection/main.go)

![Multi-Selection Prompt](.assets/multi_selection.gif)

The selection module also contains a multi-selection prompt (this is still a preview, use `go get go get github.com/erikgeiser/promptkit@main`): [Example Code](https://github.com/erikgeiser/promptkit/blob/main/examples/multi_selection_custom/main.go)

![Custom Multi-Selection Prompt](.assets/multi_selection_custom.gif)

---

## Text Input

A text input that supports editable default values: [Example Code](https://github.com/erikgeiser/promptkit/blob/main/examples/textinput/main.go)

![Text Input Prompt](.assets/textinput.gif)

---

Custom validation is also supported: [Example Code](https://github.com/erikgeiser/promptkit/blob/main/examples/textinput_custom/main.go)

![Custom Text Input Prompt](.assets/textinput_custom.gif)

---

The text input can also be used in hidden mode for password prompts: [Example Code](https://github.com/erikgeiser/promptkit/blob/main/examples/textinput_hidden/main.go)

![Hidden Text Input Prompt](.assets/textinput_hidden.gif)

The text input prompt also supports auto-completion: [Example Code](https://github.com/erikgeiser/promptkit/blob/main/examples/textinput_completion/main.go)

![Text Input Prompt With Auto-Completion](.assets/textinput_completion.gif)

For longer inputs, there is also a textara (this is still a preview, use `go get go get github.com/erikgeiser/promptkit@main`): [Example Code](https://github.com/erikgeiser/promptkit/blob/main/examples/textarea/main.go)

![Text Area Prompt](.assets/textarea.gif)

It can also be customized (this is still a preview, use `go get go get github.com/erikgeiser/promptkit@main`): [Example Code](https://github.com/erikgeiser/promptkit/blob/main/examples/textarea_custom/main.go)

![Custom Text Area Prompt](.assets/textarea_custom.gif)

---

## Confirmation Prompt

A confirmation prompt for binary questions: [Example Code](https://github.com/erikgeiser/promptkit/blob/main/examples/confirmation/main.go)

![Confirmation Prompt](.assets/confirmation.gif)

The confirmation prompt can be customized as well: [Example Code](https://github.com/erikgeiser/promptkit/blob/main/examples/confirmation_custom/main.go):

![Custom Confirmation Prompt](.assets/confirmation_custom.gif)

## Widget

The prompts in this library can also be used as [bubbletea](https://github.com/charmbracelet/bubbletea) widgets: [Example Code](https://github.com/erikgeiser/promptkit/blob/main/examples/bubbletea_widget/main.go)

## Acknowledgements

This library is built on top of many great libraries, especially the following:

- https://github.com/charmbracelet/bubbletea
- https://github.com/charmbracelet/bubbles
- https://github.com/muesli/termenv
- https://github.com/muesli/reflow
