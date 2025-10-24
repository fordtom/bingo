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
        # Build the Go bot with CGO (sqlite) enabled
        bingoPkg = pkgs.buildGoModule {
          pname = "bingo";
          version = "0.1.0";
          src = ./.;

          vendorHash = "sha256-+5/WOmVNAM+JRsh6KtzdQ+sHlcShnP/csIkwiwJ1WmI=";
          go = pkgs.go;

          env.CGO_ENABLED = "1";
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
          ];
          env = {
            CGO_ENABLED = "1";
          };
        };
      }
    );
}
