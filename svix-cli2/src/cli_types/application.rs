use crate::cli_types::Ordering;
use clap::Args;
use svix::api;

#[derive(Args, Clone)]
pub struct ApplicationListOptions {
    #[arg(long)]
    pub iterator: Option<String>,
    #[arg(long)]
    pub limit: Option<i32>,
    #[arg(long)]
    pub order: Option<Ordering>,
}

impl From<ApplicationListOptions> for api::ApplicationListOptions {
    fn from(
        ApplicationListOptions {
            iterator,
            limit,
            order,
        }: ApplicationListOptions,
    ) -> Self {
        Self {
            iterator,
            limit,
            order: order.map(Into::into),
        }
    }
}
