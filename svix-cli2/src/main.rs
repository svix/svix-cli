use anyhow::Result;
use clap::{Args, Parser, Subcommand, ValueEnum};
use colored_json::{ColorMode, Output};
use concolor_clap::{Color, ColorChoice};
use serde::Serialize;
use svix::api::{ApplicationListOptions, Ordering};

const VERSION: &str = env!("CARGO_PKG_VERSION");

#[derive(Parser)]
#[clap(color = concolor_clap::color_choice())]
struct Cli {
    #[command(flatten)]
    color: Color,
    #[command(subcommand)]
    command: RootCommands,
}

impl Cli {
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
    Authentication,
    /// List, create & modify event types
    Endpoint,
    /// List & create messages
    EventType,
    /// List, create & modify endpoints
    /// Export data from your Svix Organization
    Export,
    /// Import data to your Svix Organization
    Import,
    /// List integrations by app id
    Integration,
    /// Forward webhook requests a local url
    Listen,
    /// Interactively configure your Svix API credentials
    Login,
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

#[derive(Args)]
#[command(args_conflicts_with_subcommands = true)]
#[command(flatten_help = true)]
struct ApplicationArgs {
    #[command(subcommand)]
    command: ApplicationCommands,
}

#[derive(Copy, Clone, PartialEq, Eq, PartialOrd, Ord, ValueEnum)]
pub enum Ordering_ {
    Ascending,
    Descending,
}

impl From<Ordering_> for Ordering {
    fn from(value: Ordering_) -> Self {
        match value {
            Ordering_::Ascending => Ordering::Ascending,
            Ordering_::Descending => Ordering::Descending,
        }
    }
}

// FIXME: there's surely a way to use the type as-is from the libs with clap.
//  Workaround: make this type and a From impl to convert for lib usage.
//  Updating codegen to derive `Args`
#[derive(Args, Clone)]
pub struct ApplicationListOptions_ {
    #[arg(long)]
    pub iterator: Option<String>,
    #[arg(long)]
    pub limit: Option<i32>,
    #[arg(long)]
    pub order: Option<Ordering_>,
}

impl From<ApplicationListOptions_> for ApplicationListOptions {
    fn from(
        ApplicationListOptions_ {
            iterator,
            limit,
            order,
        }: ApplicationListOptions_,
    ) -> Self {
        ApplicationListOptions {
            iterator,
            limit,
            order: order.map(Into::into),
        }
    }
}

// FIXME: build these via codegen from the spec, along with the rust lib.
#[derive(Subcommand)]
enum ApplicationCommands {
    /// List current applications
    List(ApplicationListOptions_),
    /// Creates a new application
    Create { body: String },
    /// Get an application by id
    Get { id: String },
    /// Update an application by id
    Update { id: String, body: String },
    /// Deletes an application by id
    Delete { id: String },
}

fn print_json_output<T>(val: &T, color_mode: ColorMode) -> Result<()>
where
    T: Serialize,
{
    // FIXME: factor the writer out? Will that help with testing?
    let mut writer = std::io::stdout().lock();
    colored_json::write_colored_json_with_mode(val, &mut writer, color_mode)?;
    Ok(())
}

impl ApplicationCommands {
    // FIXME: codegen an exec() method that takes the args and a client and does the thing?
    //   Not sure if we need to pass in a printer or how the output should work if we can't
    //   have a typed return here.
    //   This might not make sense but let's roll with it for now.
    async fn exec(&self, client: &svix::api::Svix, color_mode: ColorMode) -> Result<()> {
        match self {
            ApplicationCommands::List(options) => {
                let resp = client
                    .application()
                    .list(Some(options.clone().into()))
                    .await?;

                print_json_output(&resp, color_mode)?;
            }
            ApplicationCommands::Create { body } => todo!("application create"),
            ApplicationCommands::Get { id } => todo!("application get"),
            ApplicationCommands::Update { id, body } => todo!("application update"),
            ApplicationCommands::Delete { id } => todo!("application delete"),
        }
        Ok(())
    }
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
        RootCommands::Authentication => todo!("Commands::Authentication"),
        RootCommands::EventType => todo!("Commands::EventType"),
        RootCommands::Endpoint => todo!("Commands::Endpoint"),
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
