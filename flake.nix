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

        tex = with pkgs.texlive; [
          (combine { inherit scheme-basic latexmk; })
        ];
      in
      {
        packages = {
          inherit resumegen;

          default = mkResumeFromOptions {
            enable = true;
            name = "default";
            identity = {
              name = "Ujaan Das";
              email = "ujaandas03@gmail.com";
              linkedin = "linkedin.com/in/ujaandas";
              github = "github.com/ujaandas";
            };
            sections = [
              {
                title = "Education";
                entries = [
                  "edu_warwick"
                  "edu_hkust"
                ];
              }
              {
                title = "Experience";
                entries = [
                  "work_stellerus_swe_2025"
                  "work_hkust_castle_2024"
                  "work_stellerus_sde_2023"
                ];
              }
              {
                title = "Projects";
                entries = [
                  "proj_dissertation"
                  "proj_yywm"
                  "proj_snip"
                  "proj_follow_me_robot"
                ];
              }
              {
                title = "Skills";
                entries = [ "skills_default" ];
              }
            ];
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
