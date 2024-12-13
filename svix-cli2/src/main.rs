use crate::cmds::api::authentication::AuthenticationArgs;
use crate::cmds::api::endpoint::EndpointArgs;
use crate::cmds::api::event_type::EventTypeArgs;
use crate::cmds::api::integration::IntegrationArgs;
use crate::cmds::api::message::MessageArgs;
use crate::config::Config;
use crate::signature::SignatureArgs;
use anyhow::Result;
use clap::{Parser, Subcommand};
use clap_complete::Shell;
use cmds::api::application::ApplicationArgs;
use colored_json::{ColorMode, Output};
use concolor_clap::{Color, ColorChoice};

const VERSION: &str = env!("CARGO_PKG_VERSION");
const DEFAULT_SERVER_URL: &str = "https://api.svix.com";

mod cli_types;
mod cmds;
mod completion;
mod config;
mod json;
mod login;
mod signature;

#[derive(Parser)]
#[command(version, about, long_about = None)]
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
#[derive(Subcommand)]
enum RootCommands {
    /// List, create & modify applications
    Application(ApplicationArgs),
    /// Manage authentication tasks such as getting dashboard URLs
    Authentication(AuthenticationArgs),
    /// List, create & modify endpoints
    Endpoint(EndpointArgs),
    /// List, create & modify event types
    EventType(EventTypeArgs),
    /// Export data from your Svix Organization
    Export,
    /// Outputs shell completions for a variety of shells
    GenerateCompletions { shell: Shell },
    /// Import data to your Svix Organization
    Import,
    /// List integrations by app id
    Integration(IntegrationArgs),
    /// Forward webhook requests to a local url
    Listen,
    /// Interactively configure your Svix API credentials
    Login,
    /// List & create messages
    Message(MessageArgs),
    /// List, lookup & resend message attempts
    // FIXME: need codegen for this
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
    // XXX: cfg can give an Err in certain situations.
    // Assigning the variable here since several match arms need a `&Config` but the rest of them
    // won't care/are still usable if the config doesn't exist.
    // To this, the `?` is deferred until the point inside a given match arm needs the config value.
    let cfg = Config::load();
    match cli.command {
        // Local-only things
        RootCommands::Version => println!("{VERSION}"),
        RootCommands::SvixSignature(args) => args.command.exec().await?,
        RootCommands::Open => todo!("Commands::Open"),
        // Remote API calls
        RootCommands::Application(args) => {
            let client = get_client(&cfg?)?;
            args.command.exec(&client, color_mode).await?;
        }
        RootCommands::Authentication(args) => {
            let cfg = cfg?;
            let client = get_client(&cfg)?;
            args.command.exec(&client, color_mode, &cfg).await?;
        }
        RootCommands::EventType(args) => {
            let client = get_client(&cfg?)?;
            args.command.exec(&client, color_mode).await?;
        }
        RootCommands::Endpoint(args) => {
            let client = get_client(&cfg?)?;
            args.command.exec(&client, color_mode).await?;
        }
        RootCommands::Message(args) => {
            let client = get_client(&cfg?)?;
            args.command.exec(&client, color_mode).await?;
        }
        // FIXME: need codegen for this one
        RootCommands::MessageAttempt => todo!("Commands::MessageAttempt"),
        RootCommands::Import => todo!("Commands::Import"),
        RootCommands::Export => todo!("Commands::Export"),
        RootCommands::Integration(args) => {
            let client = get_client(&cfg?)?;
            args.command.exec(&client, color_mode).await?;
        }

        RootCommands::Listen => todo!("Commands::Listen"),
        RootCommands::Login => login::prompt()?,
        RootCommands::GenerateCompletions { shell } => completion::generate(&shell)?,
    }

    Ok(())
}

fn get_client(cfg: &Config) -> Result<svix::api::Svix> {
    let token = cfg.auth_token.clone().ok_or_else(|| {
        anyhow::anyhow!("No auth token set. Try running `svix login` to get started.")
    })?;
    let opts = get_client_options(cfg)?;
    Ok(svix::api::Svix::new(token, Some(opts)))
}

fn get_client_options(cfg: &Config) -> Result<svix::api::SvixOptions> {
    Ok(svix::api::SvixOptions {
        debug: false,
        server_url: cfg
            .server_url
            .clone()
            .or_else(|| Some(String::from(DEFAULT_SERVER_URL))),
        timeout: None,
    })
}
