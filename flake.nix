{
  description = "Dev environment for LogPaste";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";

    # 1.21.1 release
    go-nixpkgs.url = "github:NixOS/nixpkgs/78058d810644f5ed276804ce7ea9e82d92bee293";

    # 3.44.2 release
    sqlite-nixpkgs.url = "github:NixOS/nixpkgs/5ad9903c16126a7d949101687af0aa589b1d7d3d";

    # 20.6.1 release
    nodejs-nixpkgs.url = "github:NixOS/nixpkgs/78058d810644f5ed276804ce7ea9e82d92bee293";

    # 0.9.0 release
    shellcheck-nixpkgs.url = "github:NixOS/nixpkgs/8b5ab8341e33322e5b66fb46ce23d724050f6606";

    # 1.2.1 release
    sqlfluff-nixpkgs.url = "github:NixOS/nixpkgs/7cf5ccf1cdb2ba5f08f0ac29fc3d04b0b59a07e4";

    # 0.1.147 release
    flyctl-nixpkgs.url = "github:NixOS/nixpkgs/0a254180b4cad6be45aa46dce896bdb8db5d2930";

    # 0.3.13 release
    litestream-nixpkgs.url = "github:NixOS/nixpkgs/a343533bccc62400e8a9560423486a3b6c11a23b";
  };

  outputs = {
    self,
    flake-utils,
    go-nixpkgs,
    sqlite-nixpkgs,
    nodejs-nixpkgs,
    shellcheck-nixpkgs,
    sqlfluff-nixpkgs,
    flyctl-nixpkgs,
    litestream-nixpkgs,
  } @ inputs:
    flake-utils.lib.eachDefaultSystem (system: let
      gopkg = go-nixpkgs.legacyPackages.${system};
      go = gopkg.go_1_21;
      sqlite = sqlite-nixpkgs.legacyPackages.${system}.sqlite;
      nodejs = nodejs-nixpkgs.legacyPackages.${system}.nodejs_20;
      shellcheck = shellcheck-nixpkgs.legacyPackages.${system}.shellcheck;
      sqlfluff = sqlfluff-nixpkgs.legacyPackages.${system}.sqlfluff;
      flyctl = flyctl-nixpkgs.legacyPackages.${system}.flyctl;
      litestream = litestream-nixpkgs.legacyPackages.${system}.litestream;
    in {
      devShells.default = gopkg.mkShell.override { stdenv = gopkg.pkgsStatic.stdenv; } {
        packages = [
          gopkg.gotools
          gopkg.gopls
          gopkg.go-outline
          gopkg.gocode
          gopkg.gopkgs
          gopkg.gocode-gomod
          gopkg.godef
          gopkg.golint
          go
          sqlite
          nodejs
          shellcheck
          sqlfluff
          flyctl
          litestream
        ];

        shellHook = ''
          export GOROOT="${go}/share/go"

          echo "shellcheck" "$(shellcheck --version | grep '^version:')"
          sqlfluff --version
          echo "litestream" "$(litestream version)"
          fly version | cut -d ' ' -f 1-3
          echo "node" "$(node --version)"
          echo "npm" "$(npm --version)"
          echo "sqlite" "$(sqlite3 --version | cut -d ' ' -f 1-2)"
          go version
        '';
      };

      formatter = gopkg.alejandra;
    });
}
