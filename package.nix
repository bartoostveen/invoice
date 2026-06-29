{
  buildGoModule,
  lib,
  installShellFiles,
  ...
}:

buildGoModule (finalAttrs: {
  pname = "invoice";
  version = "0.0.1";

  src = ./.;

  vendorHash = "sha256-uj1yLyuJVoeB15lFgqLuhNTNHYe3lYLmOKFVQkZg7TI=";

  nativeBuildInputs = [
    installShellFiles
  ];

  postInstall = ''
    installShellCompletion --cmd invoice \
        --bash <($out/bin/invoice completion bash) \
        --fish <($out/bin/invoice completion fish) \
        --zsh <($out/bin/invoice completion zsh)

    $out/bin/invoice man > invoice.1
    installManPage invoice.1 
  '';

  meta = {
    description = "Generate invoices from the command line";
    homepage = "https://github.com/maaslalani/invoice";
    license = lib.licenses.mit;
  };
})
