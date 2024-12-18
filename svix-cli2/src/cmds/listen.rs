use crate::config::{get_config_file_path, Config};
use anyhow::{Context, Result};
use clap::Args;

#[derive(Args)]
pub struct ListenArgs {
    /// The local URL to forward webhooks to
    url: url::Url,
    /// Disables History Logging
    #[clap(long)]
    no_logging: bool,
}

impl ListenArgs {
    pub async fn exec(self, cfg: &Config) -> Result<()> {
        let token = match cfg.relay_token.as_ref() {
            None => {
                let token = crate::relay::token::generate_token()?;
                let mut updated_cfg = cfg.clone();
                updated_cfg.auth_token = Some(token.clone());

                let cfg_path = get_config_file_path()?;
                if let Err(e) = updated_cfg.save_to_disk(&cfg_path).context(format!(
                    "failed to save relay token to config file at `{}`",
                    cfg_path.as_os_str().to_str().unwrap_or_default()
                )) {
                    eprintln!("{e}");
                }
                token
            }
            Some(token) => token.clone(),
        };
        crate::relay::listen(
            self.url,
            token,
            !self.no_logging,
            cfg.relay_debug_url.as_deref(),
            cfg.relay_disable_security.unwrap_or_default(),
        )
        .await?;
        Ok(())
    }
}
