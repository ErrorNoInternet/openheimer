pub mod sqlite;

use crate::minecraft;
use async_trait::async_trait;

#[async_trait]
pub trait Database {
    async fn new(location: String) -> Result<Self, String>
    where
        Self: Sized;
    async fn add_server(server: minecraft::server::Server) -> Result<(), String>;
    async fn add_player(player: minecraft::player::Player) -> Result<(), String>;
    async fn add_targeted_server(server: minecraft::server::Server) -> Result<(), String>;
    async fn add_targeted_player(player: minecraft::player::Player) -> Result<(), String>;

    // TODO: Add CLI configuration options
}
