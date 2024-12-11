use crate::cli_types::Ordering;
use clap::Args;
use svix::api;

#[derive(Args, Clone)]
pub struct EndpointListOptions {
    #[arg(long)]
    pub iterator: Option<String>,
    #[arg(long)]
    pub limit: Option<i32>,
    #[arg(long)]
    pub order: Option<Ordering>,
}

impl From<EndpointListOptions> for api::EndpointListOptions {
    fn from(
        EndpointListOptions {
            iterator,
            limit,
            order,
        }: EndpointListOptions,
    ) -> Self {
        Self {
            iterator,
            limit,
            order: order.map(Into::into),
        }
    }
}
