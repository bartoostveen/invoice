{
  description = "Generic devshell flake";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs?ref=nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";

    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "aarch64-darwin"
        "x86_64-darwin"
      ];

      imports = [
        inputs.treefmt-nix.flakeModule
      ];

      perSystem =
        {
          self',
          pkgs,
          ...
        }:
        {
          treefmt = {
            programs.nixfmt.enable = true;
            programs.gofmt.enable = true;
            programs.goimports.enable = true;
            programs.shellcheck.enable = true;
            programs.yamlfmt.enable = true;
          };

          devShells.default = pkgs.mkShell {
            packages = with pkgs; [
              go
              gopls
              gotools
              vhs
            ];

            env = {
            };

            shellHook = ''
              export GOPATH=$PWD/.gopath
              export PATH=$GOPATH/bin:$PATH
              mkdir -p $GOPATH
              go telemetry off # This doesn't restore original state
            '';
          };

          packages.default = pkgs.callPackage ./package.nix { };
          packages.invoice = self'.packages.default;

          packages.tape = pkgs.writeShellApplication {
            name = "tape";

            runtimeInputs = with pkgs; [
              vhs
              bash
              self'.packages.default
            ];

            text = ''
              vhs invoice.tape
            '';
          };

          packages.extract-font =
            pkgs.writers.writePython3Bin "extract-font"
              {
                libraries = with pkgs.python3Packages; [
                  fontforge
                ];
              }
              ''
                import fontforge
                import shutil
                import os
                import sys

                inter = "${pkgs.inter}"
                file = inter + "/share/fonts/truetype/Inter.ttc"

                destination = "invoice/inter"
                if (len(sys.argv) > 1):
                    destination = sys.argv[1]

                try:
                    shutil.rmtree(destination)
                except FileNotFoundError:
                    pass

                os.makedirs(destination)

                for font in fontforge.fontsInFile(file):
                    f = fontforge.open(u'%s(%s)' % (file, font))
                    f.generate('%s/%s.ttf' % (destination, font))

                shutil.copy(inter + "/share/fonts/truetype/InterVariable.ttf", destination)
              '';
        };

      flake = {
        inherit inputs;
      };
    };
}
