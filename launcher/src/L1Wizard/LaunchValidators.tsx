import { useWizardStore } from './store';
import NextPrev from './ui/NextPrev';
import { useState } from 'react';
import Note from './ui/Note';

const dockerCommand = (subnetID: string) => `mkdir -p ~/.avalanchego/staking; docker run -it -d \\
  --name avalanchego \\
  --network host \\
  -v ~/.avalanchego:/home/avalanche/.avalanchego \\
  -e AVALANCHEGO_NETWORK_ID=fuji \\
  -e AVALANCHEGO_PARTIAL_SYNC_PRIMARY_NETWORK=true \\
  -e AVALANCHEGO_TRACK_SUBNETS=${subnetID} \\
  -e HOME=/home/avalanche \\
  --user $(id -u):$(id -g) \\
  containerman17/avalanchego-subnetevm:v1.12.1_v0.7.0`


export default function LaunchValidators() {
  const { subnetId, chainId, evmChainId, nodesCount } = useWizardStore();
  const [isBootstrapped, setIsBootstrapped] = useState(false);

  const checkMark = (
    <svg
      className="inline-block w-5 h-5 text-green-500 ml-1"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 16 12"
    >
      <path
        stroke="currentColor"
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth="2"
        d="M1 5.917 5.724 10.5 15 1.5"
      />
    </svg>
  );

  return (
    <>
      <h1 className="text-2xl font-medium mb-6">Launch L1 Validators</h1>

      <h3 className="mb-4 font-medium">
        Launch this on each of your {nodesCount} validator node
        {nodesCount > 1 ? 's' : ''}:
      </h3>
      <pre className="bg-gray-100 p-4 rounded-md mb-4">
        {dockerCommand(subnetId)}
      </pre>

      <Note>
        <code className="font-mono bg-blue-100 px-1 py-0.5 rounded">{subnetId}</code> is the subnet ID
      </Note>


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
      </p>


      <p className="mb-4">
        It will take approximately 3 to 10 minutes to bootstrap as it needs to
        sync the P-Chain, which is quick. Note that these nodes are not
        archival; you should read how to set up archival nodes if needed.
      </p>

      <h3 className="mb-4 font-medium">
        You can check the node logs by running:
      </h3>
      <pre className="bg-gray-100 p-4 rounded-md mb-4">
        docker logs -f avalanchego
      </pre>

      <h3 className="mb-4 font-medium">
        To test if the node bootstrapped successfully, run:
      </h3>
      <pre className="bg-gray-100 p-4 rounded-md mb-4">
        {`curl -X POST --data '{ 
"jsonrpc":"2.0", "method":"eth_chainId", "params":[], "id":1 
}' -H 'content-type:application/json;' \\
127.0.0.1:9650/ext/bc/${chainId}/rpc`}
      </pre>

      <Note>
        <code className="font-mono bg-blue-100 px-1 py-0.5 rounded">{chainId}</code> is the chain ID
      </Note>


      <p className="mb-4">
        If your node has bootstrapped successfully, you should see this response:
      </p>

      <pre className="bg-gray-100 p-4 rounded-md mb-4">
        {`{"jsonrpc":"2.0","id":1,"result":"0x${evmChainId.toString(16)}"}`}
      </pre>


      <Note>
        <code className="font-mono bg-blue-100 px-1 py-0.5 rounded">0x{evmChainId.toString(16)}</code> is the hex representation of the chain ID {evmChainId}
      </Note>

      <p className="mb-8">
        Once the validator
        {nodesCount > 1 ? 's' : ''} are bootstrapped, you will launch the RPC
        node{nodesCount > 1 ? 's' : ''} in the next step.
      </p>

      <div className="mb-4">
        <div className="flex items-center">
          <input
            type="checkbox"
            id="bootstrapConfirm"
            checked={isBootstrapped}
            onChange={(e) => setIsBootstrapped(e.target.checked)}
            className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500"
          />
          <label htmlFor="bootstrapConfirm" className="ml-2">
            I confirm my node{nodesCount > 1 ? 's are' : ' is'} bootstrapped and returning the expected chain ID
          </label>
        </div>
        <p className="text-sm text-gray-500 pl-6">Click this to continue</p>
      </div>


      <NextPrev nextDisabled={!isBootstrapped} currentStepName="launch-validators" />
    </>
  );
}
