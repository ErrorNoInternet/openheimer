pub mod configuration;

use crate::metadata;
use clap::{Parser, Subcommand};
use tracing::level_filters::LevelFilter;

/// The final word in Minecraft server scanners
#[derive(Debug, Parser)]
#[command(version = metadata::format(), about, subcommand_required = true)]
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
    Start,

    /// Manipulate openheimer configuration files
    Configuration {
        #[command(subcommand)]
        subcommand: ConfigurationSubcommand,
    },
}

#[derive(Debug, Subcommand)]
pub enum ConfigurationSubcommand {
    /// Generate a default configuration file
    Default,

    /// Add missing options to your existing configuration file
    Fill,
}
