{
  description = "lazygit";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, }:
    let
      goVersion = 24; # Go 1.24
      supportedSystems =
        [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
    in flake-utils.lib.eachSystem supportedSystems (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ self.overlays.default ];
        };
      in {
        # Development shell
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # Go toolchain
            go
            gotools
            golangci-lint

            # Git for development and testing
            git

            # Build tools
            gnumake
          ];

          # Development environment
          shellHook = ''
            echo "lazygit development environment"
            echo "Go version: $(go version)"
            echo "Git version: $(git --version)"
            echo "Make version: $(make --version)"
            echo ""
          '';

          # Environment variables for development
          CGO_ENABLED = "0";
        };
      }) // {
        # Global overlay
        overlays.default = final: prev: {
          go = final."go_1_${toString goVersion}";
        };
      };
}
