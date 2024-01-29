pub mod arguments;
pub mod configuration;
pub mod core;
pub mod database;
pub mod metadata;
pub mod minecraft;

use crate::arguments::MainSubcommand;
use arguments::Arguments;
use clap::Parser;

pub fn main() {
    let arguments = Arguments::parse();
    if let Some(ref subcommand) = arguments.subcommand {
        match subcommand {
            MainSubcommand::Start => core::main(&arguments),
            MainSubcommand::Configuration { subcommand } => {
                arguments::configuration::parse(&arguments, subcommand);
            }
        };
    }
}
