{
  description = "openheimer";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    rust-overlay.url = "github:oxalica/rust-overlay";
  };

  outputs = {
    nixpkgs,
    flake-utils,
    rust-overlay,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs {
        inherit system;
        overlays = [rust-overlay.overlays.default];
      };
      inherit (pkgs) pkgsStatic;
      inherit (pkgs) pkgsCross;

      rust = pkgs.rust-bin.nightly.latest.default.override {
        targets = [
          "x86_64-unknown-linux-gnu"
          "x86_64-unknown-linux-musl"
          "i686-unknown-linux-musl"
          "x86_64-pc-windows-gnu"
        ];
        extensions = [
          "rust-src"
          "rust-analyzer-preview"
        ];
      };
    in rec {
      devShell = pkgs.mkShell {
        name = "openheimer";
        buildInputs = with pkgs; [
          clang
          libgit2
          mold
          pkgsCross.mingwW64.buildPackages.gcc
          rust
        ];

        CARGO_TARGET_X86_64_PC_WINDOWS_GNU_RUSTFLAGS = "-L native=${pkgsCross.mingwW64.windows.mingw_w64_pthreads}/lib";
        LIBZ_SYS_STATIC = 1;
        OPENSSL_DIR = pkgsStatic.openssl.dev;
        OPENSSL_LIB_DIR = "${pkgsStatic.openssl.out}/lib";
        OPENSSL_STATIC = 1;
        PKG_CONFIG_ALL_STATIC = true;
        PKG_CONFIG_ALLOW_CROSS = true;
        RUST_BACKTRACE = 1;
      };

      packages.openheimer = pkgs.rustPlatform.buildRustPackage {
        pname = "openheimer";
        version = "dev";

        cargoLock.lockFile = ./Cargo.lock;
        src = pkgs.lib.cleanSource ./.;

        nativeBuildInputs = with pkgs; [
          clang
          libgit2
          mold
          rust
        ];
      };
      defaultPackage = packages.openheimer;
    });
}
