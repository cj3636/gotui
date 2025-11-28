# Security Policy

## Supported Versions

Currently supporting the latest version only.

| Version | Supported          |
| ------- | ------------------ |
| latest  | :white_check_mark: |

## Security Considerations

### API Tokens

- Store GitHub and GitLab API tokens securely in `config.yaml`
- Never commit API tokens to version control
- Use minimal required permissions for API tokens
  - GitHub: `repo` scope for private repositories
  - GitLab: `read_api`, `read_repository` for private projects
- API tokens are never logged or displayed in the UI

### Network Requests

The application makes HTTP requests to the following services:
- **wttr.in** - Weather data (no authentication required)
- **ipinfo.io** - IP geolocation data (no authentication required)
- **api.github.com** - GitHub API (optional authentication)
- **gitlab.com** - GitLab API (optional authentication)

All requests are made over HTTPS.

### File Access

The application may read files specified in the configuration:
- Text files via `text_file` configuration
- Markdown files via `markdown_file` configuration
- Configuration file (`config.yaml`)

File access is limited to:
- Reading only (no write operations)
- Files explicitly specified in configuration
- Local filesystem only

### System Information

The application accesses system information for monitoring:
- CPU usage (read-only via gopsutil)
- Memory usage (read-only via gopsutil)
- Disk usage (read-only via standard tools)
- SMART data (read-only, requires smartmontools)

No system modifications are made.

## Reporting a Vulnerability

If you discover a security vulnerability in GoTUI, please:

1. **Do NOT** open a public issue
2. Email the maintainers directly (see repository contacts)
3. Provide detailed information about the vulnerability
4. Allow time for a fix to be developed and deployed

We will:
- Acknowledge receipt within 48 hours
- Provide an estimated timeline for a fix
- Credit you in the security advisory (unless you prefer anonymity)

## Security Best Practices

When using GoTUI:

1. **Protect your configuration file**
   ```bash
   chmod 600 ~/.config/gotui/config.yaml
   ```

2. **Use environment variables for sensitive data** (alternative to config file)
   ```bash
   export GITHUB_TOKEN="your_token"
   export GITLAB_TOKEN="your_token"
   ```

3. **Regularly rotate API tokens**

4. **Review file paths in configuration** before running

5. **Keep dependencies updated**
   ```bash
   go get -u ./...
   go mod tidy
   ```

## Audit Information

- Last security audit: 2025-11-28
- CodeQL analysis: No vulnerabilities found
- Dependencies: Regularly updated from Charmbracelet ecosystem

## Contact

For security concerns, please contact the repository maintainers through GitHub.
