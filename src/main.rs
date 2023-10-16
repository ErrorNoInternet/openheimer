mod configuration;
mod database;
mod logger;
mod minecraft;

use clap::{ArgAction, Parser};
use configuration::Configuration;
use logger::LogMessageType::*;
use std::str::FromStr;

/// The final word in Minecraft server scanners
#[derive(Parser)]
#[command(author, version, about, long_about = None)]
struct Arguments {
    /// Level of verbosity (can be supplied up to 3 times)
    #[arg(short, long, action = ArgAction::Count)]
    verbosity: u8,
}

fn main() {
    let arguments = Arguments::parse();
    let logger = logger::Logger {
        verbosity: arguments.verbosity as usize,
    };

    logger.log_message(
        Information,
        &format!("Starting OpenHeimer v{}...", env!("CARGO_PKG_VERSION")),
    );

    let configuration_string;
    if std::path::Path::new("openheimer.toml").exists() {
        logger.log_message(Verbose1, "Reading configuration from openheimer.toml...");
        configuration_string = match std::fs::read_to_string("openheimer.toml") {
            Ok(file_data) => file_data,
            Err(error) => {
                logger.log_error("Unable to read openheimer.toml", &error);
                return;
            }
        };
    } else {
        logger.log_message(
            Information,
            "Saving default configuration to openheimer.toml...",
        );
        let default_configuration_string = Configuration::default().to_string();
        configuration_string = default_configuration_string.clone();
        match std::fs::write("openheimer.toml", &default_configuration_string) {
            Ok(_) => (),
            Err(error) => {
                logger.log_error("Unable to write to openheimer.toml", &error);
                return;
            }
        };
    }

    logger.log_message(Verbose3, "Parsing configuration from string...");
    let configuration = match Configuration::from_str(&configuration_string) {
        Ok(configuration) => configuration,
        Err(error) => {
            logger.log_error("Unable to parse configuration", &error);
            return;
        }
    };
    logger.log_message(
        Verbose3,
        &format!("Loaded configuration:\n{configuration:?}"),
    )
}
