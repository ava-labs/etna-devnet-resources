import NextPrev from "./ui/NextPrev";

export default function Welcome() {
    return <>
        <h1 className="text-2xl font-medium mb-6">Welcome</h1>

        <p className="mb-4">
            This L1 Launcher will help you getting you're L1 set up. It is aimed for developers that want to launch and maintain the validators and RPC nodes of their L1 on their own infrastructure or provider. This tool is free to use and completely open source.
        </p>

        <ul className="mb-4 list-disc list-inside">
            <li><strong>Create genesis:</strong> Define the L1s native token allocations and precompiles</li>
            <li><strong>Prepare Validators:</strong> Set up the infrastructure for the validators</li>
            <li><strong>Generate keys:</strong> Generate the node credentials and provide them</li>
            <li><strong>Create an L1:</strong> Create the P-Chain records of the L1</li>
            <li><strong>Launch validators:</strong> Launch validators tracking the L1</li>
            <li><strong>Launch an RPC node:</strong> Launch an RPC node for your L1</li>
            <li><strong>Open RPC port:</strong>Open the RPC port & and optinally set up a SSL certificate</li>
            <li><strong>Add to wallet:</strong> Add your L1 to your wallet</li>
            <li><strong>Deploy contracts:</strong> Complete the setup of the L1 by deploying the Validator Manager contracts</li>
        </ul>

        <p className="mb-4">
            If you're looking for a full-service solution that includes hosting, monitoring and maintenance of the L1's validators for you and offers many additional features check out <a href="https://avacloud.io/" target="_blank" className="text-blue-500 hover:text-blue-700 underline">AvaCloud</a>.
        </p>

        <NextPrev nextDisabled={false} currentStepName="welcome" />
    </>;
}
