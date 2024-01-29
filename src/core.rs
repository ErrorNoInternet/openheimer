use crate::{arguments::Arguments, configuration::Configuration, logging, metadata};
use std::str::FromStr;
use tracing::{debug, info, trace, warn};

pub fn main(arguments: &Arguments) {
    let (is_default, options) = get_options(arguments);

    let binding = options.clone();
    match tracing_appender::rolling::Builder::new()
        .filename_prefix(binding.logger.prefix)
        .filename_suffix(binding.logger.suffix)
        .rotation(options.logger.rotation.clone().into())
        .max_log_files(options.logger.max_log_files)
        .build(binding.logger.directory)
    {
        Ok(file_appender) => logging::set_up_logging(file_appender),
        Err(error) => {
            eprintln!("unable to set up rolling logger: {error:#?}");
        }
    };

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
