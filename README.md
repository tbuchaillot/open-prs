# Open PRs Tracker

## Overview
This project is a Go application designed to fetch and display the number of open pull requests for each contributor in a specified GitHub repository. It supports outputting the data to the console or to a CSV file, including the URLs of the pull requests.

## Features
- Fetch open pull requests from a specified GitHub repository.
- Count open pull requests per contributor.
- Output the data to the console or a CSV file.
- Sort contributors by the number of open pull requests in descending order.

## Requirements
- Go 1.19 or higher (as specified in `go.mod`).

## Installation
```console
go install github.com/tbuchaillot/open-prs
```

## Usage
```console
open-prs -org <org_name> -repository <repository_name> -output <output_type> -token <github_token>
```
Where:
```
  -org string
        Github organization name (required)
  -repository string
        Github repository name (required)
  -output string
        ouput type (stdout,csv)  (default "stdout")
  -token string
        Github Personal Access Token (PAT) (required for private repos)
```

### G Personal Access Token
You can generate a GitHub Personal Access Token (PAT) by following the instructions [here](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token).

## Contributing
Contributions are welcome! Please feel free to submit a pull request.
