use std::{env, time::Duration};

use anyhow::Result;
use rand::RngExt;
use tokio::time::sleep;
use tracing_subscriber::EnvFilter;

use crate::{
    claimer::claim::{ClaimResult, claim},
    config::Config,
    constants::APP_VERSION,
};

mod claimer;
mod config;
mod constants;

#[tokio::main]
async fn main() -> Result<()> {
    if cfg!(debug_assertions) {
        tracing_subscriber::fmt()
            .with_max_level(tracing::Level::DEBUG)
            .init();
    } else {
        tracing_subscriber::fmt()
            .with_env_filter(
                EnvFilter::try_from_default_env().unwrap_or_else(|_| EnvFilter::new("info")),
            )
            .init();
    }

    tracing::info!("Pity Patrol v{} started!", APP_VERSION);

    assert!(
        !(env::var("GITHUB_ACTIONS").is_ok() || env::var("GITLAB_CI").is_ok()),
        "Unauthorized environment."
    );

    let config = Config::from_env()?;

    tracing::debug!("Config: {:?}", config);

    tracing::info!("{} account/s configured", config.accounts.len());

    let mut rng = rand::rng();

    for (index, account) in config.accounts.iter().enumerate() {
        let identifier = if let Some(identifier) = &account.identifier {
            format!(
                "Account #{} {} [{}]",
                index + 1,
                identifier,
                account.game_name()
            )
        } else {
            format!("Account #{} [{}]", index + 1, account.game_name())
        };

        tracing::info!("{} claiming...", identifier);

        match claim(&config, account).await {
            Ok(ClaimResult::Claimed) => {
                tracing::info!("{} claimed reward successfully!", identifier);
            }
            Ok(ClaimResult::AlreadyClaimed) => {
                tracing::info!("{} has already claimed reward", identifier);
            }
            Err(err) => {
                tracing::error!("{} could not claim rewards because: {:?}", identifier, err);
            }
        }

        // sleep a random time between querying accounts
        let delay = rng.random_range(500..2000);
        sleep(Duration::from_millis(delay)).await;
    }

    Ok(())
}
