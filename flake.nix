{
  description = "Dev environment for LogPaste";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";

    # 1.21.1 release
    go_dep.url = "github:NixOS/nixpkgs/78058d810644f5ed276804ce7ea9e82d92bee293";

    # 20.6.1 release
    nodejs_dep.url = "github:NixOS/nixpkgs/78058d810644f5ed276804ce7ea9e82d92bee293";

    # 0.9.0 release
    shellcheck_dep.url = "github:NixOS/nixpkgs/8b5ab8341e33322e5b66fb46ce23d724050f6606";

    # 1.2.1 release
    sqlfluff_dep.url = "github:NixOS/nixpkgs/7cf5ccf1cdb2ba5f08f0ac29fc3d04b0b59a07e4";

    # 0.1.147 release
    flyctl_dep.url = "github:NixOS/nixpkgs/0a254180b4cad6be45aa46dce896bdb8db5d2930";
  };

  outputs = { self, flake-utils, go_dep, nodejs_dep, shellcheck_dep, sqlfluff_dep, flyctl_dep }@inputs :
    flake-utils.lib.eachDefaultSystem (system:
    let
      go_dep = inputs.go_dep.legacyPackages.${system};
      nodejs_dep = inputs.nodejs_dep.legacyPackages.${system};
      shellcheck_dep = inputs.shellcheck_dep.legacyPackages.${system};
      sqlfluff_dep = inputs.sqlfluff_dep.legacyPackages.${system};
      flyctl_dep = inputs.flyctl_dep.legacyPackages.${system};
    in
    {
      devShells.default = go_dep.mkShell.override { stdenv = go_dep.pkgsStatic.stdenv; } {
        packages = [
          go_dep.gopls
          go_dep.gotools
          go_dep.go_1_21
          nodejs_dep.nodejs_20
          shellcheck_dep.shellcheck
          sqlfluff_dep.sqlfluff
          flyctl_dep.flyctl
        ];

        shellHook = ''
          echo "shellcheck" "$(shellcheck --version | grep '^version:')"
          sqlfluff --version
          fly version | cut -d ' ' -f 1-3
          echo "node" "$(node --version)"
          echo "npm" "$(npm --version)"
          go version
        '';
      };
    });
}
