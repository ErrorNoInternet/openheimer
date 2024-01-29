use serde::{Deserialize, Serialize};
use std::str::FromStr;

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
    pub prefix: String,
}

impl Default for Logger {
    fn default() -> Self {
        Self {
            directory: "openheimer/logs".into(),
            prefix: "openheimer.log".into(),
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
        toml::to_string_pretty(&self).unwrap()
    }
}
