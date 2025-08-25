{
  description = "A simple terminal UI for git commands";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    let
      supportedSystems =
        [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];

      gitCommit = self.rev or self.dirtyRev or "dev";
      version = "2.2.1";

    in flake-utils.lib.eachSystem supportedSystems (system:
      let
        pkgs = import nixpkgs { inherit system; };

        lazygit = pkgs.buildGoModule rec {
          pname = "lazygit";
          inherit version;

          src = ./.;

          vendorHash = null;

          # Disable integration tests that require specific environment
          doCheck = false;

          nativeBuildInputs = with pkgs; [ git makeWrapper ];

          buildInputs = with pkgs; [ git ];

          ldflags = [
            "-s"
            "-w"
            "-X main.commit=${gitCommit}"
            "-X main.version=${version}"
            "-X main.buildSource=nix"
          ];

          postInstall = ''
            wrapProgram $out/bin/lazygit \
              --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.git ]}
          '';

          meta = with pkgs.lib; {
            description = "A simple terminal UI for git commands";
            homepage = "https://github.com/jesseduffield/lazygit";
            license = licenses.mit;
            maintainers = [ "jesseduffield" ];
            platforms = platforms.unix;
            mainProgram = "lazygit";
          };
        };

      in {
        packages = {
          default = lazygit;
          inherit lazygit;
        };

        apps = {
          default = flake-utils.lib.mkApp {
            drv = lazygit;
            name = "lazygit";
          };
          lazygit = flake-utils.lib.mkApp {
            drv = lazygit;
            name = "lazygit";
          };
        };

        devShells.default = pkgs.mkShell {
          name = "lazygit-dev";

          buildInputs = with pkgs; [
            # Go toolchain
            go_1_24
            gotools
            golangci-lint

            # Development tools
            git
            gnumake
          ];

          shellHook = ''
            echo "Lazygit development environment"
            echo "Go version: $(go version)"
            echo "Git version: $(git --version)"
            echo ""
          '';

          # Environment variables for development
          CGO_ENABLED = "0";
        };

        # Formatting check
        formatter = pkgs.nixpkgs-fmt;

        # Development checks
        checks = {
          # Ensure the package builds
          build = lazygit;

          # Format check
          format = pkgs.runCommand "check-format" {
            buildInputs = [ pkgs.nixpkgs-fmt ];
          } ''
            nixpkgs-fmt --check ${./.}
            touch $out
          '';
        };
      }) // {
        # Global overlay for other flakes to use
        overlays.default = final: prev: {
          lazygit = self.packages.${final.system}.lazygit;
        };

        # CI/CD support
        hydraJobs = self.packages;
      };
}
