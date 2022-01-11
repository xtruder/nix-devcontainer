{
  description = "dnix-devcontainer";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-21.05";
  };

  outputs = { self, nixpkgs  }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in
    {
      devShell.${system} = pkgs.mkShell {
        nativeBuildInputs = with pkgs; [
          # needed by nix
          nixpkgs-fmt
          rnix-lsp

          # for building docker images
          docker-client
          gnumake
          git

          # go dev
          go
        ];
      };
    };
}
