# MyGo
My Go experiments

## Play Pongo on an iPad

The `pongo2` game can run in Safari as a WebAssembly app. From the repository
root, build it and start a local web server:

```sh
cd pongo2
./build-web.sh
cd web/dist
python3 -m http.server 8000 --bind 0.0.0.0
```

Open the server URL on the iPad. On the same Wi-Fi network, use
`http://<computer-ip-address>:8000`. If this repository is open in GitHub
Codespaces, forward port 8000 and open its forwarded HTTPS URL in Safari.

The browser build embeds its font and supports touch input; the desktop build
continues to use mouse input.
