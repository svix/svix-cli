use crate::cli_types::application::ApplicationListOptions;
use crate::cli_types::PostOptions;
use crate::json::JsonOf;
use clap::{Args, Subcommand};
use colored_json::ColorMode;
use svix::api::ApplicationIn;

#[derive(Args)]
#[command(args_conflicts_with_subcommands = true)]
#[command(flatten_help = true)]
pub struct ApplicationArgs {
    #[command(subcommand)]
    pub command: ApplicationCommands,
}

// FIXME: build these via codegen from the spec, along with the rust lib.
#[derive(Subcommand)]
pub enum ApplicationCommands {
    /// Creates a new application
    Create {
        application_in: JsonOf<ApplicationIn>,
        #[clap(flatten)]
        post_options: Option<PostOptions>,
    },
    /// Deletes an application by id
    Delete { id: String },
    /// Get an application by id
    Get { id: String },
    /// List current applications
    List(ApplicationListOptions),
    /// Update an application by id
    Update {
        id: String,
        application_in: JsonOf<ApplicationIn>,
        #[clap(flatten)]
        post_options: Option<PostOptions>,
    },
}

impl ApplicationCommands {
    // FIXME: codegen an exec() method that takes the args and a client and does the thing?
    //   Not sure if we need to pass in a printer or how the output should work if we can't
    //   have a typed return here.
    //   This might not make sense but let's roll with it for now.
    pub async fn exec(self, client: &svix::api::Svix, color_mode: ColorMode) -> anyhow::Result<()> {
        match self {
            ApplicationCommands::List(options) => {
                let resp = client.application().list(Some(options.into())).await?;

                crate::json::print_json_output(&resp, color_mode)?;
            }
            ApplicationCommands::Create {
                application_in,
                post_options,
            } => {
                let resp = client
                    .application()
                    .create(application_in.into_inner(), post_options.map(Into::into))
                    .await?;

                crate::json::print_json_output(&resp, color_mode)?;
            }
            ApplicationCommands::Get { id } => {
                let resp = client.application().get(id).await?;
                crate::json::print_json_output(&resp, color_mode)?;
            }
            ApplicationCommands::Update {
                id,
                application_in,
                post_options,
            } => {
                let resp = client
                    .application()
                    .update(
                        id,
                        application_in.into_inner(),
                        post_options.map(Into::into),
                    )
                    .await?;

                crate::json::print_json_output(&resp, color_mode)?;
            }
            ApplicationCommands::Delete { id } => {
                client.application().delete(id).await?;
            }
        }
        Ok(())
    }
}
