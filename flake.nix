{
    description = "A page crawler written in go.";

    inputs = {
        flake-compat = {
            url = "github:edolstra/flake-compat";
            flake = false;
        };
    };

    outputs = {
        self,
        nixpkgs,
        ...
    }:
    let 
        inherit (nixpkgs) lib;
        supportedSystems = [ "x86_64-linux" "aarch64-linux" "aarch64-darwin" "x86_64-darwin" ];
        forAllSystems = lib.genAttrs supportedSystems;
        nixpkgsFor = forAllSystems(system: import nixpkgs {
            inherit system;
            configure.AllowUnfree = true;
        });
    in {

        packages = forAllSystems(system: 
        let
            pkgs = nixpkgsFor.${system};
        in {
            gocrawler = {
                pname = "gocrawler";
                version = "1.5.1";
                src = lib.cleanSource ./.;
                vendorSha256 = "sha256-cPOZ+95ajSi5AJL9aTegtI/7dre0nRB52v2pY6HD0P0=";
            };

            default = self.packages.${system}.gocrawler;
        });

        apps = forAllSystems(system: 
        let
            packages = self.packages.${system};
        in {
            gocrawler = {
                type = "app";
                program = "${packages.default}/bin/gocrawler";
            };

            default = self.apps.${system}.gocrawler;
        });

        devShells = forAllSystems(system: 
        let
            pkgs = nixpkgsFor.${system};
        in {
            default = pkgs.mkShell {
                nativeBuildInputs = with pkgs; [
                    gnumake
                    go
                ];
            };
        });
    };
}
