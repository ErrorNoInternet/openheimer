use serde::{Deserialize, Serialize};
use std::str::FromStr;
use tracing_appender::rolling;

#[derive(Clone, Debug, Default, Deserialize, Serialize)]
#[serde(default)]
pub struct Configuration {
    pub logger: Logger,
    pub database: Database,
}

#[derive(Clone, Debug, Deserialize, Serialize)]
#[serde(default)]
pub struct Logger {
    pub directory: String,
    pub rotation: LoggerRotation,
    pub max_log_files: usize,
    pub prefix: String,
    pub suffix: String,
}

impl Default for Logger {
    fn default() -> Self {
        Self {
            directory: "openheimer/logs".into(),
            rotation: LoggerRotation::default(),
            max_log_files: 100,
            prefix: "openheimer".into(),
            suffix: "log".into(),
        }
    }
}

#[derive(Clone, Debug, Default, Deserialize, Serialize)]
pub enum LoggerRotation {
    Never,
    Minutely,
    Hourly,
    #[default]
    Daily,
}

impl From<LoggerRotation> for tracing_appender::rolling::Rotation {
    fn from(value: LoggerRotation) -> Self {
        match value {
            LoggerRotation::Never => rolling::Rotation::NEVER,
            LoggerRotation::Minutely => rolling::Rotation::MINUTELY,
            LoggerRotation::Hourly => rolling::Rotation::HOURLY,
            LoggerRotation::Daily => rolling::Rotation::DAILY,
        }
    }
}

#[derive(Clone, Debug, Deserialize, Serialize)]
#[serde(default)]
pub struct Database {
    provider: String,
    location: String,
}

impl Default for Database {
    fn default() -> Self {
        Self {
            provider: "sqlite".into(),
            location: "openheimer.db".into(),
        }
    }
}

impl FromStr for Configuration {
    type Err = String;

    fn from_str(string: &str) -> Result<Self, Self::Err> {
        match toml::from_str(string) {
            Ok(configuration) => Ok(configuration),
            Err(error) => Err(format!("deserialization error: {error:?}")),
        }
    }
}

impl ToString for Configuration {
    fn to_string(&self) -> String {
        toml::to_string_pretty(&self).unwrap().trim().to_string()
    }
}
