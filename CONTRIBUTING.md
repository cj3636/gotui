# Contributing to GoTUI

Thank you for your interest in contributing to GoTUI! This document provides guidelines for contributing to the project.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-username/gotui.git`
3. Create a new branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Test your changes: `go build && ./gotui`
6. Commit your changes: `git commit -m "Add your feature"`
7. Push to your fork: `git push origin feature/your-feature-name`
8. Open a Pull Request

## Development Guidelines

### Code Style

- Follow standard Go formatting: `gofmt` and `go vet`
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions small and focused

### Widget Development

When creating a new widget:

1. Create a new file in `internal/widgets/`
2. Implement the `Widget` interface:
   ```go
   type Widget interface {
       Init() tea.Cmd
       Update(msg tea.Msg) (Widget, tea.Cmd)
       View() string
       SetSize(width, height int)
       Title() string
   }
   ```
3. Embed `BaseWidget` for common functionality
4. Use appropriate emojis in the widget title
5. Handle errors gracefully
6. Add configuration options to `config.yaml` if needed

### Testing

- Test your widget with different terminal sizes
- Verify that resizing works correctly
- Test with and without configuration options
- Ensure your widget doesn't block the UI

### Documentation

- Update README.md if you add new features
- Add comments to your code
- Update config.yaml.example if you add new configuration options

## Pull Request Process

1. Ensure your code builds without errors
2. Update documentation as needed
3. Write a clear PR description explaining your changes
4. Link any relevant issues
5. Be responsive to feedback

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers
- Focus on constructive feedback
- Help others learn and grow

## Questions?

If you have questions, feel free to:
- Open an issue
- Start a discussion
- Reach out to the maintainers

Thank you for contributing to GoTUI!
