use super::{Arguments, ConfigurationSubcommand};
use crate::configuration::Configuration;
use std::str::FromStr;

pub fn parse(arguments: &Arguments, subcommand: &ConfigurationSubcommand) {
    match subcommand {
        ConfigurationSubcommand::Default => default(),
        ConfigurationSubcommand::Fill => fill(arguments),
        ConfigurationSubcommand::Validate => validate(arguments),
    }
}

fn default() {
    println!("{}", Configuration::default());
}

fn fill(arguments: &Arguments) {
    let default = Configuration::default();

    let file_contents = match std::fs::read_to_string(arguments.configuration_file.clone()) {
        Ok(file_contents) => file_contents,
        Err(error) => {
            eprintln!("unable to read configuration file: {error:#?}");
            default.to_string()
        }
    };

    let options = match Configuration::from_str(file_contents.as_str()) {
        Ok(options) => options,
        Err(error) => {
            eprintln!("unable to parse configuration file: {error:#?}");
            default
        }
    };

    println!("{options}");
}

fn validate(arguments: &Arguments) {
    match std::fs::read_to_string(arguments.configuration_file.clone()) {
        Ok(file_contents) => match Configuration::from_str(file_contents.as_str()) {
            Ok(_) => println!("no errors found"),
            Err(error) => {
                eprintln!("unable to parse configuration file: {error:#?}");
            }
        },
        Err(error) => {
            eprintln!("unable to read configuration file: {error:#?}");
        }
    };
}
