mod configuration;
mod database;
mod logger;
mod minecraft;

use clap::{Arg, ArgAction, Command};
use configuration::Configuration;
use logger::LogMessageType::*;
use std::str::FromStr;

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
            Verbose1,
            "Saving default configuration to openheimer.toml...",
        );
        let defualt_configuration_string = Configuration::default().to_string();
        configuration_string = defualt_configuration_string.clone();
        match std::fs::write("openheimer.toml", &defualt_configuration_string) {
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
        &format!("Loaded configuration:\n{}", configuration.to_string()),
    )
}
