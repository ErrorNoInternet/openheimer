use serde::Serialize;

#[derive(Clone, Default, Serialize)]
pub struct Player {
    first_seen: u64,
    last_seen: u64,

    name: String,
    uuid: String,
}

impl Player {}
