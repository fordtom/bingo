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
      in {
        packages.default = pkgs.buildGoModule {
          pname = "bingo-bot";
          version = "0.1.0";
          src = ./.;
          vendorHash = null; # Use go.sum instead of vendoring
          CGO_ENABLED = 1;
          nativeBuildInputs = with pkgs; [gcc];
          buildInputs = with pkgs; [sqlite];
        };

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            gopls
            gcc
            sqlite
          ];
          env = {
            CGO_ENABLED = 1;
          };
        };
      }
    );
}
