{
  description = "openheimer";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
    rust-overlay.url = "github:oxalica/rust-overlay";
  };

  outputs = {
    flake-parts,
    nixpkgs,
    rust-overlay,
    self,
    ...
  } @ inputs:
    flake-parts.lib.mkFlake {inherit inputs;} {
      systems = [
        "aarch64-linux"
        "x86_64-linux"
      ];

      perSystem = {
        system,
        pkgs,
        ...
      }: let
        inherit (pkgs) pkgsCross pkgsStatic;
        rust = pkgs.rust-bin.nightly.latest.default.override {
          extensions = [
            "rust-src"
            "rust-analyzer-preview"
          ];
        };
      in {
        _module.args.pkgs = import nixpkgs {
          inherit system;
          overlays = [rust-overlay.overlays.default];
        };

        devShells.default = pkgs.mkShell {
          name = "openheimer";

          buildInputs = with pkgs; [
            clang
            curl
            libgit2
            mold
            pkgsCross.mingwW64.buildPackages.gcc
            rust
            taplo
            unzip
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

        packages = rec {
          openheimer = pkgs.rustPlatform.buildRustPackage {
            pname = "openheimer";
            version = self.shortRev or self.dirtyShortRev;

            src = pkgs.lib.cleanSource ./.;
            cargoLock.lockFile = ./Cargo.lock;

            nativeBuildInputs = with pkgs; [
              clang
              libgit2
              mold
              rust
            ];
          };
          default = openheimer;
        };
      };
    };

  nixConfig = {
    extra-substituters = ["https://errornobinaries.cachix.org/"];
    extra-trusted-public-keys = ["errornobinaries.cachix.org-1:84oagGNCIsXxBTYmfTiP+lvWje7lIS294iqAtCpFsbU="];
  };
}
