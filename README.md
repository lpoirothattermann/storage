# Manage your encrypted directories easily.
`storage` is a small application that has the goal of managing your encrypted directories with Age.

## How it works
`storage` manages a list of directories, which are called `states`. This list is defined under a configuration file located at `~/.config/storage/storage.toml`.

### Configuration file example
```toml
[states]

[states.example_state]
private_key_path = "~/key.txt"
encrypted_path = "~/.storage/"
decrypted_path = "~/Documents/storage"
```
In this example, `example_state` is the name of the state. `encrypted_path` is the path where your archive will be stored. Archives contain only one root directory which has the name of the state. Archives are tarballs that are compressed and obviously encrypted. In this example, the path of the archive will be `~/.storage/example_state.tar.gz.age`.
The next and last field is `decrypted_path`, it represents the path where your archive will be decrypted, uncompressed, and unarchived. The final path will be `~/Documents/storage/example_state/` in this example.

### Init state
The first thing you want to do, even before modifying your configuration file, is to init a new archive that will be used in your state.
```bash
$ storage init ~/key.txt ~/.storage/ example_state
```
This command will create your encrypted archive `~/.storage/example_state.tar.gz.age` with the Age key located at `~/key.txt`. You can now create the new `state` with the name `example_state` in your configuration file and start using it.

### Open state
When you want to access your files/directories that are stored in a `state`, you need to `open` the `state`.
```bash
$ storage open example_state
```
This command will decrypt, uncompress, and unarchive your state to where you have indicated in the configuration file, `~/Documents/storage/example_state/` in this example. You can now read and write in this directory.

### Close state
When you want to convert your directory back to an encrypted archive, you need to `close` the `state`.
```bash
$ storage close example_state
```
This command will do the opposite of the `open` command. It will archive, compress, and encrypt your directory. Now your archive contains all the modifications that you made.Also `~/Documents/storage/example_state` is automatically deleted.

#### Information notice
A state needs to be `close` to be `open`, and vice versa.
