{
    description = "A page crawler written in go.";

    inputs = {
        nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    };

    outputs = {
        self,
        nixpkgs
    }:
    let 
        system = "x86_64-linux";
        pkgs = import nixpkgs {
            inherit system;
            config.allowUnfree = true;
        };
        inherit (nixpkgs) lib;
    in {

        packages.${system} = {
            go-crawler = pkgs.buildGoModule rec {
                pname = "go-crawler";
                version = "1.5.0";

                src = lib.cleanSource ./.;

                vendorSha256 = "sha256-cPOZ+95ajSi5AJL9aTegtI/7dre0nRB52v2pY6HD0P0=";
            };

            default = self.packages.${system}.go-crawler;
        };

        apps.${system} =
        let
            packages = self.packages.${system};
        in {
            default = {
                type = "app";
                program = "${packages.default}/bin/gocrawler";
            };
        };

        devShells.${system}.default = pkgs.mkShell {
            nativeBuildInputs = with pkgs; [
                gnumake
                go
            ];
        };
    };
}
