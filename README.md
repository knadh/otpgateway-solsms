## Solutions Infini SMS provider for [otpgateway](https://github.com/knadh/otpgateway).
This is a Provider plugin for [otpgateway](https://github.com/knadh/otpgateway) that sends SMSes using [SolutionsInfini](https://www.solutionsinfini.com), an Indian SMS gateway.

## Build
- Run `make build` to produce `solsms.prov`

## Usage
- Add this configuration to otpgateway's config.toml
```toml
[provider.solsms]
subject = "Verification"
template = "static/sms.txt"
template_type = "text"
config = '{"APIKey": "YourSolutionsInfiniKey", "Sender": "YourID"}'
```

- Run `./otpgateway --prov solsms.prov`