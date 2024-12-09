use clap::{Parser, Subcommand};

const VERSION: &str = env!("CARGO_PKG_VERSION");

#[derive(Parser)]
#[clap(color = concolor_clap::color_choice())]
struct Cli {
    #[command(flatten)]
    color: concolor_clap::Color,
    #[command(subcommand)]
    command: Commands,
}

// N.b Ordering matters here for how clap presents the help.
// FIXME: double-check Go cli. Seems like cobra may sort the items in the help lexigraphically
#[derive(Subcommand)]
enum Commands {
    /// Get the version of the Svix CLI
    Version,
    /// Interactively configure your Svix API credentials
    Login,
    /// List, create & modify applications
    Application(ApplicationArgs),
    /// Manage authentication tasks such as getting dashboard URLs
    Authentication,
    /// List, create & modify event types
    EventType,
    /// List, create & modify endpoints
    Endpoint,
    /// List & create messages
    Message,
    /// List, lookup & resend message attempts
    MessageAttempt,
    /// Verify the signature of a webhook message
    Verify,
    /// Quickly open Svix pages in your browser
    Open,
    /// Forward webhook requests a local url
    Listen,
    /// Import data to your Svix Organization
    Import,
    /// Export data from your Svix Organization
    Export,
    /// List integrations by app id
    Integration,
}

#[derive(Subcommand)]
enum ApplicationArgs {
    /// List current applications
    List,
    /// Creates a new application
    Create,
    /// Get an application by id
    Get,
    /// Update an application by id
    Update,
    /// Deletes an application by id
    Delete,
}

fn main() {
    let cli = Cli::parse();

    match &cli.command {
        // Local-only things
        Commands::Version => println!("{VERSION}"),
        Commands::Verify => todo!("Commands::Verify"),
        Commands::Open => todo!("Commands::Open"),

        // Remote API calls
        Commands::Application(cmd) => match cmd {
            ApplicationArgs::List => todo!("application list"),
            ApplicationArgs::Create => todo!("application create"),
            ApplicationArgs::Get => todo!("application get"),
            ApplicationArgs::Update => todo!("application update"),
            ApplicationArgs::Delete => todo!("application delete"),
        },
        Commands::Authentication => todo!("Commands::Authentication"),
        Commands::EventType => todo!("Commands::EventType"),
        Commands::Endpoint => todo!("Commands::Endpoint"),
        Commands::Message => todo!("Commands::Message"),
        Commands::MessageAttempt => todo!("Commands::MessageAttempt"),
        Commands::Import => todo!("Commands::Import"),
        Commands::Export => todo!("Commands::Export"),
        Commands::Integration => todo!("Commands::Integration"),

        // FIXME: make login/listen play subcommands?
        Commands::Listen => todo!("Commands::Listen"),
        Commands::Login => todo!("Commands::Login"),
    }
}

fn get_client() -> svix::api::Svix {
    // XXX: Go client will exit if the token is not set. May need to rewrangle the flow.
    let token = String::new(); // FIXME: read from config
    let opts = get_client_options();
    svix::api::Svix::new(token, Some(opts))
}

fn get_client_options() -> svix::api::SvixOptions {
    // FIXME: read options from config file
    // FIXME: validate server url
    svix::api::SvixOptions::default()
}
