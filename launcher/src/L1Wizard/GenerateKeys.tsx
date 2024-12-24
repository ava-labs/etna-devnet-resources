const dockerInstallation = `sudo yum update -y
sudo yum -y install docker

sudo service docker start

sudo usermod -a -G docker ec2-user
sudo chmod 666 /var/run/docker.sock
docker version
`

const dockerCommand = `mkdir -p ~/.avalanchego/staking; docker run -it -d \\
  --name avalanchego \\
  --network host \\
  -v ~/.avalanchego:/home/avalanche/.avalanchego \\
  -e AVALANCHEGO_NETWORK_ID=fuji \\
  -e AVALANCHEGO_PARTIAL_SYNC_PRIMARY_NETWORK=true \\
  -e HOME=/home/avalanche \\
  --user $(id -u):$(id -g) \\
  containerman17/easy-avalanchego:${CONTAINER_VERSION}`

const popRequest = `curl -X POST --data '{ 
    "jsonrpc":"2.0", 
    "id"     :1, 
    "method" :"info.getNodeID" 
}' -H 'content-type:application/json;' 127.0.0.1:9650/ext/info`

const stopScript = `docker stop avalanchego; docker rm avalanchego`

import { CONTAINER_VERSION } from './constants';
import { useWizardStore } from './store';
import NextPrev from './ui/NextPrev';

export default function GenerateKeys() {
    const { nodesCount, setNodesCount } = useWizardStore();
    const nodeCounts = [1, 3, 5];

    return <>
        <h1 className="text-2xl font-medium mb-6">
            Set up your Validator Nodes
        </h1>

        <p className="mb-4">
            This step will guide you through generating keys for your validator nodes. We recommend performing this step on the infrastructure where you plan to run your nodes.
        </p>

        <p className="mb-4">
            <strong>Requirements for validator nodes:</strong>
            <ul className="list-disc list-inside ml-4">
                <li>16GB RAM (you might try with 8GB)</li>
                <li>8 cores CPU (you might try 4 cores)</li>
                <li>
                    At least 100GB of any disk space (EBS or SSD), except for HDD
                </li>
                <li>
                    <strong>⚠️ Important:</strong> make sure port 9651 is open on your node!
                </li>
            </ul>
            If you are hosting the validators on AWS you can use t2.2xlarge EC2 instances.
        </p>

        <h3 className="mb-4 font-medium">Docker</h3>
        <p className="mb-4">
            We will retrieve the binary images of <a href='https://github.com/ava-labs/avalanchego' target='_blank'>AvalancheGo</a> from the Docker Hub. Make sure you have Docker installed on your system. To install Docker on an AWS machine, run the following commands:
        </p>

        <pre className="bg-gray-100 p-4 rounded-md mb-4">{dockerInstallation}</pre>

        <p className="mb-4">
            If you do not want to use Docker, you can follow the instructions here <a href="https://github.com/ava-labs/avalanchego?tab=readme-ov-file#installation" target="_blank">here</a>.
        </p>

        <h3 className="mb-4 font-medium">How many nodes do you want to run?</h3>
        <ul className="mb-4 items-center w-full text-sm font-medium text-gray-900 bg-white border border-gray-200 rounded-lg sm:flex">
            {nodeCounts.map((count) => (
                <li key={count} className="w-full border-b border-gray-200 sm:border-b-0 sm:border-r last:border-r-0 ">
                    <div className="flex items-center ps-3">
                        <input
                            id={`nodes-${count}`}
                            type="radio"
                            checked={nodesCount === count}
                            onChange={() => setNodesCount(count)}
                            name="nodes-count"
                            className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 focus:ring-0"
                        />
                        <label htmlFor={`nodes-${count}`} className="w-full py-3 ms-2 text-sm font-medium text-gray-900">
                            {count} {count === 1 ? 'Node' : 'Nodes'}
                            {count === 1 && (
                                <span className="ms-2 bg-red-100 text-red-800 text-xs font-medium px-2.5 py-0.5 rounded-full">
                                    Dev
                                </span>
                            )}
                            {count === 3 && (
                                <span className="ms-2 bg-green-100 text-green-800 text-xs font-medium px-2.5 py-0.5 rounded-full">
                                    Testnet
                                </span>
                            )}
                            {count === 5 && (
                                <span className="ms-2 bg-blue-100 text-blue-800 text-xs font-medium px-2.5 py-0.5 rounded-full">
                                    Mainnet
                                </span>
                            )}
                        </label>
                    </div>
                </li>
            ))}
        </ul>

        <h1 className="text-2xl font-medium mb-6">Generate Keys</h1>

        <p className="mb-4">
            For creating the L1, we need to know the node IDs, the BLS public keys and the proof of possession (POP) of the nodes. To generate the keys for the nodes, we will briefly start the nodes and request the keys from them. Afterwards, we will immediately stop the nodes.
        </p>
        
        <h3 className="mb-4 font-medium">Run this on {nodesCount === 1 ? "the" : "every"} node:</h3>
        <pre className="bg-gray-100 p-4 rounded-md mb-4">{dockerCommand}</pre>
        <p className="mb-4">
            This runs an avalanchego node in docker. The node, while starting, generates its own keys if they are not present.
            You can find them at <code>~/.avalanchego/staking/</code> in your local system.
        </p>

        <h3 className="mb-4 font-medium">Then request node credentials:</h3>
        <pre className="bg-gray-100 p-4 rounded-md mb-4">{popRequest}</pre>

        <p className="mb-4">Save the responses. The response will contain fields <code>nodeID</code> and <code>nodePOP</code> (proof of possession). We will need them to convert the subnet to L1.</p>

        <h3 className="mb-4 font-medium">Stop and remove the nodes:</h3>
        <pre className="bg-gray-100 p-4 rounded-md mb-4">{stopScript}</pre>

        <p className="mb-4">
            Please don't forget to stop the nodes or subsequent steps will fail.
        </p>

        <NextPrev nextDisabled={false} currentStepName="generate-keys" />
    </>
}
