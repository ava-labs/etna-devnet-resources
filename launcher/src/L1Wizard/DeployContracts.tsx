import { useState, useEffect } from 'react';
import { useWizardStore } from './store';
import { getWalletAddress } from './wallet';
import PoAValidatorManager from "../../contract_compiler/compiled/PoAValidatorManager.json"
import ValidatorMessages from "../../contract_compiler/compiled/ValidatorMessages.json"
import { createPublicClient, createWalletClient, custom, http, Chain, defineChain, keccak256, encodeAbiParameters, parseAbiParameters, toHex } from 'viem';
import { privateKeyToAccount } from 'viem/accounts';

const PROXY_ADMIN_ADDRESS = '0xC0fFEE1234567890aBCdeF1234567890abcDef34' as const;
const PROXY_ADDRESS = '0xfEeDC0DE00000000000000000000000000000000' as const;

// Function selectors
const GET_IMPLEMENTATION_SELECTOR = '0xf9633eab'; // keccak256('getProxyImplementation(address)').slice(0, 10)
const UPGRADE_TO_SELECTOR = '0x3659cfe6'; // keccak256('upgradeTo(address)').slice(0, 10)

interface DeploymentStatus {
    status: 'not_started' | 'deploying' | 'error' | 'success';
    error?: string;
    address?: `0x${string}`;
    txHash?: `0x${string}`;
    currentImplementation?: `0x${string}`;
}

interface DeploymentState {
    validatorMessages: DeploymentStatus;
    validatorManager: DeploymentStatus;
    proxy: DeploymentStatus;
}


