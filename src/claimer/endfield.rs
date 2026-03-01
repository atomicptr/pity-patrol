use std::time::{SystemTime, UNIX_EPOCH};

use crate::{
    claimer::claim::ClaimResult,
    config::{Account, Config, Game},
};
use anyhow::Result;
use hmac::{Hmac, Mac};
use md5::{Digest, Md5};
use reqwest::header::{HeaderMap, HeaderValue};
use sha2::Sha256;

const BASE_URL: &str = "https://zonai.skport.com";
const CLAIM_URL: &str = "/web/v1/game/endfield/attendance";
const REFRESH_URL: &str = "/web/v1/auth/refresh";
const PLATFORM: &str = "3";
const V_NAME: &str = "1.0.0";

type HmacSha256 = Hmac<Sha256>;

pub async fn claim(config: &Config, account: &Account) -> Result<ClaimResult> {
    let client = reqwest::Client::builder()
        .user_agent(
            config
                .user_agent
                .clone()
                .unwrap_or(crate::constants::USER_AGENT.to_string()),
        )
        .build()?;

    let Game::Endfield {
        credentials,
        sk_game_role,
    } = &account.game;

    let token = refresh_token(&client, credentials).await?;

    let timestamp = SystemTime::now()
        .duration_since(UNIX_EPOCH)?
        .as_secs()
        .to_string();

    let sign = generate_sign(CLAIM_URL, "", &timestamp, &token);

    let mut headers = HeaderMap::new();
    headers.insert("cred", HeaderValue::from_str(credentials)?);
    headers.insert("sk-game-role", HeaderValue::from_str(sk_game_role)?);
    headers.insert("platform", HeaderValue::from_static(PLATFORM));
    headers.insert("vName", HeaderValue::from_static(V_NAME));
    headers.insert("timestamp", HeaderValue::from_str(&timestamp)?);
    headers.insert("sign", HeaderValue::from_str(&sign)?);
    headers.insert("Content-Type", HeaderValue::from_static("application/json"));
    headers.insert("sk-language", HeaderValue::from_static("en"));

    let res = client
        .post(format!("{BASE_URL}{CLAIM_URL}"))
        .headers(headers)
        .send()
        .await?;

    let json: serde_json::Value = res.json().await?;

    tracing::debug!("Response: {json:?}");

    let code = json["code"].as_i64().unwrap_or(-1);

    match code {
        0 => Ok(ClaimResult::Claimed),
        10001 => Ok(ClaimResult::AlreadyClaimed),
        _ => {
            anyhow::bail!(json["message"].to_string())
        }
    }
}

async fn refresh_token(client: &reqwest::Client, credentials: &str) -> Result<String> {
    let resp = client
        .get(format!("{BASE_URL}{REFRESH_URL}"))
        .header("cred", credentials)
        .header("platform", PLATFORM)
        .header("vName", V_NAME)
        .send()
        .await?;

    let json: serde_json::Value = resp.json().await?;

    if json["code"] != 0 {
        anyhow::bail!("Token refresh failed: {}", json["message"]);
    }

    Ok(json["data"]["token"].as_str().unwrap_or("").to_string())
}

fn generate_sign(path: &str, body: &str, timestamp: &str, token: &str) -> String {
    let header_json = format!(
        r#"{{"platform":"{PLATFORM}","timestamp":"{timestamp}","dId":"","vName":"{V_NAME}"}}"#,
    );

    let data_to_sign = format!("{path}{body}{timestamp}{header_json}");

    let mut mac =
        HmacSha256::new_from_slice(token.as_bytes()).expect("HMAC can take key of any size");
    mac.update(data_to_sign.as_bytes());
    let hmac_result = hex::encode(mac.finalize().into_bytes());
    let mut hasher = Md5::new();
    hasher.update(hmac_result.as_bytes());
    hex::encode(hasher.finalize())
}
