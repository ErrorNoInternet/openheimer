use super::Database;
use crate::minecraft::{player::Player, server::Server};
use async_trait::async_trait;
use sqlx::sqlite::SqlitePoolOptions;

pub struct SqliteDatabase {
    database: sqlx::SqlitePool,
}

#[async_trait]
impl Database for SqliteDatabase {
    async fn new(location: String) -> Result<Self, String> {
        let database = match SqlitePoolOptions::new()
            .connect(format!("sqlite://{location}?mode=rwc").as_str())
            .await
        {
            Ok(database) => database,
            Err(error) => return Err(error.to_string()),
        };
        Ok(SqliteDatabase { database })
    }

    async fn add_server(server: Server) -> Result<(), String> {
        // TODO: Add server
        Ok(())
    }

    async fn add_player(player: Player) -> Result<(), String> {
        // TODO: Add player
        Ok(())
    }

    async fn add_targeted_server(server: Server) -> Result<(), String> {
        // TODO: Add server
        Ok(())
    }

    async fn add_targeted_player(player: Player) -> Result<(), String> {
        // TODO: Add player
        Ok(())
    }
}
