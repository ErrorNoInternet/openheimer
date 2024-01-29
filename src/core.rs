use crate::{arguments::Arguments, configuration::Configuration, metadata};
use std::{process::exit, str::FromStr};
use tracing::{debug, info, trace, warn};
use tracing_appender::rolling::RollingFileAppender;
use tracing_subscriber::prelude::*;

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
        Ok(file_appender) => set_up_logging(file_appender),
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
    if let Some(configuration_file) = &arguments.configuration_file {
        trace!("reading configuration file...");
        let file_contents = match std::fs::read_to_string(configuration_file) {
            Ok(file_contents) => file_contents,
            Err(error) => {
                eprintln!("unable to read configuration file: {error:#?}\n");
                exit(1);
            }
        };

        trace!("parsing configuration file...");
        match Configuration::from_str(file_contents.as_str()) {
            Ok(options) => (false, options),
            Err(error) => {
                eprintln!("unable to read configuration file: {error:#?}\n");
                exit(1);
            }
        }
    } else {
        (true, Configuration::default())
    }
}

fn set_up_logging(file_appender: RollingFileAppender) {
    let file_layer = tracing_subscriber::fmt::layer()
        .with_file(true)
        .with_line_number(true)
        .with_writer(file_appender);
    let console_layer = tracing_subscriber::fmt::layer()
        .with_file(true)
        .with_line_number(true);

    let subscriber = tracing_subscriber::Registry::default()
        .with(file_layer)
        .with(console_layer);
    match tracing::subscriber::set_global_default(subscriber) {
        Ok(()) => (),
        Err(error) => eprintln!("unable to set up logging: {error:#?}"),
    };
}

#[cfg(test)]
mod test {
    use super::set_up_logging;
    use crate::configuration::Configuration;
    use tracing::{debug, error, info, trace, warn};

    #[test]
    fn log_messages() {
        let options = Configuration::default();
        set_up_logging(
            tracing_appender::rolling::Builder::new()
                .filename_prefix(options.logger.prefix)
                .filename_suffix(options.logger.suffix)
                .rotation(options.logger.rotation.clone().into())
                .max_log_files(options.logger.max_log_files)
                .build(options.logger.directory)
                .expect("should have been able to create logger files"),
        );

        error!("h");
        warn!("e");
        info!("l");
        debug!("l");
        trace!("o");
    }
}
