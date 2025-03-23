# Known Issues and Limitations

## WhatsApp Library Limitations

The current implementation uses the `github.com/Rhymen/go-whatsapp` library, which has several limitations:

1. **Limited History Access**: The library doesn't provide direct methods to access historical messages. The `GetMessages` and `GetChats` functions in our code are placeholders that will need custom implementation based on your specific requirements.

2. **WhatsApp API Changes**: WhatsApp frequently updates its protocols and security measures, which may break third-party libraries. Make sure to monitor for updates and be prepared to adapt.

## Potential Alternatives

If you need more robust WhatsApp data access, consider:

1. **Using WhatsApp Export Files**: WhatsApp allows users to export chat history. You could modify this tool to parse these export files instead of directly connecting to WhatsApp.

2. **Browser Automation**: Use a Selenium or Puppeteer-based approach to automate WhatsApp Web for data extraction.

3. **Official WhatsApp Business API**: For business use cases, consider using the official WhatsApp Business API.

## Dependency Management Issues

### Missing go.sum Entries

When you see errors like:

```
missing go.sum entry for module providing package github.com/spf13/viper
```

These happen when Go can't find checksums for dependencies in your go.sum file. Solutions:

1. **Initialize dependencies with go mod tidy**:

   ```bash
   go mod tidy
   ```

2. **Explicitly download each dependency**:

   ```bash
   go get github.com/Rhymen/go-whatsapp@v0.1.1
   go get github.com/go-telegram-bot-api/telegram-bot-api/v5@v5.5.1
   go get github.com/spf13/cobra@v1.7.0
   go get github.com/spf13/viper@v1.16.0
   ```

3. **Clean module cache** (if persistent issues):
   ```bash
   go clean -modcache
   go mod tidy
   ```

### Checksum Mismatch

If you encounter checksum mismatches when building the project:

```
verifying github.com/spf13/viper@v1.16.0: checksum mismatch
        downloaded: h1:rGGH0XDZhdUOryiDWjmIvUSWpbNqisK8Wk0Vyefw8hc=
        go.sum:     h1:vawHUee0VqjUf7VdnU56IQBpO2qr4qk7s8S3tzJvTW8=
```

Fix it by removing the go.sum file and letting Go regenerate it:

```bash
rm go.sum
go mod download
```

## Telegram Rate Limits

Be aware that Telegram imposes rate limits on bot API calls. The current implementation includes basic throttling (100ms delay between messages), but you might need to adjust this based on your usage patterns.

## Future Improvements

Consider implementing:

1. A proper message sorting mechanism in the `ExportMessages` function
2. Better error handling and retries for network failures
3. Support for more message types and formats
4. A web interface for easier configuration
