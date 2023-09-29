{
  description = "Dev environment for LogPaste";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/release-23.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = import nixpkgs { inherit system; };

      pkgs_for_go = import (builtins.fetchTarball {
        # 1.18.4 release
        url = "https://github.com/NixOS/nixpkgs/archive/81094ccd6a0aa13ed176c815a60c4e25b49f072d.tar.gz";
        sha256 = "1h7bx9fnda66g651l8mymjzyb0y9km5b9sgwy2dla3dz4pvbk1zd";
      }) {inherit system; };
    in
    {
      devShells.default = pkgs_for_go.mkShell.override { stdenv = pkgs_for_go.pkgsStatic.stdenv; } {
        packages = with pkgs; [
          gopls
          gotools
          pkgs_for_go.go
        ];

        shellHook = ''
          go version
        '';
      };
    });
}
