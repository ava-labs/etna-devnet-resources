const dockerCommand = `mkdir -p ~/.avalanchego/staking; docker run -it -d \\
  --name avalanchego \\
  --network host \\
  -v ~/.avalanchego:/home/avalanche/.avalanchego \\
  -e AVALANCHEGO_NETWORK_ID=fuji \\
  -e AVALANCHEGO_PARTIAL_SYNC_PRIMARY_NETWORK=true \\
  -e HOME=/home/avalanche \\
  --user $(id -u):$(id -g) \\
  containerman17/avalanchego-subnetevm:v1.12.1_v0.7.0`

const popRequest = `curl -X POST --data '{ 
    "jsonrpc":"2.0", 
    "id"     :1, 
    "method" :"info.getNodeID" 
}' -H 'content-type:application/json;' 127.0.0.1:9650/ext/info`

const stopScript = `docker stop avalanchego; docker rm avalanchego`

import { useWizardStore } from './store';

export default function GenerateKeys() {
    const { nodesCount, setNodesCount, advanceFrom } = useWizardStore();
    const nodeCounts = [1, 3, 5];

    return <>
        <h1 className="text-2xl font-medium mb-6">Generate Keys</h1>

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
                            {count === 3 && (
                                <span className="ms-2 bg-blue-100 text-blue-800 text-xs font-medium px-2.5 py-0.5 rounded-full">
                                    Recommended
                                </span>
                            )}
                        </label>
                    </div>
                </li>
            ))}
        </ul>
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

        <div className="flex justify-between">
            <button
                onClick={() => advanceFrom('generate-keys', 'down')}
                className="px-4 py-2 rounded-md bg-gray-100 hover:bg-gray-200"
            >
                Previous
            </button>
            <button
                onClick={() => advanceFrom('generate-keys')}
                className="px-4 py-2 rounded-md bg-blue-500 text-white hover:bg-blue-600"
            >
                Continue
            </button>
        </div>
    </>
}
