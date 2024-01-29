use super::{Arguments, ConfigurationSubcommand};
use crate::configuration::Configuration;
use std::str::FromStr;

pub fn parse(arguments: &Arguments, subcommand: &ConfigurationSubcommand) {
    match subcommand {
        ConfigurationSubcommand::Default => default(),
        ConfigurationSubcommand::Fill => fill(arguments),
    }
}

fn default() {
    println!("{}", Configuration::default().to_string());
}

fn fill(arguments: &Arguments) {
    let default = Configuration::default();
    let options = if let Some(configuration_file) = &arguments.configuration_file {
        let file_contents = match std::fs::read_to_string(configuration_file) {
            Ok(file_contents) => file_contents,
            Err(error) => {
                eprintln!("unable to read configuration file: {error:#?}");
                default.to_string()
            }
        };

        match Configuration::from_str(file_contents.as_str()) {
            Ok(options) => options,
            Err(error) => {
                eprintln!("unable to parse configuration file: {error:#?}");
                default
            }
        }
    } else {
        default
    };

    println!("{}", options.to_string());
}
