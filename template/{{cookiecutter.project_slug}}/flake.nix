{
    description = "{{cookiecutter.project_slug}}";

    inputs = {
      nixpkgs.url = "github:nixos/nixpkgs/{{cookiecutter.nixpkgs_branch}}";
    };

    outputs = { self, nixpkgs }: let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      devShell.${system} = pkgs.mkShell {
        buildInputs = with pkgs; [
        ];

        shellHook = ''
        '';
      };
    };
}