use serde::Serialize;

#[derive(Clone, Default, Serialize)]
pub struct Server {
    pub address: String,
    pub offline_mode: bool,
    pub whitelisted: bool,
    pub first_seen: u64,
    pub last_seen: u64,
    pub extra_data_json: String,

    pub version_name: String,
    pub protocol_version: i32,
    pub favicon: String,
    pub description: String,
    pub description_json: String,
    pub players_online: i32,
    pub players_max: i32,
    pub players_json: String,
    pub mod_info_type: String,
    pub mod_info_list_json: String,
    pub enforces_secure_chat: bool,
}
