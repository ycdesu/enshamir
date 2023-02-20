# enshamir

`enshamir` is a powerful command line utility designed for splitting and combining secrets using Shamir's Secret Sharing Scheme, which is a cryptographic algorithm used to protect sensitive information. 

With `enshamir`, you can safely split a secret into multiple shares, each of which contains a piece of the original secret, and then combine those shares to reconstruct the original secret. To ensure maximum security, `enshamir` encrypts the secret using AES-256-GCM before splitting it, making it virtually impossible for unauthorized parties to access the secret.

## Installation

### Homebrew

```shell
brew install ycdesu/tap/enshamir
```

### Go
```shell
go install github.com/ycdesu/enshamir/cmd/enshamir@latest
```

## Splitting a Secret

To split a secret into shares using `enshamir`, follow these steps:

1. Prepare the secret file that you want to split. You can create a test file at `/tmp/test-secret` by running:

```shell
echo "my secret" > /tmp/test-secret
```

2. Split the secret into 5 shares using the following command. This will save the shares in `/tmp/splitted-secrets`:

```shell
go run ./cmd/enshamir split \
  --secret-file /tmp/test-secret \
  --parts 5 \
  --threshold 3 \
  --output-dir /tmp/splitted-secrets
```

In this example, we are splitting the secret into 5 shares, and we need at least 3 shares to reconstruct the secret.
Note that `enshamir` encrypts the secret before splitting it, so you will be asked to provide a password to encrypt the secret.
The password is hashed to a 32-byte key by `argon2id`, and a random 16-byte salt is also generated to add an extra layer of security.

3. When the splitting process is complete, `enshamir` will generate the following files under `/tmp/splitted-secrets`:

```shell
/tmp/splitted-secrets/MUST-BACK-UP-SALT     # a random 16-byte salt file used to hash the user password
/tmp/splitted-secrets/SPLITTED-SECRET-1     # share 1
/tmp/splitted-secrets/SPLITTED-SECRET-2     # share 2
/tmp/splitted-secrets/SPLITTED-SECRET-3     # share 3
/tmp/splitted-secrets/SPLITTED-SECRET-4     # share 4
/tmp/splitted-secrets/SPLITTED-SECRET-5     # share 5
```

It is **important** to back up the `MUST-BACK-UP-SALT` file and at least 3 shares to be able to recover the secret later.
The `MUST-BACK-UP-SALT` file is necessary to reconstruct the 32-byte key used to encrypt the secret, while at least 3 shares are needed to recover the encrypted secret.

4. To verify that the generated shares are valid, you can immediately combine them using the `enshamir` combine command.
If the secret is properly recovered, this will ensure that the shares were generated correctly.

## Combining Shares

Once you have split the secret into shares, it's recommended to combine the shares as soon as possible.

Assuming you have a salt file and three shares, and the secret content is "my secret". The shares are saved in the following files:

```shell
/tmp/MUST-BACK-UP-SALT
/tmp/parts-of-secrets/SPLITTED-SECRET-1
/tmp/parts-of-secrets/SPLITTED-SECRET-3
/tmp/parts-of-secrets/SPLITTED-SECRET-5
```

To combine the shares, use the `enshamir` combine command with the following arguments:

```shell
go run ./cmd/enshamir combine \
  --salt-file <path-to-salt-file> \
  --shares-dir <path-to-shares-directory> \
  --secret-file <path-to-combined-secret>
```

Here, you should replace `path-to-salt-file` with the actual path to your salt file, `path-to-shares-directory` with the 
actual path to your directory containing the split shares, and `path-to-combined-secret` with the actual path to the file 
where you want to store the combined secret.

So the combining command is:

```shell
go run ./cmd/enshamir combine \
  --salt-file /tmp/MUST-BACK-UP-SALT \
  --shares-dir /tmp/parts-of-secrets \
  --secret-file /tmp/combined-secret
```

After running the `enshamir` combine command, the combined secret will be saved to `path-to-combined-secret`. You can 
verify the content of this file by running:

```shell
cat /tmp/combined-secret
```

The content of the file should be `my secret`.
