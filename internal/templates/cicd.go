package templates

// GithubActionCiTemplate contains the GitHub Actions CI workflow
const GithubActionCiTemplate = `name: CI

on:
  push:
    branches: [ main, master ]
  pull_request:
    branches: [ main, master ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Install dependencies
      run: go mod tidy

    - name: Run Build
      run: go build -v ./...

    - name: Run Tests
      run: go test -v ./...
`
