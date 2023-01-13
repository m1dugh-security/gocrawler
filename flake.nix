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
                version = "1.4.7";

                src = ./.;

                vendorHash = "sha256-pQpattmS9VmO3ZIQUFn66az8GSmB4IvYhTTCFn6SUmo=";
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

    };
}
