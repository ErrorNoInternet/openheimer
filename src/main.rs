mod configuration;
mod database;
mod logger;
mod minecraft;

use logger::log_message;
use logger::LogMessageType::*;

fn main() {
    log_message(
        Information,
        &format!("Starting OpenHeimer v{}...", env!("CARGO_PKG_VERSION")),
    )
}
