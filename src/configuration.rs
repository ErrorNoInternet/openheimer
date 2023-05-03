use std::str::FromStr;

#[derive(Clone)]
pub struct DatabaseConfiguration {
    provider: String,
    location: String,
}

#[derive(Clone)]
pub struct Configuration {
    database: DatabaseConfiguration,
}

impl Default for Configuration {
    fn default() -> Self {
        Configuration {
            database: DatabaseConfiguration {
                provider: "sqlite".to_string(),
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
