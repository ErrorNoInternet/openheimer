use super::Error;
use crate::minecraft::{player::Player, server::Server};
use async_trait::async_trait;
use sqlx::sqlite::SqlitePoolOptions;

pub struct Database {
    pool: sqlx::SqlitePool,
}

#[async_trait]
impl super::Database for Database {
    async fn new(location: String) -> Result<Self, Error> {
        let pool = match SqlitePoolOptions::new()
            .connect(format!("sqlite://{location}?mode=rwc").as_str())
            .await
        {
            Ok(pool) => pool,
            Err(error) => return Err(Error::PoolCreate(error)),
        };
        Ok(Database { pool })
    }

    async fn add_server(&mut self, server: &Server) -> Result<(), Error> {
        let serialized_server = serde_json::to_string(server).unwrap();
        let mut connection = match self.pool.acquire().await {
            Ok(connection) => connection,
            Err(error) => return Err(Error::AcquirePool(error)),
        };
        match sqlx::query("INSERT INTO servers ( serialized_server ) VALUES ( ?1 )")
            .bind(serialized_server)
            .execute(&mut *connection)
            .await
        {
            Ok(_) => Ok(()),
            Err(error) => Err(Error::InsertQuery(error)),
        }
    }

    async fn add_player(&mut self, player: &Player) -> Result<(), Error> {
        let serialized_player = serde_json::to_string(player).unwrap();
        let mut connection = match self.pool.acquire().await {
            Ok(connection) => connection,
            Err(error) => return Err(Error::AcquirePool(error)),
        };
        match sqlx::query("INSERT INTO players ( serialized_player ) VALUES ( ?1 )")
            .bind(serialized_player)
            .execute(&mut *connection)
            .await
        {
            Ok(_) => Ok(()),
            Err(error) => Err(Error::InsertQuery(error)),
        }
    }

    async fn add_targeted_server(&mut self, server: &Server) -> Result<(), Error> {
        let serialized_server = serde_json::to_string(server).unwrap();
        let mut connection = match self.pool.acquire().await {
            Ok(connection) => connection,
            Err(error) => return Err(Error::AcquirePool(error)),
        };
        match sqlx::query("INSERT INTO targeted_servers ( serialized_server ) VALUES ( ?1 )")
            .bind(serialized_server)
            .execute(&mut *connection)
            .await
        {
            Ok(_) => Ok(()),
            Err(error) => Err(Error::InsertQuery(error)),
        }
    }

    async fn add_targeted_player(&mut self, player: &Player) -> Result<(), Error> {
        let serialized_player = serde_json::to_string(player).unwrap();
        let mut connection = match self.pool.acquire().await {
            Ok(connection) => connection,
            Err(error) => return Err(Error::AcquirePool(error)),
        };
        match sqlx::query("INSERT INTO targeted_players ( serialized_player ) VALUES ( ?1 )")
            .bind(serialized_player)
            .execute(&mut *connection)
            .await
        {
            Ok(_) => Ok(()),
            Err(error) => Err(Error::InsertQuery(error)),
        }
    }
}
