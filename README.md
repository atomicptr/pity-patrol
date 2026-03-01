# Pity Patrol

A tool to help you claim web based daily login rewards in your favorite gacha games

**NOTE**: This isn't intended for use inside GitHub Actions or GitLab CI. Using it there violates the TOS, so don't.

## Supported Games

- [Arknights: Endfield](https://endfield.gryphline.com)

## Install

Create a config file (see below) and then run:

```bash
$ docker run --rm -v /path/to/config/dir:/app/config quay.io/atomicptr/pity-patrol:latest
```

## Configuration

Pity Patrol reads a TOML config file located in

- Linux: ``$XDG_CONFIG_HOME/.config/pity-patrol/config.toml``
- MacOS: ``$HOME/Library/Application Support/pity-patrol/config.toml``
- Windows: ``%APPDATA%\\pity-patrol\\config.toml``

Or a path defined by the env var ``PITY_PATROL_CONFIG``

Here is a list of all configuration options

```toml
user-agent = "Pity Patrol" # Ability to set a custom user agent, keep empty for default (Chrome)

# and the most important thing you can add as many accounts as you like
[[accounts]]
# The game identifier, this is used to decide which game
game = "endfield"

# Account identifier, this will be used in logs and with reporters so you can differentiate different accounts
identifier = "My Endfield Account"

# Read "Getting Credentias > Arknights: Endfield" to see how to get these
credentials = "xxxxx"
sk-game-role = "xxxxx"
```

## Getting Credentials

### Arknights: Endfield

1. Open [https://game.skport.com/endfield/sign-in](https://game.skport.com/endfield/sign-in) (Make sure you are logged in)
2. Open the Browser Console (F12)
3. Go to the Network tab (make sure its recording, there should be a red dot)
4. Search for "attendance"
5. On the right click to "Headers"
6. Look for "Sk-Game-Role" (which is 'sk-game-role') and "Cred" (which is 'credentials')
7. Add them to the config

![Endfield Image Tutorial](./.github/endfield_tutorial.png)

## License

AGPLv3
