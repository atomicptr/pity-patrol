use std::{env, fs, path::PathBuf};

use anyhow::{Context, Result};
use serde::Deserialize;

#[derive(Debug, Deserialize)]
pub struct Config {
    #[serde(rename = "user-agent")]
    pub user_agent: Option<String>,
    pub accounts: Vec<Account>,
}

#[derive(Debug, Deserialize)]
pub struct Account {
    pub identifier: Option<String>,

    #[serde(flatten)]
    pub game: Game,
}

impl Account {
    pub fn game_name(&self) -> &str {
        match self.game {
            Game::Endfield {
                credentials: _,
                sk_game_role: _,
            } => "Arknights: Endfield",
        }
    }
}

#[derive(Debug, Deserialize)]
#[serde(tag = "game", rename_all = "lowercase")]
pub enum Game {
    Endfield {
        credentials: String,

        #[serde(rename = "sk-game-role")]
        sk_game_role: String,
    },
}

impl Config {
    pub fn from_env() -> Result<Config> {
        if let Some(path) = env::var("PITY_PATROL_CONFIG")
            .ok()
            .map(PathBuf::from)
            .filter(|p| p.exists())
        {
            return Self::from_path(path);
        }

        Self::from_config_dir()
    }

    fn from_config_dir() -> Result<Config> {
        let config_file = dirs::config_dir()
            .context("could not find config dir")?
            .join("pity-patrol")
            .join("config.toml");

        if config_file.exists() {
            return Self::from_path(config_file);
        }

        anyhow::bail!("Could not find file: {}", config_file.display())
    }

    pub fn from_path(path: PathBuf) -> Result<Config> {
        if !path.exists() {
            anyhow::bail!("Could not find file: {}", path.display());
        }

        let data = fs::read_to_string(path)?;

        let config: Config = toml::from_str(data.as_str())?;

        Ok(config)
    }
}
