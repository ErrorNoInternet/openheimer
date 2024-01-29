pub mod configuration;

use crate::metadata;
use clap::{Parser, Subcommand};
use tracing::level_filters::LevelFilter;

/// The final word in Minecraft server scanners
#[derive(Debug, Parser)]
#[command(author, version = metadata::format(), about, long_about = None)]
pub struct Arguments {
    /// Amount of information to log
    #[arg(short, long, default_value = "info")]
    pub verbosity: LevelFilter,

    /// TOML configuration file path
    #[arg(short, long)]
    pub configuration_file: Option<String>,

    #[command(subcommand)]
    pub subcommand: Option<MainSubcommand>,
}

#[derive(Debug, Subcommand)]
pub enum MainSubcommand {
    Configuration {
        #[command(subcommand)]
        subcommand: ConfigurationSubcommand,
    },
}

#[derive(Debug, Subcommand)]
pub enum ConfigurationSubcommand {
    Default,
    Fill,
}
