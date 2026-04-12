{
  description = "Example development environment flake";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        resumegen = pkgs.buildGoModule {
          pname = "resume";
          version = "0.1.0";
          src = ./.;
          vendorHash = null;
        };

        resumeLib = import ./oojresume.nix {
          inherit (pkgs) lib;
          inherit pkgs resumegen;
        };

        mkResumeFromOptions =
          opts:
          (pkgs.lib.evalModules {
            modules = [
              resumeLib.module
              { oojresume = opts; }
            ];
          }).config.oojresume.package;

        resumeOptionFiles =
          let
            entries = builtins.readDir ./resumes;
            names = builtins.attrNames entries;
            nixFiles = builtins.filter (n: entries.${n} == "regular" && pkgs.lib.hasSuffix ".nix" n) names;
          in
          builtins.listToAttrs (
            map (n: {
              name = pkgs.lib.removeSuffix ".nix" n;
              value = ./resumes + "/${n}";
            }) nixFiles
          );

        resumeOptions = builtins.mapAttrs (_: path: import path) resumeOptionFiles;

        resumePackages = builtins.mapAttrs (
          variantName: opts: mkResumeFromOptions (opts // { name = opts.name or variantName; })
        ) resumeOptions;

        defaultResume =
          if resumePackages ? default then
            resumePackages.default
          else
            builtins.head (builtins.attrValues resumePackages);

        resumePackageNames = builtins.attrNames resumePackages;

        packageSelectors = pkgs.lib.concatStringsSep " " (map (name: ".#${name}") resumePackageNames);

        resumegenApp = "${
          (pkgs.writeShellApplication {
            name = "resumegen";
            text = ''
              set -euo pipefail

              target_dir="''${1:-./out}"
              mkdir -p "$target_dir"

              mapfile -t out_paths < <(nix build ${packageSelectors} --print-out-paths --no-link)
              for out_path in "''${out_paths[@]}"; do
                cp -f "$out_path"/*.pdf "$target_dir"/
              done

              echo "Copied PDFs to $target_dir"
            '';
          })
        }/bin/resumegen";

        tex = with pkgs.texlive; [
          (combine { inherit scheme-basic latexmk; })
        ];
      in
      {
        packages = {
          inherit resumegen;

          default = defaultResume;
        }
        // resumePackages;

        apps = {
          resumegen = {
            type = "app";
            program = resumegenApp;
          };

          default = {
            type = "app";
            program = resumegenApp;
          };
        };

        devShell = pkgs.mkShell {
          packages =
            with pkgs;
            [
              go
              go-tools
              gopls
              resumegen
            ]
            ++ tex;
        };
      }
    );
}
