#[derive(Clone, Default)]
pub struct Server {
    address: String,
    offline_mode: bool,
    whitelisted: bool,
    first_seen: u64,
    last_seen: u64,
    extra_data_json: String,

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
    enforces_secure_chat: bool,
}

impl Server {}
