use clap::Parser;
use tracing::level_filters::LevelFilter;

/// The final word in Minecraft server scanners
#[derive(Debug, Parser)]
#[command(author, version, about, long_about = None)]
pub struct Arguments {
    /// How much information to log
    #[arg(short, long, default_value = "info")]
    pub verbosity: LevelFilter,

    /// openheimer toml configuration file path
    #[arg(short, long)]
    pub configuration_file: Option<String>,
}
