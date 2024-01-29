pub mod sqlite;

use crate::minecraft;
use async_trait::async_trait;

pub enum Error {
    Create(sqlx::Error),
    Acquire(sqlx::Error),
    Insert(sqlx::Error),
    Serialize(serde_json::Error),
}

#[async_trait]
pub trait Database {
    async fn new(location: String) -> Result<Self, Error>
    where
        Self: Sized;

    async fn add_server(&mut self, server: &minecraft::server::Server) -> Result<(), Error>;

    async fn add_player(&mut self, player: &minecraft::player::Player) -> Result<(), Error>;

    async fn add_targeted_server(
        &mut self,
        server: &minecraft::server::Server,
    ) -> Result<(), Error>;

    async fn add_targeted_player(
        &mut self,
        player: &minecraft::player::Player,
    ) -> Result<(), Error>;
}
