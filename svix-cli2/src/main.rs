use crate::cmds::api::authentication::AuthenticationArgs;
use crate::cmds::api::endpoint::EndpointArgs;
use anyhow::Result;
use clap::{Parser, Subcommand};
use cmds::api::application::ApplicationArgs;
use colored_json::{ColorMode, Output};
use concolor_clap::{Color, ColorChoice};

const VERSION: &str = env!("CARGO_PKG_VERSION");

mod cli_types;
mod cmds;
mod json;

#[derive(Parser)]
#[clap(color = concolor_clap::color_choice())]
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
// FIXME: double-check Go cli. Seems like cobra may sort the items in the help lexigraphically
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
    /// Verify the signature of a webhook message
    Verify,
    /// Get the version of the Svix CLI
    Version,
}

#[tokio::main]
async fn main() -> Result<()> {
    let cli = Cli::parse();

    match &cli.command {
        // Local-only things
        RootCommands::Version => println!("{VERSION}"),
        RootCommands::Verify => todo!("Commands::Verify"),
        RootCommands::Open => todo!("Commands::Open"),
        // Remote API calls
        RootCommands::Application(args) => {
            let client = get_client()?;
            args.command.exec(&client, cli.color_mode()).await?;
        }
        RootCommands::Authentication(args) => {
            let client = get_client()?;
            args.command.exec(&client, cli.color_mode()).await?;
        }
        RootCommands::EventType => todo!("Commands::EventType"),
        RootCommands::Endpoint(args) => {
            let client = get_client()?;
            args.command.exec(&client, cli.color_mode()).await?;
        }
        RootCommands::Message => todo!("Commands::Message"),
        RootCommands::MessageAttempt => todo!("Commands::MessageAttempt"),
        RootCommands::Import => todo!("Commands::Import"),
        RootCommands::Export => todo!("Commands::Export"),
        RootCommands::Integration => todo!("Commands::Integration"),

        // FIXME: make login/listen play subcommands?
        RootCommands::Listen => todo!("Commands::Listen"),
        RootCommands::Login => todo!("Commands::Login"),
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
