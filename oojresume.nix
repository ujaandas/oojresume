{
  lib,
  pkgs,
  resumegen,
}:
let
  inherit (lib)
    mkEnableOption
    mkIf
    mkOption
    optionalAttrs
    types
    ;

  mkResumeJson =
    cfg:
    builtins.toJSON [
      {
        identity = {
          name = cfg.identity.name;
          email = cfg.identity.email;
          phone = cfg.identity.phone;
          linkedin = cfg.identity.linkedin;
          github = cfg.identity.github;
          website = cfg.identity.website;
        };
        sections = map (
          s:
          {
            title = s.title;
            entries = s.entries;
          }
          // optionalAttrs (s.entryVSpace != null) {
            entryVSpace = s.entryVSpace;
          }
          // optionalAttrs (s.sectionVSpace != null) {
            sectionVSpace = s.sectionVSpace;
          }
        ) cfg.sections;
      }
    ];

  mkResumePackage =
    cfg:
    let
      resumeJson = pkgs.writeText "resume.json" (mkResumeJson cfg);
      texPkgs =
        with pkgs.texlive;
        combine {
          inherit
            scheme-basic
            latexmk
            collection-fontsrecommended
            geometry
            xcharter
            xstring
            xkeyval
            mweights
            fontaxes
            enumitem
            hyperref
            titlesec
            ;
        };
    in
    pkgs.stdenv.mkDerivation {
      pname = "resume-${cfg.name}";
      version = "0.1.0";
      src = ./.;

      nativeBuildInputs = [
        resumegen
        texPkgs
      ];

      buildPhase = ''
        runHook preBuild

        mkdir -p build
        cp -r ${./tmpl} build/tmpl
        cp ${./tex/preamble.tex} build/preamble.tex
        cp ${resumeJson} build/resume.json

        ${resumegen}/bin/resume \
          -dir build \
          -config resume.json \
          -tmpl build/tmpl \
          -out build/out

        cp build/preamble.tex build/out/preamble.tex

        (
          cd build/out
          latexmk -pdf -interaction=nonstopmode -halt-on-error main.tex
        )

        runHook postBuild
      '';

      installPhase = ''
        runHook preInstall

        mkdir -p $out
        cp build/out/main.tex $out/${cfg.name}.tex
        cp build/out/main.pdf $out/${cfg.name}.pdf

        runHook postInstall
      '';
    };

  module =
    { config, ... }:
    let
      cfg = config.oojresume;
    in
    {
      options.oojresume = {
        enable = mkEnableOption "resume build";

        name = mkOption {
          type = types.str;
          default = "default";
          description = "Output resume name.";
        };

        identity = {
          name = mkOption {
            type = types.str;
            default = "";
            description = "Full name.";
          };
          email = mkOption {
            type = types.str;
            default = "";
            description = "Email address.";
          };
          phone = mkOption {
            type = types.str;
            default = "";
            description = "Phone number.";
          };
          linkedin = mkOption {
            type = types.str;
            default = "";
            description = "LinkedIn path without protocol.";
          };
          github = mkOption {
            type = types.str;
            default = "";
            description = "GitHub path without protocol.";
          };
          website = mkOption {
            type = types.str;
            default = "";
            description = "Website URL/path.";
          };
        };

        sections = mkOption {
          type = types.listOf (
            types.submodule {
              options = {
                title = mkOption {
                  type = types.str;
                  description = "Section title.";
                };
                entries = mkOption {
                  type = types.listOf types.str;
                  default = [ ];
                  description = "Template names for this section.";
                };
                entryVSpace = mkOption {
                  type = types.nullOr types.int;
                  default = null;
                  description = "Optional spacing in pt inserted after each entry, e.g. -4.";
                };
                sectionVSpace = mkOption {
                  type = types.nullOr types.int;
                  default = null;
                  description = "Optional spacing in pt inserted after the section, e.g. -8.";
                };
              };
            }
          );
          default = [ ];
          description = "Ordered section list for the rendered resume.";
        };

        generatedJson = mkOption {
          type = types.nullOr types.str;
          default = null;
          description = "Rendered JSON config from options.";
        };

        package = mkOption {
          type = types.nullOr types.package;
          default = null;
          description = "Nix package that builds the resume PDF.";
        };
      };

      config = mkIf cfg.enable {
        oojresume.generatedJson = mkResumeJson cfg;
        oojresume.package = mkResumePackage cfg;
      };
    };
in
{
  inherit module mkResumePackage;
}
