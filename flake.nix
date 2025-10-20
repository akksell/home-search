{
  description = "Development Environment";
  
  nixConfig = {
    bash-prompt = "\\[\\e[96m\\]\\A\\[\\e[0m\\] \\[\\e[35m\\][NIX]\\[\\e[0m\\] \\[\\e[1m\\]$(__git_ps1 \"(%s)\")\\[\\e[0m\\]\n\\w \$ ";
  };

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/98ff3f9af2684f6136c24beef08f5e2033fc5389"; # nixos 25.05
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachSystem flake-utils.lib.defaultSystems (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          config = {
            allowUnfree = true;
          };
          overlays = [];
        };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # Go development
            go
            gotools
            gopls
            go-outline
            gopkgs
            gocode-gomod
            godef
            golint
            
            # Python development
            python314
            
            # Protocol Buffers
            buf
            
            # Build tooling
            bazel_7
            
            # Google Cloud SDK
            google-cloud-sdk
            
            # Terraform
            terraform
            
            # Additional development tools
            git
            curl
            jq
            gnumake
            
          ];

          shellHook = ''
            # source git-related files for prompt and completion
            source ${pkgs.git}/share/git/contrib/completion/git-prompt.sh
            source ${pkgs.git}/share/git/contrib/completion/git-completion.bash
            export GIT_PS1_SHOWDIRTYSTATE=1

            echo "🏠 HomeSearch Shell Environment"
            echo "=================================="
            
            # Set up Go environment
            export GOPATH=$HOME/go
            export PATH=$PATH:$GOPATH/bin
          '';

          # Environment variables
          env = {
            # Go configuration
            CGO_ENABLED = "1";
            
            # Python configuration
            PYTHONDONTWRITEBYTECODE = "1";
          };
        };
      }

    );

}