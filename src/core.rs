use crate::{arguments::Arguments, configuration::Configuration, logging, metadata};
use std::str::FromStr;
use tracing::{debug, info, trace, warn};

pub fn main(arguments: &Arguments) {
    trace!("getting configuration options...");
    let (is_default, options) = get_options(arguments);

    trace!("setting up logger...");
    logging::set_up_logging(arguments.verbosity, options.logger.clone());

    info!("openheimer {}", metadata::format());
    if is_default {
        warn!("falling back to default configuration!");
    }
    debug!("using configuration: {options:#?}");
}

fn get_options(arguments: &Arguments) -> (bool, Configuration) {
    trace!("reading configuration file...");
    let file_contents = match std::fs::read_to_string(arguments.configuration_file.clone()) {
        Ok(file_contents) => file_contents,
        Err(error) => {
            warn!("unable to read configuration file: {error:#?}");
            return (true, Configuration::default());
        }
    };

    trace!("parsing configuration file...");
    match Configuration::from_str(file_contents.as_str()) {
        Ok(options) => (false, options),
        Err(error) => {
            warn!("unable to parse configuration file: {error:#?}");
            (true, Configuration::default())
        }
    }
}
