use serde::Serialize;
use std::str::FromStr;

#[derive(Clone, Serialize)]
pub struct DatabaseConfiguration {
    provider: String,
    location: String,
}

#[derive(Clone, Serialize)]
pub struct Configuration {
    database: DatabaseConfiguration,
}

impl Default for Configuration {
    fn default() -> Self {
        Configuration {
            database: DatabaseConfiguration {
                provider: "sqlite3".to_string(),
                location: "openheimer.db".to_string(),
            },
        }
    }
}

impl FromStr for Configuration {
    type Err = String;

    fn from_str(string: &str) -> Result<Self, Self::Err> {
        Ok(Default::default())
    }
}

impl ToString for Configuration {
    fn to_string(&self) -> String {
        toml::to_string_pretty(&self).unwrap()
    }
}