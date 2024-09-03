curl -LOs $(curl -s https://api.github.com/repos/bashidogames/gevm/releases/latest | grep browser_download_url | grep gevm-linux-amd64.zip | cut -d '"' -f 4)
unzip -qo gevm-linux-amd64.zip
rm gevm-linux-amd64.zip
mkdir -p ~/.local/bin
mv gevm ~/.local/bin
