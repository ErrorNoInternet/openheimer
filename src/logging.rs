use tracing_appender::rolling::RollingFileAppender;
use tracing_subscriber::prelude::*;

pub fn set_up_logging(file_appender: RollingFileAppender) {
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
