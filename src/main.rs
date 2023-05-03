mod configuration;
mod database;
mod logger;
mod minecraft;

use clap::{Arg, ArgAction, Command};
use logger::LogMessageType::*;

fn main() {
    let matches = Command::new("openheimer")
        .about("The final word in Minecraft server scanners")
        .arg(
            Arg::new("verbosity")
                .short('v')
                .long("verbose")
                .help("Sets the level of verbosity")
                .action(ArgAction::Count),
        )
        .get_matches();

    let logger = logger::Logger {
        verbosity: matches.get_count("verbosity") as usize,
    };

    logger.log_message(
        Information,
        &format!("Starting OpenHeimer v{}...", env!("CARGO_PKG_VERSION")),
    );
}
