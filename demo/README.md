# DEMO

Here we provide a few scripts to demonstrate how to interact with evm
nodes.

You might need to install NodeJS and dependencies first:

```bash
# install node version manager
$ curl -o- https://raw.githubusercontent.com/creationix/nvm/v0.33.5/install.sh | bash
# use nvm to install stable version of node
$ nvm install node stable
# install dependencies for this demo
evm/demo$ npm install
```

To start a testnet, execute the deploy commands from the `deploy/` directory:

ex:

```bash
evm$ cd deploy
evm/deploy$ make CONSENSUS=dag1 NODES=4
```

Then, in an other terminal, start the interactive demo:

```bash
$ ./demo.sh ../deploy/terraform/local/ips.dat
```

The ips.dat file, generated during the deploy phase, tells the demo program
where to reach the nodes.

In this case, we are using DAG1 consensus, so it is interesting to monitor
the dag1 nodes:

```bash
$ ./watch.sh ../deploy/terraform/local/ips.dat
```

After the demo, destroy the testnet by running `make stop` from `deploy/`:

```bash
$ cd ../deploy
$ make stop
```
