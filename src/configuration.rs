use serde::{Deserialize, Serialize};
use std::str::FromStr;
use toml;

#[derive(Debug, Deserialize, Serialize)]
#[serde(default)]
pub struct DatabaseConfiguration {
    provider: String,
    location: String,
}

impl Default for DatabaseConfiguration {
    fn default() -> Self {
        Self {
            provider: "sqlite3".to_string(),
            location: "openheimer.db".to_string(),
        }
    }
}

#[derive(Debug, Default, Deserialize, Serialize)]
#[serde(default)]
pub struct Configuration {
    database: DatabaseConfiguration,
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
