{
    description = "debian-nix-devcontainer";

    inputs = {
      nixpkgs.url = "github:nixos/nixpkgs/nixos-20.09";
    };

    outputs = { self, nixpkgs }: let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      devShell.${system} = pkgs.mkShell {
        buildInputs = with pkgs; [
          docker
          cookiecutter
        ];
      };
    };
}