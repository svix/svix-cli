//! Wrapper types around those in the SDK so they can act as args in the `clap` command tree.
// FIXME: in the new codegen, we can look to add a `cli` feature that conditionally adds these
//   attributes or trait impls to the wrapped types.
//   At that point we can use the SDK types directly and remove this module.

use clap::{Args, ValueEnum};
use svix::api;

pub mod application;

#[derive(Copy, Clone, PartialEq, Eq, PartialOrd, Ord, ValueEnum)]
pub enum Ordering {
    Ascending,
    Descending,
}

impl From<Ordering> for api::Ordering {
    fn from(value: Ordering) -> Self {
        match value {
            Ordering::Ascending => api::Ordering::Ascending,
            Ordering::Descending => api::Ordering::Descending,
        }
    }
}
