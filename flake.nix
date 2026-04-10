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
        tex = with pkgs.texlive; [
          (combine { inherit scheme-basic latexmk; })
        ];
      in
      {
        devShell = pkgs.mkShell {
          packages =
            with pkgs;
            [
              go
              go-tools
              gopls
            ]
            ++ tex;
        };
      }
    );
}
