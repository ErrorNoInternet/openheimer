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
    #[arg(short, long, default_value = "openheimer.toml")]
    pub configuration_file: String,

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
    /// Generate a new configuration file template
    Default,

    /// Add missing fields to an existing configuration file
    Fill,

    /// Check the validity of an existing configuration file
    Validate,
}
