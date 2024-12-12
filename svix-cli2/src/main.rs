use crate::cmds::api::authentication::AuthenticationArgs;
use crate::cmds::api::endpoint::EndpointArgs;
use crate::signature::SignatureArgs;
use anyhow::Result;
use clap::{Parser, Subcommand};
use clap_complete::Shell;
use cmds::api::application::ApplicationArgs;
use colored_json::{ColorMode, Output};
use concolor_clap::{Color, ColorChoice};

const VERSION: &str = env!("CARGO_PKG_VERSION");

mod cli_types;
mod cmds;
mod completion;
mod json;
mod signature;

#[derive(Parser)]
#[clap(color = concolor_clap::color_choice(), bin_name="svix-cli")]
struct Cli {
    #[command(flatten)]
    color: Color,
    #[command(subcommand)]
    command: RootCommands,
}

impl Cli {
    /// Converts the selected `ColorChoice` from the CLI to a `ColorMode` as used by the JSON printer.
    /// When the color choice is "auto", this considers whether stdout is a tty or not so that
    /// color codes are only produced when actually writing directly to a terminal.
    fn color_mode(&self) -> ColorMode {
        match self.color.color {
            ColorChoice::Auto => ColorMode::Auto(Output::StdOut),
            ColorChoice::Always => ColorMode::On,
            ColorChoice::Never => ColorMode::Off,
        }
    }
}

// N.b Ordering matters here for how clap presents the help.
#[derive(Subcommand)]
enum RootCommands {
    /// List, create & modify applications
    Application(ApplicationArgs),
    /// Manage authentication tasks such as getting dashboard URLs
    Authentication(AuthenticationArgs),
    /// List, create & modify endpoints
    Endpoint(EndpointArgs),
    /// List, create & modify event types
    EventType,
    /// Export data from your Svix Organization
    Export,
    /// Outputs shell completions for a variety of shells
    GenerateCompletions { shell: Shell },
    /// Import data to your Svix Organization
    Import,
    /// List integrations by app id
    Integration,
    /// Forward webhook requests to a local url
    Listen,
    /// Interactively configure your Svix API credentials
    Login,
    /// List & create messages
    Message,
    /// List, lookup & resend message attempts
    MessageAttempt,
    /// Quickly open Svix pages in your browser
    Open,
    /// Verifying and signing webhooks with the Svix signature scheme
    SvixSignature(SignatureArgs),
    /// Get the version of the Svix CLI
    Version,
}

#[tokio::main]
async fn main() -> Result<()> {
    let cli = Cli::parse();
    let color_mode = cli.color_mode();
    match cli.command {
        // Local-only things
        RootCommands::Version => println!("{VERSION}"),
        RootCommands::SvixSignature(args) => args.command.exec().await?,
        RootCommands::Open => todo!("Commands::Open"),
        // Remote API calls
        RootCommands::Application(args) => {
            let client = get_client()?;
            args.command.exec(&client, color_mode).await?;
        }
        RootCommands::Authentication(args) => {
            let client = get_client()?;
            args.command.exec(&client, color_mode).await?;
        }
        RootCommands::EventType => todo!("Commands::EventType"),
        RootCommands::Endpoint(args) => {
            let client = get_client()?;
            args.command.exec(&client, color_mode).await?;
        }
        RootCommands::Message => todo!("Commands::Message"),
        RootCommands::MessageAttempt => todo!("Commands::MessageAttempt"),
        RootCommands::Import => todo!("Commands::Import"),
        RootCommands::Export => todo!("Commands::Export"),
        RootCommands::Integration => todo!("Commands::Integration"),

        // FIXME: make login/listen play subcommands?
        RootCommands::Listen => todo!("Commands::Listen"),
        RootCommands::Login => todo!("Commands::Login"),
        RootCommands::GenerateCompletions { shell } => {
            completion::generate(&shell)?;
        }
    }

    Ok(())
}

fn get_client() -> Result<svix::api::Svix> {
    // XXX: Go client will exit if the token is not set. May need to rewrangle the flow.
    // FIXME: read from config

    // FIXME: don't hardcode ;)
    let token = env!("LOCAL_CLOUD_TOKEN").to_string();
    let opts = get_client_options()?;
    Ok(svix::api::Svix::new(token, Some(opts)))
}

fn get_client_options() -> Result<svix::api::SvixOptions> {
    // FIXME: read options from config file
    // FIXME: validate server url

    Ok(svix::api::SvixOptions {
        debug: false,
        // FIXME: don't hardcode ;)
        server_url: Some(env!("SVIX_ROOT").to_string()),
        timeout: None,
    })
}
