{
  description = "openheimer";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
    rust-overlay.url = "github:oxalica/rust-overlay";
  };

  outputs = {
    nixpkgs,
    flake-parts,
    rust-overlay,
    ...
  } @ inputs:
    flake-parts.lib.mkFlake {inherit inputs;} {
      systems = ["x86_64-linux" "aarch64-linux" "aarch64-darwin"];
      perSystem = {
        system,
        pkgs,
        ...
      }: let
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
        _module.args.pkgs = import inputs.nixpkgs {
          inherit system;
          overlays = [rust-overlay.overlays.default];
        };

        devShells.default = pkgs.mkShell {
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
        packages.default = packages.openheimer;
      };
    };
}
