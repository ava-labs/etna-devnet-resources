const dockerCommand = `mkdir -p ~/.avalanchego/staking; docker run -it -d \\
  --name avalanchego \\
  --network host \\
  -v ~/.avalanchego:/home/avalanche/.avalanchego \\
  -e AVALANCHEGO_NETWORK_ID=fuji \\
  -e AVALANCHEGO_HTTP_ALLOWED_HOSTS=* \\
  -e AVALANCHEGO_HTTP_HOST=0.0.0.0 \\
  -e AVALANCHEGO_PARTIAL_SYNC_PRIMARY_NETWORK=true \\
  -e HOME=/home/avalanche \\
  --user $(id -u):$(id -g) \\
  containerman17/avalanchego-subnetevm:v1.12.1_v0.7.0`

const popRequest = `curl -X POST --data '{ 
    "jsonrpc":"2.0", 
    "id"     :1, 
    "method" :"info.getNodeID" 
}' -H 'content-type:application/json;' 127.0.0.1:9650/ext/info`

import { useWizardStore } from './store';

export default function GenerateKeys() {
    const { nodesCount, setNodesCount } = useWizardStore();
    const nodeCounts = [1, 3, 5];

    return <>
        <h1 className="text-2xl font-medium mb-6">Generate Keys</h1>

        <h3 className="mb-4 font-medium">How many nodes do you want to run?</h3>
        <ul className="mb-8 items-center w-full text-sm font-medium text-gray-900 bg-white border border-gray-200 rounded-lg sm:flex dark:bg-gray-700 dark:border-gray-600 dark:text-white">
            {nodeCounts.map((count) => (
                <li key={count} className="w-full border-b border-gray-200 sm:border-b-0 sm:border-r last:border-r-0 dark:border-gray-600">
                    <div className="flex items-center ps-3">
                        <input
                            id={`nodes-${count}`}
                            type="radio"
                            checked={nodesCount === count}
                            onChange={() => setNodesCount(count)}
                            name="nodes-count"
                            className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 focus:ring-0 dark:focus:ring-0 dark:bg-gray-600 dark:border-gray-500"
                        />
                        <label htmlFor={`nodes-${count}`} className="w-full py-3 ms-2 text-sm font-medium text-gray-900 dark:text-gray-300">
                            {count} {count === 1 ? 'Node' : 'Nodes'}
                            {count === 3 && (
                                <span className="ms-2 bg-blue-100 text-blue-800 text-xs font-medium px-2.5 py-0.5 rounded-full dark:bg-blue-900 dark:text-blue-300">
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
        <h3 className="mb-4 font-medium">Then request node credentials:</h3>
        <pre className="bg-gray-100 p-4 rounded-md mb-4">{popRequest}</pre>
    </>
}