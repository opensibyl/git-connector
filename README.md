# git-connector

Receive webhook, clone, analyze and upload.

## What is it

- receive webhooks from gitlab/github/...
- check auth
- clone in-memory
- analyze
- upload

## Usage

You can download our prebuild executions in [release page](https://github.com/opensibyl/git-connector/releases).

And deploy it with:

```bash
./git-connector --port 9448 --url http://127.0.0.1:9876 --gitlab_user YOUR_USERNAME --gitlab_pwd YOUR_PASSWORD_OR_TOKEN
```

## Support SCM

At the most time, this project was used for building some private workflows. So we support GitLab firstly.

| SCM Platform | Support? |
|--------------|----------|
| GitLab       | Yes      |
| GitHub       | Not yet  |

PR/issues are always welcome if you need some other support.

## Why

Offer a standard way to keep mirror between SCM and sibyl system.

## License

[Apache 2.0](LICENSE)
