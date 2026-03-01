use anyhow::Result;

use crate::{
    claimer::endfield,
    config::{Account, Config, Game},
};

#[derive(Debug)]
pub enum ClaimResult {
    Claimed,
    AlreadyClaimed,
}

pub async fn claim(config: &Config, account: &Account) -> Result<ClaimResult> {
    match &account.game {
        Game::Endfield {
            credentials: _,
            sk_game_role: _,
        } => endfield::claim(config, account),
    }
    .await
}
