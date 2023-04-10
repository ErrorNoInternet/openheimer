#[derive(Clone, Default)]
pub struct Player {
    name: String,
    uuid: String,
}

impl Player {
    fn from_json(json: String) -> Result<Self, String> {
        // TODO: Parse JSON
        Ok(Player::default())
    }
}
