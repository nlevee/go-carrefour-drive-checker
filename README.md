# Carrefour Drive Checker :alarm_clock:

Ce script permet la recherche de créneau disponible dans les Drives Carrefour.

## Installation

Récupérer le dernier binaire ici [https://github.com/nlevee/go-carrefour-drive-checker/releases/latest]()

Par exemple sous linux (x64) :

```
wget https://github.com/nlevee/go-carrefour-drive-checker/releases/download/v0.1.1/go-carrefour-drive-checker_v0.1.1_linux_amd64.tar.gz
tar xvzf go-carrefour-drive-checker_v0.1.1_linux_amd64.tar.gz
```

## Utilisation d'une liste de proxy

Pour utiliser une liste de proxy il faut ajouter l'option suivante aux options définies plus bas :

```bash
./go-carrefour-drive-checker -proxies http://proxy_1:port_1,http://proxy_2:port_2
```

A chaque appel vers l'API carrefour, un des serveur sera utilisé.

## Usage

Le script va tourner en continue et va afficher sur la console si un créneau est disponible

Pour lancer la recherche par code postal :

```bash
./go-carrefour-drive-checker -cp [CODE POSTAL]
```

Avec proxy : 

```bash
./go-carrefour-drive-checker -cp [CODE POSTAL] -proxies http://proxy_1:port_1,http://proxy_2:port_2
```

## Usage API

Pour rendre accessible la recherche de créneau via une mini API :

```bash
./go-carrefour-drive-checker -port 8089 -host 0.0.0.0 &
```

Avec proxy : 

```bash
./go-carrefour-drive-checker -port 8089 -host 0.0.0.0 -proxies http://proxy_1:port_1,http://proxy_2:port_2 &
```

Pour avoir la liste des clés de drive disponible :

```bash
curl 127.0.0.1:8089/stores?postalCode=[CODE POSTAL]
```

Pour ajouter un scrapper sur un store :

```bash
curl -XPUT 127.0.0.1:8089/scrappers/[ID DU DRIVE]
```

Pour checker l'état d'un drive :

```bash
curl 127.0.0.1:8089/scrappers/[ID DU DRIVE]
```
