const dockerInstallation = `sudo yum update -y
sudo yum -y install docker

sudo service docker start

sudo usermod -a -G docker ec2-user
sudo chmod 666 /var/run/docker.sock
docker version
`

import { useWizardStore } from './store';
import NextPrev from './ui/NextPrev';

export default function PrepareValidators() {
    const { nodesCount, setNodesCount } = useWizardStore();

    const nodeConfigurations = {
        1: { colorClass: 'bg-red-100 text-red-800', badge: 'Dev' },
        3: { colorClass: 'bg-green-100 text-green-800', badge: 'Testnet' },
        5: { colorClass: 'bg-blue-100 text-blue-800', badge: 'Mainnet' },
    };

    return <>
        <h1 className="text-2xl font-medium mb-6">
            Prepare your Validator Nodes
        </h1>

        <p className="mb-4">
            This step will guide you through preparing your validator nodes. for long-lived environment we recommend using a infrastructure provider such as AWS or Google Cloud. For short-lived environments you can use your local machine or any other machine you have access to.
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
            If you do not want to use Docker, you can follow the instructions <a href="https://github.com/ava-labs/avalanchego?tab=readme-ov-file#installation" target="_blank">here</a>.
        </p>

        <h3 className="mb-4 font-medium">How many nodes do you want to run?</h3>
        <ul className="mb-4 items-center w-full text-sm font-medium text-gray-900 bg-white border border-gray-200 rounded-lg sm:flex">
            {Object.entries(nodeConfigurations).map(([count, config]) => {
                return (
                    <li key={count} className="w-full border-b border-gray-200 sm:border-b-0 sm:border-r last:border-r-0 ">
                        <div className="flex items-center ps-3">
                            <input
                                id={`nodes-${count}`}
                                type="radio"
                                checked={nodesCount === Number(count)}
                                onChange={() => setNodesCount(Number(count))}
                                name="nodes-count"
                                className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 focus:ring-0"
                            />
                            <label htmlFor={`nodes-${count}`} className="w-full py-3 ms-2 text-sm font-medium text-gray-900">
                                {count} {count === '1' ? 'Node' : 'Nodes'}
                                {config.badge && (
                                    <span className={`ms-2 ${config.colorClass} text-xs font-medium px-2.5 py-0.5 rounded-full`}>
                                        {config.badge}
                                    </span>
                                )}
                            </label>
                        </div>
                    </li>
                );
            })}
        </ul>

        <NextPrev nextDisabled={!nodesCount} currentStepName="prepare-validators" />
    </>
}
