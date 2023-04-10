#[derive(Clone, Default)]
pub struct Server {
    version_name: String,
    protocol_version: i32,
    favicon: String,
    description: String,
    description_json: String,
    players_online: i32,
    players_max: i32,
    players_json: String,
    mod_info_type: String,
    mod_info_list_json: String,
    offline_mode: bool,
    whitelisted: bool,
    extra_data_json: String,
}

impl Server {
    fn from_json(json: String) -> Result<Self, String> {
        // TODO: Parse JSON
        Ok(Server::default())
    }
}
