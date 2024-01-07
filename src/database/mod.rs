pub mod sqlite;

use crate::minecraft;
use async_trait::async_trait;

pub enum DatabaseError {
    SqlxPoolCreateFailed(sqlx::Error),
    SqlxAcquirePoolFailed(sqlx::Error),
    SqlxInsertQueryFailed(sqlx::Error),
}

#[async_trait]
pub trait Database {
    async fn new(location: String) -> Result<Self, DatabaseError>
    where
        Self: Sized;
    async fn add_server(&mut self, server: &minecraft::server::Server)
        -> Result<(), DatabaseError>;
    async fn add_player(&mut self, player: &minecraft::player::Player)
        -> Result<(), DatabaseError>;
    async fn add_targeted_server(
        &mut self,
        server: &minecraft::server::Server,
    ) -> Result<(), DatabaseError>;
    async fn add_targeted_player(
        &mut self,
        player: &minecraft::player::Player,
    ) -> Result<(), DatabaseError>;
}
