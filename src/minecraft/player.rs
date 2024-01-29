use serde::Serialize;

#[derive(Clone, Default, Serialize)]
pub struct Player {
    pub first_seen: u64,
    pub last_seen: u64,

    pub name: String,
    pub uuid: String,
}
