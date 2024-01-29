pub mod build_information {
    include!(concat!(env!("OUT_DIR"), "/built.rs"));
}

pub const VERSION: &str = env!("CARGO_PKG_VERSION");
pub const GIT_HASH: &str = if let Some(git_hash) = build_information::GIT_COMMIT_HASH_SHORT {
    git_hash
} else {
    "unknown commit"
};

pub fn format() -> String {
    format!("v{VERSION} ({GIT_HASH})")
}
