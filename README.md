# GitPullExecutor

GitPullExecutor is a lightweight tool that continuously pulls an already cloned git repository and executes a given command whenever a change is pulled.


## Building

```bash
go build
```

## Usage Prerequisites

The repository must be cloned using an RSA key with SHA-2, as SHA-1 is no longer allowed. Using a SHA-1 key will result in the following error message:

```bash
Unable to pull repository: unknown error: ERROR: You\'re using an RSA key with SHA-1, which is no longer allowed. Please use a newer client or a different key type.. Retrying...
```

To create a SHA-2 key:

```bash
# Create ecdsa-sha2-nistp521 key
ssh-keygen -t ecdsa -b 521 -C "your_email@example.com"
# Then add the key as default to your GitHub account 
```

## Usage
See `git_pull_executor --help` for detailed usage.