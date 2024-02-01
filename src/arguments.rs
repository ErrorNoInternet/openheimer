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
    #[command(visible_aliases = ["s"])]
    Start,

    /// Manipulate openheimer configuration files
    #[command(visible_aliases = ["c", "config"])]
    Configuration {
        #[command(subcommand)]
        subcommand: ConfigurationSubcommand,
    },
}

#[derive(Debug, Subcommand)]
pub enum ConfigurationSubcommand {
    /// Generate a new configuration file template
    #[command(visible_aliases = ["d", "def", "new"])]
    Default,

    /// Add missing fields to an existing configuration file
    #[command(visible_aliases = ["f", "merge"])]
    Fill,

    /// Check the validity of an existing configuration file
    #[command(visible_aliases = ["v", "val", "check"])]
    Validate,
}
