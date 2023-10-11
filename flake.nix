{
    description = "openheimer flake";

    inputs = {
        nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
        mozilla.url = "github:mozilla/nixpkgs-mozilla";
        flake-utils.url = "github:numtide/flake-utils";
    };

    outputs = { self, nixpkgs, mozilla, flake-utils }:
        (flake-utils.lib.eachDefaultSystem (system:
            let
                overlays = [ self.inputs.mozilla.overlays.rust ];
                pkgs = import nixpkgs { inherit system overlays; };
                channel = pkgs.rustChannelOf {
                    date = "2023-10-11";
                    channel = "nightly";
                    sha256 = "sha256-gq7H6KCWVbf5rp6ceZVomtz/DOxM40i4TeWCIKxNAr8=";
                };
                rust = (channel.rust.override {
                    targets = [
                        "x86_64-unknown-linux-gnu"
                        "x86_64-unknown-linux-musl"
                    ];
                    extensions = [ "rust-src" ];
                });
            in rec
            {
                devShells.default = pkgs.mkShell {
                    name = "rust-environment";
                    nativeBuildInputs = [];
                    buildInputs = [ rust ];

                    PKG_CONFIG_ALLOW_CROSS = true;
                    PKG_CONFIG_ALL_STATIC = true;
                    LIBZ_SYS_STATIC = 1;
                };

                packages.openheimer = pkgs.rustPlatform.buildRustPackage {
                    pname = "openheimer";
                    version = "2.0.0-alpha";
                    cargoLock.lockFile = ./Cargo.lock;
                    src = pkgs.lib.cleanSource ./.;
                    nativeBuildInputs = [ pkgs.pkg-config ];
                    buildInputs = [ pkgs.udev ];
                };
                defaultPackage = packages.openheimer;
            }
        ));
}
