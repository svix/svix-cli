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
    /// List current applications
    List(ApplicationListOptions),
    /// Creates a new application
    Create {
        application_in: JsonOf<ApplicationIn>,
        #[clap(flatten)]
        post_options: Option<PostOptions>,
    },
    /// Get an application by id
    Get { id: String },
    /// Update an application by id
    Update { id: String, body: String },
    /// Deletes an application by id
    Delete { id: String },
}

impl ApplicationCommands {
    // FIXME: codegen an exec() method that takes the args and a client and does the thing?
    //   Not sure if we need to pass in a printer or how the output should work if we can't
    //   have a typed return here.
    //   This might not make sense but let's roll with it for now.
    pub async fn exec(
        &self,
        client: &svix::api::Svix,
        color_mode: ColorMode,
    ) -> anyhow::Result<()> {
        match self {
            ApplicationCommands::List(options) => {
                let resp = client
                    .application()
                    .list(Some(options.clone().into()))
                    .await?;

                crate::json::print_json_output(&resp, color_mode)?;
            }
            ApplicationCommands::Create {
                application_in,
                post_options,
            } => {
                let resp = client
                    .application()
                    .create(
                        application_in.clone().into_inner(),
                        post_options.clone().map(Into::into),
                    )
                    .await?;

                crate::json::print_json_output(&resp, color_mode)?;
            }
            ApplicationCommands::Get { id } => todo!("application get"),
            ApplicationCommands::Update { id, body } => todo!("application update"),
            ApplicationCommands::Delete { id } => todo!("application delete"),
        }
        Ok(())
    }
}