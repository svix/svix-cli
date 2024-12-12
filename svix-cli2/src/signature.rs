use clap::{Args, Subcommand};
use colored_json::ColorMode;

#[derive(Args)]
#[command(args_conflicts_with_subcommands = true)]
#[command(flatten_help = true)]
pub struct SignatureArgs {
    #[command(subcommand)]
    pub command: SignatureCommands,
}

#[derive(Subcommand)]
pub enum SignatureCommands {
    Sign,
    /// Verify the signature of a webhook message
    Verify,
}

impl SignatureCommands {
    pub async fn exec(self, _color_mode: ColorMode) -> anyhow::Result<()> {
        match self {
            SignatureCommands::Sign => todo!("SignatureCommands::Sign"),
            SignatureCommands::Verify => todo!("SignatureCommands::Verify"),
        }
        Ok(())
    }
}
