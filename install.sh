curl -LOs $(curl -s https://api.github.com/repos/bashidogames/gdvm/releases/latest | grep browser_download_url | grep gdvm-linux-amd64.zip | cut -d '"' -f 4)
unzip -qo gdvm-linux-amd64.zip
rm gdvm-linux-amd64.zip
mkdir -p ~/.local/bin
mv gdvm ~/.local/bin