export default function DeployContracts() {
    const [deploymentState, setDeploymentState] = useState<DeploymentState>({
        validatorMessages: { status: 'not_started' },
        validatorManager: { status: 'not_started' },
        proxy: { status: 'not_started' },
    });

    const { evmChainId, rpcAddress, rpcLocationType, rpcDomainType, chainId } = useWizardStore();

    const getRpcEndpoint = () => {
        if (rpcLocationType === 'local') return `http://localhost:8080/ext/bc/${chainId}/rpc`;
        if (rpcDomainType === 'no-domain') return `https://${rpcAddress}.nip.io/ext/bc/${chainId}/rpc`;
        return `https://${rpcAddress}/ext/bc/${chainId}/rpc`;
    };

    const getChainConfig = () => defineChain({
        id: Number(evmChainId),
        name: chainId,
        nativeCurrency: {
            decimals: 18,
            name: 'Native Token',
            symbol: 'TOKEN',
        },
        rpcUrls: {
            default: { http: [getRpcEndpoint()] },
            public: { http: [getRpcEndpoint()] },
        },
    });

    const handleCopyToClipboard = async (text: string) => {
        try {
            await navigator.clipboard.writeText(text);
        } catch (err) {
            console.error('Failed to copy:', err);
        }
    };

    const handleDeployValidatorMessages = async () => {
        setDeploymentState(prev => ({
            ...prev,
            validatorMessages: { status: 'deploying' }
        }));

        try {
            if (!window.ethereum) {
                throw new Error('No wallet detected');
            }

            const chain = getChainConfig();
            const publicClient = createPublicClient({
                chain,
                transport: http(),
            });

            const walletClient = createWalletClient({
                chain,
                transport: custom(window.ethereum),
            });

            const [address] = await walletClient.requestAddresses();

            // Get bytecode and ensure it has 0x prefix
            const bytecode = ValidatorMessages.bytecode.object;
            const bytecodeWithPrefix = bytecode.startsWith('0x') ? bytecode : `0x${bytecode}`;

            // Deploy ValidatorMessages contract
            const hash = await walletClient.deployContract({
                chain,
                abi: ValidatorMessages.abi,
                bytecode: bytecodeWithPrefix as `0x${string}`,
                account: address,
            });

            // Wait for deployment and get contract address
            const receipt = await publicClient.waitForTransactionReceipt({ hash });

            if (!receipt.contractAddress) {
                throw new Error('Contract address not found in receipt');
            }

            setDeploymentState(prev => ({
                ...prev,
                validatorMessages: {
                    status: 'success',
                    address: receipt.contractAddress as `0x${string}`,
                    txHash: hash,
                }
            }));
        } catch (err: any) {
            console.error('Deployment failed:', err);
            setDeploymentState(prev => ({
                ...prev,
                validatorMessages: {
                    status: 'error',
                    error: err.message || 'Failed to deploy ValidatorMessages'
                }
            }));
        }
    };

    const calculateLibraryHash = (libraryPath: string) => {
        // Calculate keccak256 of the fully qualified library name
        const hash = keccak256(
            new TextEncoder().encode(libraryPath)
        ).slice(2);
        // Take first 34 characters (17 bytes)
        return hash.slice(0, 34);
    };

    const handleDeployValidatorManager = async () => {
        if (!deploymentState.validatorMessages.address) {
            throw new Error('ValidatorMessages address not found');
        }

        setDeploymentState(prev => ({
            ...prev,
            validatorManager: { status: 'deploying' }
        }));

        try {
            if (!window.ethereum) {
                throw new Error('No wallet detected');
            }

            const chain = getChainConfig();
            const publicClient = createPublicClient({
                chain,
                transport: http(),
            });

            const walletClient = createWalletClient({
                chain,
                transport: custom(window.ethereum),
            });

            const [address] = await walletClient.requestAddresses();

            // Get bytecode and ensure it has 0x prefix
            const bytecode = PoAValidatorManager.bytecode.object;
            const bytecodeWithPrefix = bytecode.startsWith('0x') ? bytecode : `0x${bytecode}`;

            // Get library path from linkReferences
            const libraryPath = `${Object.keys(PoAValidatorManager.bytecode.linkReferences)[0]}:${Object.keys(Object.values(PoAValidatorManager.bytecode.linkReferences)[0])[0]}`
            const libraryHash = calculateLibraryHash(libraryPath);
            const libraryPlaceholder = `__$${libraryHash}$__`;


            // Replace library placeholder with actual address using split/join
            const linkedBytecode = bytecodeWithPrefix
                .split(libraryPlaceholder)
                .join(deploymentState.validatorMessages.address.slice(2).padStart(40, '0'));

            if (linkedBytecode.includes("$__")) {
                throw new Error("Failed to replace library placeholder with actual address")
            }

            // Deploy ValidatorManager contract with linked bytecode
            const hash = await walletClient.deployContract({
                chain,
                abi: PoAValidatorManager.abi,
                bytecode: linkedBytecode as `0x${string}`,
                account: address,
            });

            // Wait for deployment and get contract address
            const receipt = await publicClient.waitForTransactionReceipt({ hash });

            if (!receipt.contractAddress) {
                throw new Error('Contract address not found in receipt');
            }

            setDeploymentState(prev => ({
                ...prev,
                validatorManager: {
                    status: 'success',
                    address: receipt.contractAddress as `0x${string}`,
                    txHash: hash,
                }
            }));
        } catch (err: any) {
            console.error('Deployment failed:', err);
            setDeploymentState(prev => ({
                ...prev,
                validatorManager: {
                    status: 'error',
                    error: err.message || 'Failed to deploy ValidatorManager'
                }
            }));
        }
    };

    useEffect(() => {
        checkProxyImplementation();
    }, []);

    const checkProxyImplementation = async () => {
        try {
            if (!window.ethereum) {
                throw new Error('No wallet detected');
            }

            const chain = getChainConfig();
            const publicClient = createPublicClient({
                chain,
                transport: http(),
            });

            const walletClient = createWalletClient({
                chain,
                transport: custom(window.ethereum),
            });

            const [address] = await walletClient.requestAddresses();

            // Call getProxyImplementation using staticCall with proxy address parameter
            const result = await publicClient.call({
                account: address,
                to: PROXY_ADMIN_ADDRESS,
                data: `${GET_IMPLEMENTATION_SELECTOR}${PROXY_ADDRESS.slice(2).padStart(64, '0')}` as `0x${string}`,
            });

            if (result.data) {
                const implementation = `0x${result.data.slice(-40)}` as `0x${string}`;
                setDeploymentState(prev => ({
                    ...prev,
                    proxy: {
                        ...prev.proxy,
                        status: 'success',
                        currentImplementation: implementation
                    }
                }));
            }
        } catch (err: any) {
            console.error('Failed to check proxy implementation:', err);
            setDeploymentState(prev => ({
                ...prev,
                proxy: {
                    status: 'error',
                    error: err.message || 'Failed to check proxy implementation'
                }
            }));
        }
    };

    const handleUpdateProxyAddress = async () => {
        if (!deploymentState.validatorManager.address) {
            throw new Error('ValidatorManager address not found');
        }

        setDeploymentState(prev => ({
            ...prev,
            proxy: { status: 'deploying' }
        }));

        try {
            if (!window.ethereum) {
                throw new Error('No wallet detected');
            }

            const chain = getChainConfig();
            const walletClient = createWalletClient({
                chain,
                transport: custom(window.ethereum),
            });

            const [address] = await walletClient.requestAddresses();

            // Encode upgrade function call
            const upgradeData = encodeAbiParameters(
                parseAbiParameters('address implementation'),
                [deploymentState.validatorManager.address]
            );

            // Call upgrade function
            const hash = await walletClient.sendTransaction({
                account: address,
                to: PROXY_ADMIN_ADDRESS,
                data: `${UPGRADE_TO_SELECTOR}${upgradeData.slice(2)}` as `0x${string}`,
            });

            // Wait for transaction to complete
            const publicClient = createPublicClient({
                chain,
                transport: http(),
            });
            await publicClient.waitForTransactionReceipt({ hash });

            // Check new implementation
            const result = await publicClient.call({
                account: PROXY_ADMIN_ADDRESS,
                to: PROXY_ADMIN_ADDRESS,
                data: `${GET_IMPLEMENTATION_SELECTOR}${deploymentState.validatorManager.address!.slice(2).padStart(64, '0')}` as `0x${string}`,
            });

            if (result.data) {
                const implementation = `0x${result.data.slice(-40)}` as `0x${string}`;
                setDeploymentState(prev => ({
                    ...prev,
                    proxy: {
                        status: 'success',
                        txHash: hash,
                        currentImplementation: implementation
                    }
                }));
            } else {
                throw new Error('Failed to get implementation address');
            }
        } catch (err: any) {
            console.error('Proxy update failed:', err);
            setDeploymentState(prev => ({
                ...prev,
                proxy: {
                    status: 'error',
                    error: err.message || 'Failed to update proxy'
                }
            }));
        }
    };

    const renderDeploymentStatus = (deployment: DeploymentStatus, title: string) => {
        const statusColors = {
            not_started: 'bg-gray-50 border-gray-200',
            deploying: 'bg-blue-50 border-blue-200',
            error: 'bg-red-50 border-red-200',
            success: 'bg-green-50 border-green-200'
        };

        const statusText = {
            not_started: 'Not deployed yet',
            deploying: 'Deploying...',
            error: 'Deployment failed',
            success: 'Deployed successfully'
        };

        return (
            <div className={`p-4 rounded-lg border ${statusColors[deployment.status]} mb-4`}>
                <div className="flex items-center justify-between mb-2">
                    <h3 className="font-medium">{title}</h3>
                    <span className={`text-sm ${deployment.status === 'error' ? 'text-red-600' : ''}`}>
                        {statusText[deployment.status]}
                    </span>
                </div>

                {deployment.error && (
                    <div className="text-sm text-red-600 mb-2">{deployment.error}</div>
                )}

                {deployment.address && (
                    <div className="mb-2">
                        <div className="text-sm text-gray-500 mb-1">Contract Address:</div>
                        <div className="flex items-center bg-white rounded p-2 border border-gray-100">
                            <code className="font-mono text-sm flex-1 break-all">{deployment.address}</code>
                            <button
                                onClick={() => handleCopyToClipboard(deployment.address!)}
                                className="ml-2 px-2 py-1 text-xs bg-gray-100 hover:bg-gray-200 rounded"
                            >
                                Copy
                            </button>
                        </div>
                    </div>
                )}

                {deployment.txHash && (
                    <div>
                        <div className="text-sm text-gray-500 mb-1">Transaction:</div>
                        <div className="flex items-center bg-white rounded p-2 border border-gray-100">
                            <code className="font-mono text-sm flex-1 break-all">{deployment.txHash}</code>
                            <button
                                onClick={() => handleCopyToClipboard(deployment.txHash!)}
                                className="ml-2 px-2 py-1 text-xs bg-gray-100 hover:bg-gray-200 rounded"
                            >
                                Copy
                            </button>
                        </div>
                    </div>
                )}

                {deployment.currentImplementation && (
                    <div className="mb-2">
                        <div className="text-sm text-gray-500 mb-1">Current Implementation:</div>
                        <div className="flex items-center bg-white rounded p-2 border border-gray-100">
                            <code className="font-mono text-sm flex-1 break-all">{deployment.currentImplementation}</code>
                            <button
                                onClick={() => handleCopyToClipboard(deployment.currentImplementation!)}
                                className="ml-2 px-2 py-1 text-xs bg-gray-100 hover:bg-gray-200 rounded"
                            >
                                Copy
                            </button>
                        </div>
                    </div>
                )}

                {deployment.status === 'not_started' && (
                    <button
                        onClick={
                            title === 'ValidatorMessages' ? handleDeployValidatorMessages :
                                title === 'ValidatorManager' ? handleDeployValidatorManager :
                                    handleUpdateProxyAddress
                        }
                        disabled={
                            title === 'ValidatorManager' && deploymentState.validatorMessages.status !== 'success' ||
                            title === 'Proxy' && deploymentState.validatorManager.status !== 'success'
                        }
                        className={`mt-2 w-full p-2 rounded ${(title === 'ValidatorManager' && deploymentState.validatorMessages.status !== 'success') ||
                            (title === 'Proxy' && deploymentState.validatorManager.status !== 'success')
                            ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                            : 'bg-blue-500 text-white hover:bg-blue-600'
                            }`}
                    >
                        Deploy
                    </button>
                )}

                {title === 'Proxy' && deployment.status === 'success' && deploymentState.validatorManager.status === 'success' && (
                    <button
                        onClick={handleUpdateProxyAddress}
                        className="mt-2 w-full p-2 rounded bg-blue-500 text-white hover:bg-blue-600"
                    >
                        Update Implementation
                    </button>
                )}
            </div>
        );
    };

    return (
        <div className="max-w-3xl mx-auto">
            <h1 className="text-2xl font-medium mb-6">Deploy Contracts</h1>
            {renderDeploymentStatus(deploymentState.validatorMessages, 'ValidatorMessages')}
            {renderDeploymentStatus(deploymentState.validatorManager, 'ValidatorManager')}
            {renderDeploymentStatus(deploymentState.proxy, 'Proxy')}
        </div>
    );
}
