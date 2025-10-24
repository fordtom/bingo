{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {
          inherit system;
        };
        lib = pkgs.lib;

        # Build the Go bot with CGO (sqlite) enabled
        bingoPkg = pkgs.buildGoModule {
          pname = "bingo";
          version = "0.1.0";
          src = ./.;

          # First build: this will fail and print the correct hash.
          # Use a typed fake SRI so evaluation doesn't bail out early.
          vendorHash = "sha256-+5/WOmVNAM+JRsh6KtzdQ+sHlcShnP/csIkwiwJ1WmI=";
          go = pkgs.go;

          # CGO/sqlite
          env.CGO_ENABLED = "1";
          nativeBuildInputs = [pkgs.pkg-config];
          buildInputs = [pkgs.sqlite];

          # Build the main package at repo root to a 'bingo' binary
          subPackages = ["."];
        };
      in {
        packages = {
          default = bingoPkg;
          bingo = bingoPkg;
        };

        # nix run
        apps.default = {
          type = "app";
          program = "${bingoPkg}/bin/bingo";
        };

        # Dev shell with Go and sqlite headers/libs
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            gopls
            gcc
            sqlite
            pkg-config
          ];
          env = {
            CGO_ENABLED = "1";
          };
        };
      }
    );
}
