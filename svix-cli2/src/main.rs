use clap::{Parser, Subcommand};

#[derive(Parser)]
#[clap(color = concolor_clap::color_choice())]
struct Cli {
    #[command(flatten)]
    color: concolor_clap::Color,
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    Version,
    Login,
    Application,
    Authentication,
    EventType,
    Endpoint,
    Message,
    MessageAttempt,
    Verify,
    Open,
    Listen,
    Import,
    Export,
    Integration,
}

fn main() {
    let cli = Cli::parse();

    match &cli.command {
        // Local-only things
        Commands::Version => todo!("Commands::Version"),
        Commands::Verify => todo!("Commands::Verify"),
        Commands::Open => todo!("Commands::Open"),

        // Remote API calls
        Commands::Application => todo!("Commands::Application"),
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
