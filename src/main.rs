mod arguments;
mod configuration;
mod core;
mod database;
mod metadata;
mod minecraft;

use crate::arguments::MainSubcommand;
use arguments::Arguments;
use clap::Parser;

fn main() {
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
