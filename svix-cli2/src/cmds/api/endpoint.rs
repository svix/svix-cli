use crate::cli_types::endpoint::EndpointListOptions;
use crate::cli_types::PostOptions;
use crate::json::JsonOf;
use clap::{Args, Subcommand};
use colored_json::ColorMode;
use svix::api::{EndpointHeadersIn, EndpointHeadersPatchIn, EndpointIn, EndpointUpdate};

#[derive(Args)]
#[command(args_conflicts_with_subcommands = true)]
#[command(flatten_help = true)]
pub struct EndpointArgs {
    #[command(subcommand)]
    pub command: EndpointCommands,
}

#[derive(Subcommand)]
pub enum EndpointCommands {
    /// Create a new endpoint
    Create {
        app_id: String,
        endpoint_in: JsonOf<EndpointIn>,
        #[clap(flatten)]
        post_options: Option<PostOptions>,
    },
    /// Delete an endpoint by id
    Delete { app_id: String, id: String },
    /// Get an endpoint by id
    Get { app_id: String, id: String },
    /// Get custom headers for endpoint by id
    GetHeaders { app_id: String, id: String },
    /// List current endpoints
    List {
        app_id: String,
        #[clap(flatten)]
        options: EndpointListOptions,
    },
    /// Patch custom headers for endpoint by id
    PatchHeaders {
        app_id: String,
        id: String,
        endpoint_headers_patch_in: JsonOf<EndpointHeadersPatchIn>,
    },
    /// Get an endpoint's secret by id
    Secret { app_id: String, id: String },
    /// Update an endpoint by id
    Update {
        app_id: String,
        id: String,
        endpoint_update: JsonOf<EndpointUpdate>,
        #[clap(flatten)]
        post_options: Option<PostOptions>,
    },
    /// Update custom headers for endpoint by id
    UpdateHeaders {
        app_id: String,
        id: String,
        endpoint_headers_in: JsonOf<EndpointHeadersIn>,
    },
}

impl EndpointCommands {
    pub async fn exec(self, client: &svix::api::Svix, color_mode: ColorMode) -> anyhow::Result<()> {
        match self {
            EndpointCommands::Create {
                app_id,
                endpoint_in,
                post_options,
            } => {
                let resp = client
                    .endpoint()
                    .create(
                        app_id,
                        endpoint_in.into_inner(),
                        post_options.map(Into::into),
                    )
                    .await?;
                crate::json::print_json_output(&resp, color_mode)?;
            }
            EndpointCommands::Delete { app_id, id } => {
                client.endpoint().delete(app_id, id).await?;
            }
            EndpointCommands::Get { app_id, id } => {
                let resp = client.endpoint().get(app_id, id).await?;
                crate::json::print_json_output(&resp, color_mode)?;
            }
            EndpointCommands::GetHeaders { app_id, id } => {
                let resp = client.endpoint().get_headers(app_id, id).await?;
                crate::json::print_json_output(&resp, color_mode)?;
            }
            EndpointCommands::List { app_id, options } => {
                let resp = client.endpoint().list(app_id, Some(options.into())).await?;
                crate::json::print_json_output(&resp, color_mode)?;
            }
            EndpointCommands::PatchHeaders {
                app_id,
                id,
                endpoint_headers_patch_in,
            } => {
                client
                    .endpoint()
                    .patch_headers(app_id, id, endpoint_headers_patch_in.into_inner())
                    .await?;
            }
            EndpointCommands::Secret { app_id, id } => {
                let resp = client.endpoint().get_secret(app_id, id).await?;
                crate::json::print_json_output(&resp, color_mode)?;
            }
            EndpointCommands::Update {
                app_id,
                id,
                endpoint_update,
                post_options,
            } => {
                let resp = client
                    .endpoint()
                    .update(
                        app_id,
                        id,
                        endpoint_update.into_inner(),
                        post_options.map(Into::into),
                    )
                    .await?;
                crate::json::print_json_output(&resp, color_mode)?;
            }
            EndpointCommands::UpdateHeaders {
                app_id,
                id,
                endpoint_headers_in,
            } => {
                client
                    .endpoint()
                    .update_headers(app_id, id, endpoint_headers_in.into_inner())
                    .await?;
            }
        }
        Ok(())
    }
}
