# enshamir

`enshamir` is a command line tool for splitting and combining secrets using Shamir's Secret Sharing Scheme. The secret
is encrypted by AES-256-GCM before splitting.

## Split

```shell
go run ./cmd/enshamir split \
  --secret-file <secret-file> \
  --parts <number-of-generated-shares> \
  --threshold <minimum-number-of-shares-to-reconstruct> \
  --output-dir <output-directory>
```

Let's say we want to split a secret into 5 shares, and 3 of them are required to recover the secret. `enshamir` encrypts
the secret before splitting it, so we need to provide a password to encrypt the secret.

There are 3 steps in the splitting process:

1. The password is hashed to a 32 bytes key by `argon2id`. A random 16 bytes salt is also generated.
2. The secret is encrypted by `AES-256-GCM` with the 32 bytes key.
3. The encrypted secret is split into 5 shares using Shamir's Secret Sharing Scheme.

You must save the random 16 bytes salt file and 3 shares to recover the secret.

The first step is preparing a secret file. We create a test file `/tmp/test-secret` by `echo`.

```shell
echo "my secret" > /tmp/test-secret
```

Then we split the secret into 5 shares which will be saved in `/tmp/splitted-secrets`. The `--threshold 3` means we need
at least 3 shares to reconstruct the secret.

`enshamir` will ask you to enter a password to encrypt the secret. The password is hashed to a 32 bytes key by `argon2id`.

```shell
go run ./cmd/enshamir split \
  --secret-file /tmp/test-secret \
  --parts 5 \
  --threshold 3 \
  --output-dir /tmp/splitted-secrets
```

There will be a `MUST-BACK-UP-SALT` file and a `shares` directory under `/tmp/splitted-secrets`. The generated files are:
```
# random 16 bytes salt which is necessary to hash the password. You must back up this file.
/tmp/splitted-secrets/MUST-BACK-UP-SALT

# 5 shares are generated. You must backup at least 3 shares.
/tmp/splitted-secrets/SPLITTED-SECRET-1
/tmp/splitted-secrets/SPLITTED-SECRET-2
/tmp/splitted-secrets/SPLITTED-SECRET-3
/tmp/splitted-secrets/SPLITTED-SECRET-4
/tmp/splitted-secrets/SPLITTED-SECRET-5
```

You **MUST** back up the `MUST-BACK-UP-SALT` file and at least 3 shares.

I strongly recommend you to verify the generated shares by combining them immediately.

## Combine

```shell
go run ./cmd/enshamir combine \
  --salt-file <salt-file-path> \
  --shares-dir <shares-directory> \
  --secret-file <output-secret-file-path>
```

After splitting the secret into shares, I strongly recommend you to combine parts of shares to the secret immediately.

Assume we have 1 salt file and 3 shares.

```
/tmp/MUST-BACK-UP-SALT

/tmp/parts-of-secrets/SPLITTED-SECRET-1
/tmp/parts-of-secrets/SPLITTED-SECRET-3
/tmp/parts-of-secrets/SPLITTED-SECRET-5
```

We have to pass the salt file path and shares directory to `enshamir` to combine the shares. `--secret-file` is the path
of combined secret file.

```shell
go run ./cmd/enshamir combine \
  --salt-file /tmp/MUST-BACK-UP-SALT \
  --shares-dir /tmp/parts-of-secrets \
  --secret-file /tmp/combined-secret
```

The combined secret is written to `/tmp/combined-secret`.

```shell
cat /tmp/combined-secret
```

The content should be `my secret`.