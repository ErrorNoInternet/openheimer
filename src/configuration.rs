use crate::logger;
use logger::LogMessageType::ConfigurationWarning;
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

    fn from_str(string: &str) -> Result<Self, String> {
        let logger = logger::Logger { verbosity: 0 };

        let configuration = string.parse::<toml::Table>().unwrap();
        let default_configuration = toml::Table::try_from(Configuration::default()).unwrap();

        let get_default_value = |tree: &[&str]| {
            let mut default_branch = default_configuration.get(tree[0]).unwrap();
            for branch_index in 1..tree.len() {
                default_branch = default_branch.get(tree[branch_index]).unwrap();
            }
            default_branch.clone()
        };

        let get_value = |tree: &[&str]| {
            let mut default_branch = default_configuration.get(tree[0]).unwrap();
            let mut branch = match configuration.get(tree[0]) {
                Some(branch) => branch,
                None => {
                    logger.log_message(
                        ConfigurationWarning,
                        &format!("Unable to find `{}`, using default!", tree[0]),
                    );
                    default_branch
                }
            };

            for branch_index in 1..tree.len() {
                default_branch = default_branch.get(tree[branch_index]).unwrap();
                branch = match branch.get(tree[branch_index]) {
                    Some(branch) => branch,
                    None => {
                        logger.log_message(
                            ConfigurationWarning,
                            &format!(
                                "Unable to find `{}`, using default!",
                                tree.split_at(branch_index + 1).0.join(".")
                            ),
                        );
                        default_branch
                    }
                };
            }

            branch.clone()
        };

        let get_string_value = |tree: &[&str]| match get_value(tree).as_str() {
            Some(value) => value.to_string(),
            None => {
                logger.log_message(
                    ConfigurationWarning,
                    &format!("`{}` is not a string, using default!", tree.join(".")),
                );
                get_default_value(tree).as_str().unwrap().to_string()
            }
        };

        Ok(Self {
            database: DatabaseConfiguration {
                provider: get_string_value(&["database", "provider"]),
                location: get_string_value(&["database", "location"]),
            },
        })
    }
}

impl ToString for Configuration {
    fn to_string(&self) -> String {
        toml::to_string_pretty(&self).unwrap()
    }
}
