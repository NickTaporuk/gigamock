# gigamock
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FNickTaporuk%2Fgigamock.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FNickTaporuk%2Fgigamock?ref=badge_shield)

Gigamock is a generic utility to be used as mocked server to provide a simplified and consistent API over various network data sources such as http REST API or grpc or graphql services or as a mocking messaging systems like a kafka via mocking response or send message to some message system.
# Conception
# Download

## Precompiled Binaries

You can download the precompiled release binary from [releases](https://github.com/NickTaporuk/gigamock/releases/) via web
or via

```bash
wget https://github.com/NickTaporuk/gigamock/releases/<version>/gigamock_<version>_<os>_<arch>
```

#### Go get

You can also use Go 1.12 or later to build the latest stable version from source:

```bash
go get github.com/NickTaporuk/gigamock
```

#### Homebrew Tap

```bash
brew install nicktaporuk/tap/gigamock
# After initial install you can upgrade the version via:
brew upgrade gigamock
```

### Scenarios
### Feature
    grpc api mock
    graphql api mock
    parse swagger api to mock scenarios
    benchmarks as performance tests(REST API, kafka topics, graphql APIs, NATS and so one)



## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FNickTaporuk%2Fgigamock.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FNickTaporuk%2Fgigamock?ref=badge_large)