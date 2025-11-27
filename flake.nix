{
  description = "A simple terminal UI for git commands";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
    flake-parts.url = "github:hercules-ci/flake-parts";
    flake-compat.url = "https://flakehub.com/f/edolstra/flake-compat/1.tar.gz";
    treefmt-nix.url = "github:numtide/treefmt-nix";
  };

  outputs =
    inputs@{ flake-parts, systems, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import systems;
      imports = [
        inputs.treefmt-nix.flakeModule
      ];

      perSystem =
        {
          pkgs,
          system,
          ...
        }:
        let
          goMod = builtins.readFile ./go.mod;
          versionMatch = builtins.match ".*go[[:space:]]([0-9]+\\.[0-9]+)(\\.[0-9]+)?.*" goMod;

          goVersion =
            if versionMatch != null then
              builtins.head versionMatch
            else
              throw "Could not extract Go version from go.mod";

          goOverlay = final: prev: {
            go = prev."go_${builtins.replaceStrings [ "." ] [ "_" ] goVersion}";
          };

          lazygit = pkgs.buildGoModule rec {
            pname = "lazygit";
            version = "dev";

            gitCommit = inputs.self.rev or inputs.self.dirtyRev or "dev";

            src = ./.;
            vendorHash = null;

            # Disable integration tests that require specific environment
            doCheck = false;

            nativeBuildInputs = with pkgs; [
              git
              makeWrapper
            ];
            buildInputs = [ pkgs.git ];

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

            meta = {
              description = "A simple terminal UI for git commands";
              homepage = "https://github.com/jesseduffield/lazygit";
              license = pkgs.lib.licenses.mit;
              maintainers = [ "jesseduffield" ];
              platforms = pkgs.lib.platforms.unix;
              mainProgram = "lazygit";
            };
          };
        in
        {
          _module.args.pkgs = import inputs.nixpkgs {
            inherit system;
            overlays = [ goOverlay ];
            config = { };
          };

          packages = {
            default = lazygit;
            inherit lazygit;
          };

          devShells.default = pkgs.mkShell {
            name = "lazygit-dev";

            buildInputs = with pkgs; [
              # Go toolchain
              go
              gotools

              # Development tools
              git
              gnumake
            ];

            # Environment variables for development
            CGO_ENABLED = "0";
          };

          treefmt = {
            programs.nixfmt.enable = pkgs.lib.meta.availableOn pkgs.stdenv.buildPlatform pkgs.nixfmt-rfc-style.compiler;
            programs.nixfmt.package = pkgs.nixfmt-rfc-style;
            programs.gofmt.enable = true;
          };

          checks.build = lazygit;
        };

      flake = {
        overlays.default = final: prev: {
          lazygit = inputs.self.packages.${final.system}.lazygit;
        };
      };
    };
}
