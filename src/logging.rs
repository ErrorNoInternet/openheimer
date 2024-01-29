use crate::configuration;
use tracing::level_filters::LevelFilter;
use tracing_subscriber::prelude::*;

pub fn init(verbosity: LevelFilter, logger_options: configuration::Logger) {
    let console_layer = tracing_subscriber::fmt::layer()
        .with_file(true)
        .with_line_number(true)
        .with_writer(std::io::stdout);

    match tracing_appender::rolling::Builder::new()
        .filename_prefix(logger_options.prefix)
        .filename_suffix(logger_options.suffix)
        .rotation(logger_options.rotation.clone().into())
        .max_log_files(logger_options.max_log_files)
        .build(logger_options.directory)
    {
        Ok(file_writer) => {
            let file_layer = tracing_subscriber::fmt::layer()
                .with_file(true)
                .with_line_number(true)
                .with_writer(file_writer);

            match tracing::subscriber::set_global_default(
                tracing_subscriber::Registry::default()
                    .with(verbosity)
                    .with(console_layer)
                    .with(file_layer),
            ) {
                Ok(()) => (),
                Err(error) => eprintln!("unable to set up logging: {error:#?}"),
            };
        }
        Err(error) => {
            eprintln!("unable to set up rolling file logger: {error:#?}");

            match tracing::subscriber::set_global_default(
                tracing_subscriber::Registry::default()
                    .with(verbosity)
                    .with(console_layer),
            ) {
                Ok(()) => (),
                Err(error) => eprintln!("unable to set up logging: {error:#?}"),
            };
        }
    };
}

#[cfg(test)]
mod test {
    use super::init;
    use crate::configuration::Configuration;
    use tracing::{debug, error, info, level_filters::LevelFilter, trace, warn};

    #[test]
    fn log_messages() {
        init(LevelFilter::TRACE, Configuration::default().logger);

        error!("h");
        warn!("e");
        info!("l");
        debug!("l");
        trace!("o");
    }
}
