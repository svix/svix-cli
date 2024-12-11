use crate::cli_types::PostOptions;
use crate::get_client_options;
use clap::{Args, Subcommand};
use colored_json::ColorMode;

#[derive(Args)]
#[command(args_conflicts_with_subcommands = true)]
#[command(flatten_help = true)]
pub struct AuthenticationArgs {
    #[command(subcommand)]
    pub command: AuthenticationCommands,
}

#[derive(Subcommand)]
pub enum AuthenticationCommands {
    /// Get a dashboard URL for the given app.
    /// Deprecated: use `app-portal` instead.
    #[clap(alias = "dashboard")]
    DashboardAccess {
        app_id: String,
        #[clap(flatten)]
        post_options: Option<PostOptions>,
    },
    /// Invalidates the given dashboard key
    Logout {
        dashboard_auth_token: String,
        #[clap(flatten)]
        post_options: Option<PostOptions>,
    },
}

impl AuthenticationCommands {
    pub async fn exec(self, client: &svix::api::Svix, color_mode: ColorMode) -> anyhow::Result<()> {
        match self {
            AuthenticationCommands::DashboardAccess {
                app_id,
                post_options,
            } => {
                let resp = client
                    .authentication()
                    .dashboard_access(app_id, post_options.map(Into::into))
                    .await?;

                crate::json::print_json_output(&resp, color_mode)?;
            }
            AuthenticationCommands::Logout {
                dashboard_auth_token,
                post_options,
            } => {
                // We're not using the client received by `exec()` here since the token is an
                // arg, not whatever is configured for the cli otherwise.
                let client =
                    svix::api::Svix::new(dashboard_auth_token, Some(get_client_options()?));

                client
                    .authentication()
                    .logout(post_options.map(Into::into))
                    .await?;
            }
        }
        Ok(())
    }
}
