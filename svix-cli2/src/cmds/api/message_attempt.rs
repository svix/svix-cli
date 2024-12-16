use clap::{Args, Subcommand};
use colored_json::ColorMode;

use crate::cli_types::message_attempt::{
    MessageAttemptListAttemptedDestinationsOptions, MessageAttemptListByEndpointOptions,
    MessageAttemptListByMsgOptions,
};

#[derive(Args)]
#[command(args_conflicts_with_subcommands = true)]
#[command(flatten_help = true)]
pub struct MessageAttemptArgs {
    #[command(subcommand)]
    pub command: MessageAttemptCommands,
}

#[derive(Subcommand)]
pub enum MessageAttemptCommands {
    /// List attempts by endpoint id
    ///
    /// Note that by default this endpoint is limited to retrieving 90 days' worth of data
    /// relative to now or, if an iterator is provided, 90 days before/after the time indicated
    /// by the iterator ID. If you require data beyond those time ranges, you will need to explicitly
    /// set the `before` or `after` parameter as appropriate.
    ListByEndpoint {
        app_id: String,
        endpoint_id: String,

        #[clap(flatten)]
        options: MessageAttemptListByEndpointOptions,
    },
    /// List attempts by message id
    ///
    /// Note that by default this endpoint is limited to retrieving 90 days' worth of data
    /// relative to now or, if an iterator is provided, 90 days before/after the time indicated
    /// by the iterator ID. If you require data beyond those time ranges, you will need to explicitly
    /// set the `before` or `after` parameter as appropriate.
    ListByMsg {
        app_id: String,
        msg_id: String,

        #[clap(flatten)]
        options: MessageAttemptListByMsgOptions,
    },
    /// `msg_id`: Use a message id or a message `eventId`
    Get {
        app_id: String,
        msg_id: String,
        attempt_id: String,
    },
    /// Deletes the given attempt's response body. Useful when an endpoint accidentally returned sensitive content.
    ExpungeContent {
        app_id: String,
        msg_id: String,
        attempt_id: String,
    },
    /// List endpoints attempted by a given message. Additionally includes metadata about the latest message attempt.
    /// By default, endpoints are listed in ascending order by ID.
    ListAttemptedDestinations {
        app_id: String,
        msg_id: String,

        #[clap(flatten)]
        options: MessageAttemptListAttemptedDestinationsOptions,
    },
    /// Resend a message to the specified endpoint.
    Resend {
        app_id: String,
        msg_id: String,
        endpoint_id: String,
        // FIXME: Not currently supported by the Rust lib
        //#[clap(flatten)]
        //post_options: Option<PostOptions>,
    },
}

impl MessageAttemptCommands {
    pub async fn exec(self, client: &svix::api::Svix, color_mode: ColorMode) -> anyhow::Result<()> {
        match self {
            Self::ListByEndpoint {
                app_id,
                endpoint_id,
                options,
            } => {
                let resp = client
                    .message_attempt()
                    .list_by_endpoint(app_id, endpoint_id, Some(options.into()))
                    .await?;
                crate::json::print_json_output(&resp, color_mode)?;
            }
            Self::ListByMsg {
                app_id,
                msg_id,
                options,
            } => {
                let resp = client
                    .message_attempt()
                    .list_by_msg(app_id, msg_id, Some(options.into()))
                    .await?;
                crate::json::print_json_output(&resp, color_mode)?;
            }
            Self::Get {
                app_id,
                msg_id,
                attempt_id,
            } => {
                let resp = client
                    .message_attempt()
                    .get(app_id, msg_id, attempt_id)
                    .await?;
                crate::json::print_json_output(&resp, color_mode)?;
            }
            Self::ExpungeContent {
                app_id,
                msg_id,
                attempt_id,
            } => {
                client
                    .message_attempt()
                    .expunge_content(app_id, msg_id, attempt_id)
                    .await?;
            }
            Self::ListAttemptedDestinations {
                app_id,
                msg_id,
                options,
            } => {
                let resp = client
                    .message_attempt()
                    .list_attempted_destinations(app_id, msg_id, Some(options.into()))
                    .await?;
                crate::json::print_json_output(&resp, color_mode)?;
            }
            Self::Resend {
                app_id,
                msg_id,
                endpoint_id,
            } => {
                client
                    .message_attempt()
                    .resend(app_id, msg_id, endpoint_id)
                    .await?;
            }
        }

        Ok(())
    }
}
