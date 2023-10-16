use super::{Database, DatabaseError};
use crate::minecraft::{player::Player, server::Server};
use async_trait::async_trait;
use sqlx::sqlite::SqlitePoolOptions;

pub struct SqliteDatabase {
    pool: sqlx::SqlitePool,
}

#[async_trait]
impl Database for SqliteDatabase {
    async fn new(location: String) -> Result<Self, DatabaseError> {
        let pool = match SqlitePoolOptions::new()
            .connect(format!("sqlite://{location}?mode=rwc").as_str())
            .await
        {
            Ok(pool) => pool,
            Err(error) => return Err(DatabaseError::SqlxPoolCreateFailed(error)),
        };
        Ok(SqliteDatabase { pool })
    }

    async fn add_server(&mut self, server: &Server) -> Result<(), DatabaseError> {
        let serialized_server = serde_json::to_string(server).unwrap();
        let mut connection = match self.pool.acquire().await {
            Ok(connection) => connection,
            Err(error) => return Err(DatabaseError::SqlxAcquirePoolFailed(error)),
        };
        match sqlx::query("INSERT INTO servers ( serialized_server ) VALUES ( ?1 )")
            .bind(serialized_server)
            .execute(&mut *connection)
            .await
        {
            Ok(_) => Ok(()),
            Err(error) => Err(DatabaseError::SqlxInsertQueryFailed(error)),
        }
    }

    async fn add_player(&mut self, player: &Player) -> Result<(), DatabaseError> {
        let serialized_player = serde_json::to_string(player).unwrap();
        let mut connection = match self.pool.acquire().await {
            Ok(connection) => connection,
            Err(error) => return Err(DatabaseError::SqlxAcquirePoolFailed(error)),
        };
        match sqlx::query("INSERT INTO players ( serialized_player ) VALUES ( ?1 )")
            .bind(serialized_player)
            .execute(&mut *connection)
            .await
        {
            Ok(_) => Ok(()),
            Err(error) => Err(DatabaseError::SqlxInsertQueryFailed(error)),
        }
    }

    async fn add_targeted_server(&mut self, server: &Server) -> Result<(), DatabaseError> {
        let serialized_server = serde_json::to_string(server).unwrap();
        let mut connection = match self.pool.acquire().await {
            Ok(connection) => connection,
            Err(error) => return Err(DatabaseError::SqlxAcquirePoolFailed(error)),
        };
        match sqlx::query("INSERT INTO targeted_servers ( serialized_server ) VALUES ( ?1 )")
            .bind(serialized_server)
            .execute(&mut *connection)
            .await
        {
            Ok(_) => Ok(()),
            Err(error) => Err(DatabaseError::SqlxInsertQueryFailed(error)),
        }
    }

    async fn add_targeted_player(&mut self, player: &Player) -> Result<(), DatabaseError> {
        let serialized_player = serde_json::to_string(player).unwrap();
        let mut connection = match self.pool.acquire().await {
            Ok(connection) => connection,
            Err(error) => return Err(DatabaseError::SqlxAcquirePoolFailed(error)),
        };
        match sqlx::query("INSERT INTO targeted_players ( serialized_player ) VALUES ( ?1 )")
            .bind(serialized_player)
            .execute(&mut *connection)
            .await
        {
            Ok(_) => Ok(()),
            Err(error) => Err(DatabaseError::SqlxInsertQueryFailed(error)),
        }
    }
}
